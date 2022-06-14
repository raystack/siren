package v1beta1

import (
	"fmt"
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
	"github.com/odpf/siren/pkg/secret"
	"github.com/odpf/siren/pkg/slack"
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

	encryptor, err := secret.New(c.EncryptionKey)
	if err != nil {
		return nil, fmt.Errorf("cannot initialize encryptor: %w", err)
	}

	slackClient := slack.NewClient(slack.ClientWithHTTPClient(&http.Client{}))
	receiverSecureService := receiver.NewSecureService(encryptor, repositories.ReceiverRepository)
	receiverService := receiver.NewService(receiverSecureService, slackClient)

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
			Slack: slackClient,
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
