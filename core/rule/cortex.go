package rule

import (
	"context"

	"github.com/odpf/siren/pkg/cortex"
)

//go:generate mockery --name=CortexClient -r --case underscore --with-expecter --structname CortexClient --filename cortex_client.go --output=./mocks
type CortexClient interface {
	CreateAlertmanagerConfig(cortex.AlertManagerConfig, string) error
	CreateRuleGroup(ctx context.Context, namespace string, rg cortex.RuleGroup) error
	DeleteRuleGroup(ctx context.Context, namespace, groupName string) error
	GetRuleGroup(ctx context.Context, namespace, groupName string) (*cortex.RuleGroup, error)
	ListRules(ctx context.Context, namespace string) (map[string][]cortex.RuleGroup, error)
}
