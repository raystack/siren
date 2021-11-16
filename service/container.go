package service

import (
	"github.com/odpf/siren/pkg/namespace"
	"github.com/odpf/siren/pkg/provider"
	"github.com/odpf/siren/pkg/receiver"
	"net/http"

	"github.com/grafana/cortex-tools/pkg/client"
	"github.com/odpf/siren/domain"
	"github.com/odpf/siren/pkg/alerts"
	"github.com/odpf/siren/pkg/codeexchange"
	"github.com/odpf/siren/pkg/rules"
	"github.com/odpf/siren/pkg/slacknotifier"
	"github.com/odpf/siren/pkg/templates"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type Container struct {
	TemplatesService    domain.TemplatesService
	RulesService        domain.RuleService
	AlertService        domain.AlertService
	CodeExchangeService domain.CodeExchangeService
	NotifierServices    domain.NotifierServices
	ProviderService     domain.ProviderService
	NamespaceService    domain.NamespaceService
	ReceiverService     domain.ReceiverService
}

func Init(db *gorm.DB, c *domain.Config,
	client *client.CortexClient, httpClient *http.Client) (*Container, error) {
	templatesService := templates.NewService(db)
	rulesService := rules.NewService(db, client)
	alertHistoryService := alerts.NewService(db)
	codeExchangeService, err := codeexchange.NewService(db, httpClient, c.SlackApp, c.EncryptionKey)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create codeexchange service")
	}

	slackNotifierService := slacknotifier.NewService()
	providerService := provider.NewService(db)
	namespaceService, err := namespace.NewService(db, c.EncryptionKey)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create namespace service")
	}
	receiverService, err := receiver.NewService(db, httpClient, c.EncryptionKey)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create receiver service")
	}

	return &Container{
		TemplatesService:    templatesService,
		RulesService:        rulesService,
		AlertService:        alertHistoryService,
		CodeExchangeService: codeExchangeService,
		NotifierServices: domain.NotifierServices{
			Slack: slackNotifierService,
		},
		ProviderService:  providerService,
		NamespaceService: namespaceService,
		ReceiverService:  receiverService,
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
	err = container.CodeExchangeService.Migrate()
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
	return nil
}
