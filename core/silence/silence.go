package silence

import (
	"context"
	"fmt"
	"time"

	"github.com/antonmedv/expr"
)

const TargetExpressionRuleKey = "rule"

//go:generate mockery --name=Repository -r --case underscore --with-expecter --structname SubscriptionRepository --filename subscription_repository.go --output=./mocks
type Repository interface {
	Create(context.Context, Silence) (string, error)
	List(context.Context, Filter) ([]Silence, error)
	Get(ctx context.Context, id string) (Silence, error)
	SoftDelete(ctx context.Context, id string) error
}

type Silence struct {
	ID               string         `json:"id"`
	NamespaceID      uint64         `json:"namespace_id"`
	Type             string         `json:"type"`
	TargetID         uint64         `json:"target_id"`
	TargetExpression map[string]any `json:"target_expression"`
	Creator          string         `json:"creator"`
	Comment          string         `json:"comment"`
	CreatedAt        time.Time      `json:"created_at"`
	DeletedAt        time.Time      `json:"deleted_at"`
}

func (s Silence) Validate() error {
	switch s.Type {
	case TypeSubscription:
		if s.TargetID == 0 {
			return fmt.Errorf("target id cannot be empty or zero for type '%s'", TypeSubscription)
		}
	case TypeMatchers:
		if len(s.TargetExpression) == 0 {
			return fmt.Errorf("target expression cannot be empty and should be kv labels for type '%s'", TypeMatchers)
		}
	default:
		return fmt.Errorf("unknown silence type '%s', should be '%s' or '%s'", s.Type, TypeMatchers, TypeSubscription)
	}
	return nil
}

func (s Silence) subscriptionRule() (string, error) {
	if s.Type != TypeSubscription {
		return "", fmt.Errorf("silence id '%s' type is not subscription, type is '%s' instead", s.ID, s.Type)
	}

	rule, ok := s.TargetExpression[TargetExpressionRuleKey]
	if !ok {
		return "", nil
	}

	ruleStr := fmt.Sprintf("%s", rule)

	return ruleStr, nil
}

func (s Silence) EvaluateSubscriptionRule(env interface{}) (bool, error) {
	rule, err := s.subscriptionRule()
	if err != nil {
		return false, err
	}

	if rule == "" {
		return true, nil
	}

	res, err := expr.Eval(rule, env)
	if err != nil {
		return false, err
	}

	resBool, ok := res.(bool)
	if !ok {
		return false, fmt.Errorf("rule evaluation result is not boolean: %v", res)
	}

	return resBool, nil
}
