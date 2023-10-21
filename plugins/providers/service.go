package providers

import (
	"context"
	"html/template"

	"github.com/goto/siren/core/alert"
	"github.com/goto/siren/core/namespace"
	"github.com/goto/siren/core/provider"
	"github.com/goto/siren/core/rule"
	"github.com/goto/siren/plugins"
)

// UnimplementedService is a base receiver provider service layer
type UnimplementedService struct{}

func (s *UnimplementedService) SyncRuntimeConfig(ctx context.Context, namespaceID uint64, namespaceURN string, prov provider.Provider) error {
	return nil
}

func (s *UnimplementedService) UpsertRule(ctx context.Context, ns namespace.Namespace, prov provider.Provider, rl *rule.Rule, templateToUpdate *template.Template) error {
	return plugins.ErrNotImplemented
}

func (s *UnimplementedService) SetConfig(ctx context.Context, configRaw string) error {
	return plugins.ErrNotImplemented
}

func (s *UnimplementedService) TransformToAlerts(ctx context.Context, providerID uint64, namespaceID uint64, body map[string]any) ([]alert.Alert, int, error) {
	return nil, 0, plugins.ErrNotImplemented
}
