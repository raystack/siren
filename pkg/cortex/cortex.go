package cortex

import (
	"context"

	"github.com/grafana/cortex-tools/pkg/rules/rwrulefmt"
)

//go:generate mockery --name=CortexCaller -r --case underscore --with-expecter --structname CortexCaller --filename cortex_caller.go --output=./mocks
type CortexCaller interface {
	CreateAlertmanagerConfig(ctx context.Context, cfg string, templates map[string]string) error
	CreateRuleGroup(ctx context.Context, namespace string, rg rwrulefmt.RuleGroup) error
	DeleteRuleGroup(ctx context.Context, namespace, groupName string) error
	GetRuleGroup(ctx context.Context, namespace, groupName string) (*rwrulefmt.RuleGroup, error)
	ListRules(ctx context.Context, namespace string) (map[string][]rwrulefmt.RuleGroup, error)
}

type ReceiverConfig struct {
	Receiver      string
	Type          string
	Match         map[string]string
	Configuration map[string]string
}

type AlertManagerConfig struct {
	Receivers []ReceiverConfig
}

type RuleGroup struct {
	rwrulefmt.RuleGroup
}
