package service

import (
	"github.com/odpf/siren/domain"
	"github.com/odpf/siren/pkg/alert"
	"github.com/odpf/siren/pkg/alert/alertmanager"
	"github.com/odpf/siren/pkg/rules"
	"github.com/odpf/siren/pkg/templates"
	"gorm.io/gorm"
)

type Container struct {
	TemplatesService    domain.TemplatesService
	RulesService        domain.RuleService
	AlertmanagerService domain.AlertmanagerService
}

func Init(db *gorm.DB, cortex domain.CortexConfig) (*Container, error) {
	templatesService := templates.NewService(db)
	rulesService := rules.NewService(db, cortex)
	newClient, err := alertmanager.NewClient(cortex)
	if err != nil {
		return nil, err
	}
	alertmanagerService := alert.NewService(db, newClient)
	return &Container{
		TemplatesService:    templatesService,
		RulesService:        rulesService,
		AlertmanagerService: alertmanagerService,
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
	return nil
}
