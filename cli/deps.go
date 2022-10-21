package cli

import (
	"fmt"

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
	"github.com/odpf/siren/pkg/httpclient"
	"github.com/odpf/siren/pkg/retry"
	"github.com/odpf/siren/pkg/secret"
	"github.com/odpf/siren/plugins/providers/cortex"
	"github.com/odpf/siren/plugins/receivers/file"
	"github.com/odpf/siren/plugins/receivers/httpreceiver"
	"github.com/odpf/siren/plugins/receivers/pagerduty"
	"github.com/odpf/siren/plugins/receivers/slack"
)

type ReceiverClient struct {
	SlackClient        *slack.Client
	PagerDutyClient    *pagerduty.Client
	HTTPReceiverClient *httpreceiver.Client
}

type ProviderClient struct {
	CortexClient *cortex.Client
}

func InitAPIDeps(
	logger log.Logger,
	cfg config.Config,
	pgClient *postgres.Client,
	encryptor *secret.Crypto,
	queue notification.Queuer,
) (*api.Deps, *ReceiverClient, *ProviderClient, map[string]notification.Notifier, error) {
	templateRepository := postgres.NewTemplateRepository(pgClient)
	templateService := template.NewService(templateRepository)

	alertRepository := postgres.NewAlertRepository(pgClient)
	alertHistoryService := alert.NewService(alertRepository)

	providerRepository := postgres.NewProviderRepository(pgClient)
	providerService := provider.NewService(providerRepository)

	namespaceRepository := postgres.NewNamespaceRepository(pgClient)
	namespaceService := namespace.NewService(encryptor, namespaceRepository)

	cortexClient, err := cortex.NewClient(cortex.Config{Address: cfg.Cortex.Address})
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("failed to init cortex client: %w", err)
	}
	cortexProviderService := cortex.NewProviderService(cortexClient)

	ruleRepository := postgres.NewRuleRepository(pgClient)
	ruleService := rule.NewService(
		ruleRepository,
		templateService,
		namespaceService,
		map[string]rule.RuleUploader{
			provider.TypeCortex: cortexProviderService,
		},
	)

	// plugin receiver services
	slackHTTPClient := httpclient.New(cfg.Receivers.Slack.HTTPClient)
	slackRetrier := retry.New(cfg.Receivers.Slack.Retry)
	slackClient := slack.NewClient(
		cfg.Receivers.Slack,
		slack.ClientWithHTTPClient(slackHTTPClient),
		slack.ClientWithRetrier(slackRetrier),
	)
	pagerdutyHTTPClient := httpclient.New(cfg.Receivers.Pagerduty.HTTPClient)
	pagerdutyRetrier := retry.New(cfg.Receivers.Slack.Retry)
	pagerdutyClient := pagerduty.NewClient(
		cfg.Receivers.Pagerduty,
		pagerduty.ClientWithHTTPClient(pagerdutyHTTPClient),
		pagerduty.ClientWithRetrier(pagerdutyRetrier),
	)
	httpreceiverHTTPClient := httpclient.New(cfg.Receivers.HTTPReceiver.HTTPClient)
	httpreceiverRetrier := retry.New(cfg.Receivers.Slack.Retry)
	httpreceiverClient := httpreceiver.NewClient(
		logger,
		cfg.Receivers.HTTPReceiver,
		httpreceiver.ClientWithHTTPClient(httpreceiverHTTPClient),
		httpreceiver.ClientWithRetrier(httpreceiverRetrier),
	)

	slackReceiverService := slack.NewReceiverService(slackClient, encryptor)
	httpReceiverService := httpreceiver.NewReceiverService()
	pagerDutyReceiverService := pagerduty.NewReceiverService()
	fileReceiverService := file.NewReceiverService()

	receiverRepository := postgres.NewReceiverRepository(pgClient)
	receiverService := receiver.NewService(
		receiverRepository,
		map[string]receiver.ConfigResolver{
			receiver.TypeSlack:     slackReceiverService,
			receiver.TypeHTTP:      httpReceiverService,
			receiver.TypePagerDuty: pagerDutyReceiverService,
			receiver.TypeFile:      fileReceiverService,
		},
	)

	subscriptionRepository := postgres.NewSubscriptionRepository(pgClient)
	subscriptionService := subscription.NewService(
		subscriptionRepository,
		namespaceService,
		receiverService,
		subscription.RegisterProviderPlugin(provider.TypeCortex, cortexProviderService),
	)

	// notification
	slackNotificationService := slack.NewNotificationService(slackClient, encryptor)
	pagerdutyNotificationService := pagerduty.NewNotificationService(pagerdutyClient)
	httpreceiverNotificationService := httpreceiver.NewNotificationService(httpreceiverClient)
	fileNotificationService := file.NewNotificationService()

	notifierRegistry := map[string]notification.Notifier{
		receiver.TypeSlack:     slackNotificationService,
		receiver.TypePagerDuty: pagerdutyNotificationService,
		receiver.TypeHTTP:      httpreceiverNotificationService,
		receiver.TypeFile:      fileNotificationService,
	}

	notificationService := notification.NewService(logger, queue, receiverService, subscriptionService, notifierRegistry)

	return &api.Deps{
			TemplateService:     templateService,
			RuleService:         ruleService,
			AlertService:        alertHistoryService,
			ProviderService:     providerService,
			NamespaceService:    namespaceService,
			ReceiverService:     receiverService,
			SubscriptionService: subscriptionService,
			NotificationService: notificationService,
		}, &ReceiverClient{
			SlackClient:        slackClient,
			PagerDutyClient:    pagerdutyClient,
			HTTPReceiverClient: httpreceiverClient,
		}, &ProviderClient{
			CortexClient: cortexClient,
		}, notifierRegistry,
		nil
}
