package service

import "github.com/odpf/siren/store"

type Container struct {
	TemplatesService TemplatesService
}

func Init(templatesStore *store.TemplatesStore) *Container {
	templatesService := NewTemplatesService(templatesStore)
	return &Container{
		TemplatesService: templatesService,
	}
}
