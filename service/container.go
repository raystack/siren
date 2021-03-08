package service

import (
	"github.com/odpf/siren/domain"
	"github.com/odpf/siren/templates"
	"gorm.io/gorm"
)

type Container struct {
	TemplatesService domain.TemplatesService
}

func Init(db *gorm.DB) *Container {
	templatesService := templates.NewService(db)
	return &Container{
		TemplatesService: templatesService,
	}
}

func MigrateAll(db *gorm.DB) error {
	container := Init(db)
	err := container.TemplatesService.Migrate()
	if err != nil {
		return err
	}
	return nil
}
