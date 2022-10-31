package base

import (
	"context"
	"html/template"

	"github.com/odpf/siren/core/provider"
	"github.com/odpf/siren/core/rule"
	"github.com/odpf/siren/plugins"
)

// UnimplementedService is a base receiver plugin service layer
type UnimplementedService struct{}

func (s *UnimplementedService) SyncRuntimeConfig(ctx context.Context, namespaceURN string, prov provider.Provider) error {
	return nil
}

func (s *UnimplementedService) UpsertRule(ctx context.Context, namespaceURN string, prov provider.Provider, rl *rule.Rule, templateToUpdate *template.Template) error {
	return plugins.ErrNotImplemented
}
