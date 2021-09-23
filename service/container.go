package service

import (
	"github.com/odpf/siren/pkg/slackworkspace"
	"github.com/odpf/siren/pkg/workspace"
	"net/http"

	"github.com/grafana/cortex-tools/pkg/client"
	"github.com/odpf/siren/domain"
	"github.com/odpf/siren/pkg/alert"
	"github.com/odpf/siren/pkg/alert/alertmanager"
	"github.com/odpf/siren/pkg/alert_history"
	"github.com/odpf/siren/pkg/codeexchange"
	"github.com/odpf/siren/pkg/rules"
	"github.com/odpf/siren/pkg/slacknotifier"
	"github.com/odpf/siren/pkg/templates"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type Container struct {
	TemplatesService      domain.TemplatesService
	RulesService          domain.RuleService
	AlertmanagerService   domain.AlertmanagerService
	AlertHistoryService   domain.AlertHistoryService
	CodeExchangeService   domain.CodeExchangeService
	NotifierServices      domain.NotifierServices
	SlackWorkspaceService domain.SlackWorkspaceService
	WorkspaceService      domain.WorkspaceService
}

func Init(db *gorm.DB, c *domain.Config,
	client *client.CortexClient, httpClient *http.Client) (*Container, error) {
	templatesService := templates.NewService(db)
	rulesService := rules.NewService(db, client)
	newClient, err := alertmanager.NewClient(c.Cortex)
	if err != nil {
		return nil, err
	}
	alertHistoryService := alert_history.NewService(db)
	codeExchangeService, err := codeexchange.NewService(db, httpClient, c.SlackApp, c.EncryptionKey)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create codeexchange service")
	}
	alertmanagerService := alert.NewService(db, newClient, c.SirenService, codeExchangeService)
	slackNotifierService := slacknotifier.NewService(codeExchangeService)
	slackworkspaceService := slackworkspace.NewService(codeExchangeService)
	workspaceService := workspace.NewService(db)
	return &Container{
		TemplatesService:    templatesService,
		RulesService:        rulesService,
		AlertmanagerService: alertmanagerService,
		AlertHistoryService: alertHistoryService,
		CodeExchangeService: codeExchangeService,
		NotifierServices: domain.NotifierServices{
			Slack: slackNotifierService,
		},
		SlackWorkspaceService: slackworkspaceService,
		WorkspaceService:      workspaceService,
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
	err = container.WorkspaceService.Migrate()
	if err != nil {
		return err
	}
	return nil
}
