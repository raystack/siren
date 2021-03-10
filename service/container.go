package service

import (
	"github.com/odpf/siren/domain"
	"github.com/odpf/siren/rules"
	"github.com/odpf/siren/templates"
	"gorm.io/gorm"
)

type Container struct {
	TemplatesService domain.TemplatesService
	RulesService     domain.RuleService
}

func Init(db *gorm.DB) *Container {
	templatesService := templates.NewService(db)
	rulesService := rules.NewService(db)
	return &Container{
		TemplatesService: templatesService,
		RulesService:     rulesService,
	}
}

func MigrateAll(db *gorm.DB) error {
	container := Init(db)
	err := container.TemplatesService.Migrate()
	if err != nil {
		return err
	}
	err = container.RulesService.Migrate()
	if err != nil {
		return err
	}
	return nil
}
