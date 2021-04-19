package service

import (
	"github.com/odpf/siren/domain"
	"github.com/odpf/siren/pkg/alert"
	"github.com/odpf/siren/pkg/alert/alertmanager"
	"github.com/odpf/siren/pkg/rules"
	"github.com/odpf/siren/pkg/templates"
	"github.com/odpf/siren/pkg/alert_history"
	"gorm.io/gorm"
)

type Container struct {
	TemplatesService    domain.TemplatesService
	RulesService        domain.RuleService
	AlertmanagerService domain.AlertmanagerService
	AlertHistoryService domain.AlertHistoryService
}

func Init(db *gorm.DB, cortex domain.CortexConfig) (*Container, error) {
	templatesService := templates.NewService(db)
	rulesService := rules.NewService(db, cortex)
	newClient, err := alertmanager.NewClient(cortex)
	if err != nil {
		return nil, err
	}
	alertmanagerService := alert.NewService(db, newClient)
	alertHistoryService := alert_history.NewService(db)
	return &Container{
		TemplatesService:    templatesService,
		RulesService:        rulesService,
		AlertmanagerService: alertmanagerService,
		AlertHistoryService: alertHistoryService,
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
	return nil
}
