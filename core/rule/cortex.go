package rule

import (
	"context"

	rwrulefmt "github.com/grafana/cortex-tools/pkg/rules/rwrulefmt"
	"github.com/odpf/siren/pkg/cortex"
)

//go:generate mockery --name=CortexClient -r --case underscore --with-expecter --structname CortexClient --filename cortex_client.go --output=./mocks
type CortexClient interface {
	CreateAlertmanagerConfig(cortex.AlertManagerConfig, string) error
	CreateRuleGroup(ctx context.Context, namespace string, rg rwrulefmt.RuleGroup) error
	DeleteRuleGroup(ctx context.Context, namespace, groupName string) error
	GetRuleGroup(ctx context.Context, namespace, groupName string) (*rwrulefmt.RuleGroup, error)
	ListRules(ctx context.Context, namespace string) (map[string][]rwrulefmt.RuleGroup, error)
}
