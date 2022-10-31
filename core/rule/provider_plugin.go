package rule

import (
	"context"

	"github.com/odpf/siren/core/provider"
	"github.com/odpf/siren/core/template"
)

// RuleUploader is an interface for the provider to upload rule(s).
// Provider plugin needs to implement this interface in order to
// support rule synchronization from siren to provider
//
//go:generate mockery --name=RuleUploader -r --case underscore --with-expecter --structname RuleUploader --filename rule_uploader.go --output=./mocks
type RuleUploader interface {
	UpsertRule(ctx context.Context, namespaceURN string, prov provider.Provider, rl *Rule, templateToUpdate *template.Template) error
}
