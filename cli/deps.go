package cli

import (
	"context"
	"fmt"

	"github.com/goto/salt/db"
	saltlog "github.com/goto/salt/log"
	"github.com/goto/siren/config"
	"github.com/goto/siren/core/alert"
	"github.com/goto/siren/core/log"
	"github.com/goto/siren/core/namespace"
	"github.com/goto/siren/core/notification"
	"github.com/goto/siren/core/provider"
	"github.com/goto/siren/core/receiver"
	"github.com/goto/siren/core/rule"
	"github.com/goto/siren/core/silence"
	"github.com/goto/siren/core/subscription"
	"github.com/goto/siren/core/template"
	"github.com/goto/siren/internal/api"
	"github.com/goto/siren/internal/store/postgres"
	"github.com/goto/siren/pkg/pgc"
	"github.com/goto/siren/pkg/secret"
	"github.com/goto/siren/pkg/telemetry"
	"github.com/goto/siren/plugins/providers/cortex"
	"github.com/goto/siren/plugins/receivers/file"
	"github.com/goto/siren/plugins/receivers/httpreceiver"
	"github.com/goto/siren/plugins/receivers/pagerduty"
	"github.com/goto/siren/plugins/receivers/slack"
	"github.com/goto/siren/plugins/receivers/slackchannel"
	"github.com/newrelic/go-agent/v3/newrelic"
)

func InitDeps(
	ctx context.Context,
	logger saltlog.Logger,
	cfg config.Config,
	queue notification.Queuer,
) (*api.Deps, *newrelic.Application, *pgc.Client, map[string]notification.Notifier, error) {

	telemetry.Init(ctx, cfg.Telemetry, logger)

	nrApp, err := newrelic.NewApplication(
		newrelic.ConfigAppName(cfg.Telemetry.NewRelicAppName),
		newrelic.ConfigLicense(cfg.Telemetry.NewRelicAPIKey),
	)
	if err != nil {
		logger.Warn("failed to init newrelic", "err", err)
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

	templateRepository := postgres.NewTemplateRepository(pgClient)
	templateService := template.NewService(templateRepository)

	providerRepository := postgres.NewProviderRepository(pgClient)
	providerService := provider.NewService(providerRepository)

	logRepository := postgres.NewLogRepository(pgClient)
	logService := log.NewService(logRepository)

	cortexPluginService := cortex.NewPluginService(logger, cfg.Providers.Cortex)
	alertRepository := postgres.NewAlertRepository(pgClient)
	alertService := alert.NewService(
		alertRepository,
		logService,
		map[string]alert.AlertTransformer{
			provider.TypeCortex: cortexPluginService,
		},
	)

	namespaceRepository := postgres.NewNamespaceRepository(pgClient)
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

	silenceRepository := postgres.NewSilenceRepository(pgClient)
	silenceService := silence.NewService(silenceRepository)

	// plugin receiver services
	slackPluginService := slack.NewPluginService(cfg.Receivers.Slack, encryptor)
	slackChannelPluginService := slackchannel.NewPluginService(cfg.Receivers.Slack, encryptor)
	pagerDutyPluginService := pagerduty.NewPluginService(cfg.Receivers.Pagerduty)
	httpreceiverPluginService := httpreceiver.NewPluginService(logger, cfg.Receivers.HTTPReceiver)
	filePluginService := file.NewPluginService()

	receiverRepository := postgres.NewReceiverRepository(pgClient)
	receiverService := receiver.NewService(
		receiverRepository,
		map[string]receiver.ConfigResolver{
			receiver.TypeSlack:        slackPluginService,
			receiver.TypeSlackChannel: slackChannelPluginService,
			receiver.TypeHTTP:         httpreceiverPluginService,
			receiver.TypePagerDuty:    pagerDutyPluginService,
			receiver.TypeFile:         filePluginService,
		},
	)

	subscriptionRepository := postgres.NewSubscriptionRepository(pgClient)
	subscriptionService := subscription.NewService(
		subscriptionRepository,
		logService,
		namespaceService,
		receiverService,
	)

	// notification
	notifierRegistry := map[string]notification.Notifier{
		receiver.TypeSlack:        slackPluginService,
		receiver.TypeSlackChannel: slackChannelPluginService,
		receiver.TypePagerDuty:    pagerDutyPluginService,
		receiver.TypeHTTP:         httpreceiverPluginService,
		receiver.TypeFile:         filePluginService,
	}

	idempotencyRepository := postgres.NewIdempotencyRepository(pgClient)
	notificationRepository := postgres.NewNotificationRepository(pgClient)
	notificationService := notification.NewService(
		logger,
		cfg.Notification,
		notificationRepository,
		queue,
		notifierRegistry,
		notification.Deps{
			LogService:            logService,
			IdempotencyRepository: idempotencyRepository,
			ReceiverService:       receiverService,
			SubscriptionService:   subscriptionService,
			SilenceService:        silenceService,
			AlertService:          alertService,
		},
	)

	return &api.Deps{
			TemplateService:     templateService,
			RuleService:         ruleService,
			AlertService:        alertService,
			ProviderService:     providerService,
			NamespaceService:    namespaceService,
			ReceiverService:     receiverService,
			SubscriptionService: subscriptionService,
			NotificationService: notificationService,
			SilenceService:      silenceService,
		}, nrApp, pgClient, notifierRegistry,
		nil
}
