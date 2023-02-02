package cli

import (
	"context"
	"fmt"

	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/odpf/salt/db"
	"github.com/odpf/salt/log"
	"github.com/odpf/siren/config"
	"github.com/odpf/siren/core/alert"
	"github.com/odpf/siren/core/namespace"
	"github.com/odpf/siren/core/notification"
	"github.com/odpf/siren/core/provider"
	"github.com/odpf/siren/core/receiver"
	"github.com/odpf/siren/core/rule"
	"github.com/odpf/siren/core/subscription"
	"github.com/odpf/siren/core/template"
	"github.com/odpf/siren/internal/api"
	"github.com/odpf/siren/internal/store/postgres"
	"github.com/odpf/siren/pkg/pgc"
	"github.com/odpf/siren/pkg/secret"
	"github.com/odpf/siren/pkg/telemetry"
	"github.com/odpf/siren/plugins/providers/cortex"
	"github.com/odpf/siren/plugins/receivers/file"
	"github.com/odpf/siren/plugins/receivers/httpreceiver"
	"github.com/odpf/siren/plugins/receivers/pagerduty"
	"github.com/odpf/siren/plugins/receivers/slack"
)

func InitDeps(
	ctx context.Context,
	logger log.Logger,
	cfg config.Config,
	queue notification.Queuer,
) (*api.Deps, *newrelic.Application, *pgc.Client, map[string]notification.Notifier, error) {

	telemetry.Init(ctx, cfg.Telemetry, logger)
	nrApp, err := newrelic.NewApplication(
		newrelic.ConfigAppName(cfg.Telemetry.ServiceName),
		newrelic.ConfigLicense(cfg.Telemetry.NewRelicAPIKey),
	)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	dbClient, err := db.New(cfg.DB)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	pgClient, err := pgc.NewClient(logger, dbClient)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	encryptor, err := secret.New(cfg.Service.EncryptionKey)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("cannot initialize encryptor: %w", err)
	}

	idempotencyRepository := postgres.NewIdempotencyRepository(pgClient)
	templateRepository := postgres.NewTemplateRepository(pgClient)
	templateService := template.NewService(templateRepository)

	alertRepository := postgres.NewAlertRepository(pgClient)

	providerRepository := postgres.NewProviderRepository(pgClient)
	providerService := provider.NewService(providerRepository)

	namespaceRepository := postgres.NewNamespaceRepository(pgClient)

	cortexPluginService := cortex.NewPluginService(logger, cfg.Providers.Cortex)

	alertHistoryService := alert.NewService(alertRepository, map[string]alert.AlertTransformer{
		provider.TypeCortex: cortexPluginService,
	})
	namespaceService := namespace.NewService(encryptor, namespaceRepository, providerService, map[string]namespace.ConfigSyncer{
		provider.TypeCortex: cortexPluginService,
	})

	ruleRepository := postgres.NewRuleRepository(pgClient)
	ruleService := rule.NewService(
		ruleRepository,
		templateService,
		namespaceService,
		map[string]rule.RuleUploader{
			provider.TypeCortex: cortexPluginService,
		},
	)

	// plugin receiver services
	slackPluginService := slack.NewPluginService(cfg.Receivers.Slack, encryptor)
	pagerDutyPluginService := pagerduty.NewPluginService(cfg.Receivers.Pagerduty)
	httpreceiverPluginService := httpreceiver.NewPluginService(logger, cfg.Receivers.HTTPReceiver)
	filePluginService := file.NewPluginService()

	receiverRepository := postgres.NewReceiverRepository(pgClient)
	receiverService := receiver.NewService(
		receiverRepository,
		map[string]receiver.ConfigResolver{
			receiver.TypeSlack:     slackPluginService,
			receiver.TypeHTTP:      httpreceiverPluginService,
			receiver.TypePagerDuty: pagerDutyPluginService,
			receiver.TypeFile:      filePluginService,
		},
	)

	subscriptionRepository := postgres.NewSubscriptionRepository(pgClient)
	subscriptionService := subscription.NewService(
		subscriptionRepository,
		namespaceService,
		receiverService,
	)

	// notification
	notifierRegistry := map[string]notification.Notifier{
		receiver.TypeSlack:     slackPluginService,
		receiver.TypePagerDuty: pagerDutyPluginService,
		receiver.TypeHTTP:      httpreceiverPluginService,
		receiver.TypeFile:      filePluginService,
	}

	notificationService := notification.NewService(logger, queue, idempotencyRepository, receiverService, subscriptionService, notifierRegistry)

	return &api.Deps{
			TemplateService:     templateService,
			RuleService:         ruleService,
			AlertService:        alertHistoryService,
			ProviderService:     providerService,
			NamespaceService:    namespaceService,
			ReceiverService:     receiverService,
			SubscriptionService: subscriptionService,
			NotificationService: notificationService,
		}, nrApp, pgClient, notifierRegistry,
		nil
}
