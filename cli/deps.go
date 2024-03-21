package cli

import (
	"context"
	"fmt"

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
	"github.com/goto/siren/plugins/providers"
	"github.com/goto/siren/plugins/receivers/file"
	"github.com/goto/siren/plugins/receivers/httpreceiver"
	"github.com/goto/siren/plugins/receivers/pagerduty"
	"github.com/goto/siren/plugins/receivers/slack"
	"github.com/goto/siren/plugins/receivers/slackchannel"
)

func InitDeps(
	ctx context.Context,
	logger saltlog.Logger,
	cfg config.Config,
	queue notification.Queuer,
	withProviderPlugin bool,
) (*api.Deps, *pgc.Client, map[string]notification.Notifier, *providers.PluginManager, error) {

	pgClient, err := pgc.NewClient(logger, cfg.DB)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	encryptor, err := secret.New(cfg.Service.EncryptionKey)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("cannot initialize encryptor: %w", err)
	}

	templateRepository := postgres.NewTemplateRepository(pgClient)
	templateService := template.NewService(templateRepository)

	logRepository := postgres.NewLogRepository(pgClient)
	logService := log.NewService(logRepository)

	// TODO: need to figure out the way on how to nicely load deps without plugins dependency
	// in case the caller does not need that
	var providerService api.ProviderService
	var configSyncers = make(map[string]namespace.ConfigSyncer, 0)
	var alertTransformers = make(map[string]alert.AlertTransformer, 0)
	var ruleUploaders = make(map[string]rule.RuleUploader, 0)
	var providersPluginManager *providers.PluginManager

	if withProviderPlugin {
		providersPluginManager = providers.NewPluginManager(logger, cfg.Providers)
		providerPluginClients := providersPluginManager.InitClients()
		providerPlugins, err := providersPluginManager.DispenseClients(providerPluginClients)
		if err != nil {
			return nil, nil, nil, nil, err
		}
		if err := providersPluginManager.InitConfigs(ctx, providerPlugins, cfg.Log.Level); err != nil {
			return nil, nil, nil, nil, err
		}

		var supportedProviderTypes = []string{}
		for typ := range providerPlugins {
			supportedProviderTypes = append(supportedProviderTypes, typ)
		}
		providerRepository := postgres.NewProviderRepository(pgClient)
		providerService = provider.NewService(providerRepository, supportedProviderTypes)

		if len(providerPlugins) == 0 {
			logger.Warn("no provider plugins found!")
		}

		for k, pc := range providerPlugins {
			alertTransformers[k] = pc.(alert.AlertTransformer)
			configSyncers[k] = pc.(namespace.ConfigSyncer)
			ruleUploaders[k] = pc.(rule.RuleUploader)
		}
	}

	alertRepository := postgres.NewAlertRepository(pgClient)
	alertService := alert.NewService(
		alertRepository,
		logService,
		alertTransformers,
	)

	namespaceRepository := postgres.NewNamespaceRepository(pgClient)
	namespaceService := namespace.NewService(encryptor, namespaceRepository, providerService, configSyncers)

	ruleRepository := postgres.NewRuleRepository(pgClient)
	ruleService := rule.NewService(
		ruleRepository,
		templateService,
		namespaceService,
		ruleUploaders,
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
			TemplateService:       templateService,
		},
		cfg.Service.EnableSilenceFeature,
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
	}, pgClient, notifierRegistry, providersPluginManager, nil
}
