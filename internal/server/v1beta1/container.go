package v1beta1

import (
	"net/http"

	"github.com/odpf/siren/config"
	"github.com/odpf/siren/core/alert"
	"github.com/odpf/siren/core/namespace"
	"github.com/odpf/siren/core/provider"
	"github.com/odpf/siren/core/receiver"
	"github.com/odpf/siren/core/rule"
	"github.com/odpf/siren/core/subscription"
	"github.com/odpf/siren/core/subscription/alertmanager"
	"github.com/odpf/siren/core/template"
	"github.com/odpf/siren/internal/store"
	slackclient "github.com/odpf/siren/plugins/receivers/http"
	"github.com/odpf/siren/plugins/receivers/slack"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type Container struct {
	TemplateService     TemplateService
	RulesService        RuleService
	AlertService        AlertService
	NotifierServices    NotifierServices
	ProviderService     ProviderService
	NamespaceService    NamespaceService
	ReceiverService     ReceiverService
	SubscriptionService SubscriptionService
}

func InitContainer(repositories *store.RepositoryContainer, db *gorm.DB, c *config.Config, httpClient *http.Client) (*Container, error) {
	templateService := template.NewService(repositories.TemplatesRepository)
	alertHistoryService := alert.NewService(repositories.AlertRepository)

	providerService := provider.NewService(repositories.ProviderRepository)

	encryptionTransformer, err := namespace.NewTransformer(c.EncryptionKey)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create namespace transformer")
	}
	namespaceService, err := namespace.NewService(repositories.NamespaceRepository, encryptionTransformer)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create namespace service")
	}

	cortexClient, err := rule.NewCortexClient(c.Cortex.Address)
	if err != nil {
		return nil, errors.Wrap(err, "failed to init cortex client")
	}

	rulesService := rule.NewService(
		repositories.RuleRepository,
		templateService,
		namespaceService,
		providerService,
		cortexClient,
	)
	receiverExchange := slackclient.NewSlackClient(&http.Client{})
	receiverSlackHelper, err := receiver.NewSlackHelper(receiverExchange, c.EncryptionKey)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create slack helper")
	}
	receiverService, err := receiver.NewService(repositories.ReceiverRepository, receiverSlackHelper, slack.NewService())
	if err != nil {
		return nil, errors.Wrap(err, "failed to create receiver service")
	}

	amClient, err := alertmanager.NewClient(config.CortexConfig{Address: c.Cortex.Address})
	if err != nil {
		return nil, errors.Wrap(err, "failed to create alert manager client")
	}
	subscriptionService := subscription.NewService(repositories.SubscriptionRepository, providerService, namespaceService, receiverService, amClient)

	return &Container{
		TemplateService: templateService,
		RulesService:    rulesService,
		AlertService:    alertHistoryService,
		NotifierServices: NotifierServices{
			Slack: slack.NewService(),
		},
		ProviderService:     providerService,
		NamespaceService:    namespaceService,
		ReceiverService:     receiverService,
		SubscriptionService: subscriptionService,
	}, nil
}

func (container *Container) MigrateAll(db *gorm.DB) error {
	err := container.TemplateService.Migrate()
	if err != nil {
		return err
	}
	err = container.RulesService.Migrate()
	if err != nil {
		return err
	}
	err = container.AlertService.Migrate()
	if err != nil {
		return err
	}
	err = container.ProviderService.Migrate()
	if err != nil {
		return err
	}
	err = container.NamespaceService.Migrate()
	if err != nil {
		return err
	}
	err = container.ReceiverService.Migrate()
	if err != nil {
		return err
	}
	err = container.SubscriptionService.Migrate()
	if err != nil {
		return err
	}
	return nil
}
