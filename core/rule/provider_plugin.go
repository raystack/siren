package rule

import (
	"context"

	"github.com/goto/siren/core/namespace"
	"github.com/goto/siren/core/provider"
	"github.com/goto/siren/core/template"
)

// RuleUploader is an interface for the provider to upload rule(s).
// Provider plugin needs to implement this interface in order to
// support rule synchronization from siren to provider
//

type RuleUploader interface {
	UpsertRule(ctx context.Context, ns namespace.Namespace, prov provider.Provider, rl *Rule, templateToUpdate *template.Template) error
}
