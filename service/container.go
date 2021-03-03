package service

import (
	"github.com/odpf/siren/templates"
	"gorm.io/gorm"
)

type Container struct {
	TemplatesService *templates.Service
}

func Init(db *gorm.DB) *Container {
	templatesService := templates.NewService(db)
	return &Container{
		TemplatesService: templatesService,
	}
}
