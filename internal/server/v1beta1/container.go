package v1beta1

import (
	"net/http"

	"github.com/odpf/siren/core/alert"
	"github.com/odpf/siren/core/namespace"
	"github.com/odpf/siren/core/provider"
	"github.com/odpf/siren/core/receiver"
	"github.com/odpf/siren/core/rules"
	"github.com/odpf/siren/core/subscription"
	"github.com/odpf/siren/core/template"
	"github.com/odpf/siren/domain"
	"github.com/odpf/siren/internal/store"
	slackclient "github.com/odpf/siren/plugins/receivers/http"
	"github.com/odpf/siren/plugins/receivers/slack"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type Container struct {
	TemplateService     TemplateService
	RulesService        domain.RuleService
	AlertService        AlertService
	NotifierServices    NotifierServices
	ProviderService     ProviderService
	NamespaceService    NamespaceService
	ReceiverService     ReceiverService
	SubscriptionService domain.SubscriptionService
}

func InitContainer(repositories *store.RepositoryContainer, db *gorm.DB, c *domain.Config, httpClient *http.Client) (*Container, error) {
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
	rulesService := rules.NewService(
		repositories.RuleRepository,
		templateService,
		namespaceService,
		providerService,
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
	subscriptionService, err := subscription.NewService(repositories.SubscriptionRepository, repositories.ProviderRepository,
		repositories.NamespaceRepository, repositories.ReceiverRepository, c.EncryptionKey)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create subscriptions service")
	}

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
