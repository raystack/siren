package base

import (
	"context"
	"html/template"

	"github.com/raystack/siren/core/provider"
	"github.com/raystack/siren/core/rule"
	"github.com/raystack/siren/plugins"
)

// UnimplementedService is a base receiver provider service layer
type UnimplementedService struct{}

func (s *UnimplementedService) SyncRuntimeConfig(ctx context.Context, namespaceID uint64, namespaceURN string, prov provider.Provider) error {
	return nil
}

func (s *UnimplementedService) UpsertRule(ctx context.Context, namespaceURN string, prov provider.Provider, rl *rule.Rule, templateToUpdate *template.Template) error {
	return plugins.ErrNotImplemented
}
