package base

import (
	"context"
	"html/template"

	"github.com/goto/siren/core/provider"
	"github.com/goto/siren/core/rule"
	"github.com/goto/siren/plugins"
)

// UnimplementedService is a base receiver provider service layer
type UnimplementedService struct{}

func (s *UnimplementedService) SyncRuntimeConfig(ctx context.Context, namespaceID uint64, namespaceURN string, prov provider.Provider) error {
	return nil
}

func (s *UnimplementedService) UpsertRule(ctx context.Context, namespaceURN string, prov provider.Provider, rl *rule.Rule, templateToUpdate *template.Template) error {
	return plugins.ErrNotImplemented
}
