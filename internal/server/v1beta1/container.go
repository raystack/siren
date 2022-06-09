package v1beta1

import (
	"net/http"

	"github.com/odpf/siren/core/alert"
	"github.com/odpf/siren/core/namespace"
	"github.com/odpf/siren/core/provider"
	"github.com/odpf/siren/core/receiver"
	"github.com/odpf/siren/core/subscription"
	"github.com/odpf/siren/internal/store"
	"github.com/odpf/siren/plugins/receivers/slack"

	"github.com/odpf/siren/core/rules"
	"github.com/odpf/siren/core/templates"
	"github.com/odpf/siren/domain"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type Container struct {
	TemplatesService    domain.TemplatesService
	RulesService        domain.RuleService
	AlertService        AlertService
	NotifierServices    domain.NotifierServices
	ProviderService     ProviderService
	NamespaceService    NamespaceService
	ReceiverService     domain.ReceiverService
	SubscriptionService domain.SubscriptionService
}

func InitContainer(repositories *store.RepositoryContainer, db *gorm.DB, c *domain.Config, httpClient *http.Client) (*Container, error) {
	templatesService := templates.NewService(repositories.TemplatesRepository)
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
	rulesService := rules.NewService(
		repositories.RuleRepository,
		templatesService,
		namespaceService,
		providerService,
	)
	receiverService, err := receiver.NewService(repositories.ReceiverRepository, httpClient, c.EncryptionKey)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create receiver service")
	}
	subscriptionService, err := subscription.NewService(repositories.SubscriptionRepository, repositories.ProviderRepository,
		repositories.NamespaceRepository, repositories.ReceiverRepository, c.EncryptionKey)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create subscriptions service")
	}

	return &Container{
		TemplatesService: templatesService,
		RulesService:     rulesService,
		AlertService:     alertHistoryService,
		NotifierServices: domain.NotifierServices{
			Slack: slack.NewService(),
		},
		ProviderService:     providerService,
		NamespaceService:    namespaceService,
		ReceiverService:     receiverService,
		SubscriptionService: subscriptionService,
	}, nil
}

func (container *Container) MigrateAll(db *gorm.DB) error {
	err := container.TemplatesService.Migrate()
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
