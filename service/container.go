package service

import (
	"github.com/grafana/cortex-tools/pkg/client"
	"github.com/odpf/siren/domain"
	"github.com/odpf/siren/pkg/alert"
	"github.com/odpf/siren/pkg/alert/alertmanager"
	"github.com/odpf/siren/pkg/alert_history"
	"github.com/odpf/siren/pkg/codeexchange"
	"github.com/odpf/siren/pkg/rules"
	"github.com/odpf/siren/pkg/templates"
	"gorm.io/gorm"
	"net/http"
)

type Container struct {
	TemplatesService    domain.TemplatesService
	RulesService        domain.RuleService
	AlertmanagerService domain.AlertmanagerService
	AlertHistoryService domain.AlertHistoryService
	CodeExchangeService domain.CodeExchangeService
}

func Init(db *gorm.DB, cortex domain.CortexConfig, siren domain.SirenServiceConfig, slackAppConfig domain.SlackApp,
	client *client.CortexClient, httpClient *http.Client) (*Container, error) {
	templatesService := templates.NewService(db)
	rulesService := rules.NewService(db, client)
	newClient, err := alertmanager.NewClient(cortex)
	if err != nil {
		return nil, err
	}
	alertmanagerService := alert.NewService(db, newClient, siren)
	alertHistoryService := alert_history.NewService(db)
	codeExchangeService := codeexchange.NewService(db, httpClient, slackAppConfig)
	return &Container{
		TemplatesService:    templatesService,
		RulesService:        rulesService,
		AlertmanagerService: alertmanagerService,
		AlertHistoryService: alertHistoryService,
		CodeExchangeService: codeExchangeService,
	}, nil
}

func (container *Container) MigrateAll(db *gorm.DB) error {
	err := container.TemplatesService.Migrate()
	if err != nil {
		return err
	}
	err = container.AlertmanagerService.Migrate()

	if err != nil {
		return err
	}
	err = container.RulesService.Migrate()
	if err != nil {
		return err
	}
	err = container.AlertHistoryService.Migrate()
	if err != nil {
		return err
	}
	err = container.CodeExchangeService.Migrate()
	if err != nil {
		return err
	}
	return nil
}
