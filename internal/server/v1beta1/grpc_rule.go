package v1beta1

import (
	"context"

	"github.com/odpf/siren/core/rule"
	sirenv1beta1 "github.com/odpf/siren/internal/server/proto/odpf/siren/v1beta1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

//go:generate mockery --name=RuleService -r --case underscore --with-expecter --structname RuleService --filename rule_service.go --output=./mocks
type RuleService interface {
	Upsert(context.Context, *rule.Rule) (uint64, error)
	List(context.Context, rule.Filter) ([]rule.Rule, error)
}

func (s *GRPCServer) ListRules(ctx context.Context, req *sirenv1beta1.ListRulesRequest) (*sirenv1beta1.ListRulesResponse, error) {
	rules, err := s.ruleService.List(ctx, rule.Filter{
		Name:         req.GetName(),
		Namespace:    req.GetNamespace(),
		GroupName:    req.GetGroupName(),
		TemplateName: req.GetTemplate(),
		NamespaceID:  req.GetProviderNamespace(),
	})
	if err != nil {
		return nil, s.generateRPCErr(err)
	}

	rulesProto := []*sirenv1beta1.Rule{}

	for _, rl := range rules {
		variables := make([]*sirenv1beta1.Variables, 0)
		for _, variable := range rl.Variables {
			variables = append(variables, &sirenv1beta1.Variables{
				Name:        variable.Name,
				Value:       variable.Value,
				Type:        variable.Type,
				Description: variable.Description,
			})
		}
		rulesProto = append(rulesProto, &sirenv1beta1.Rule{
			Id:                rl.ID,
			Name:              rl.Name,
			Enabled:           rl.Enabled,
			GroupName:         rl.GroupName,
			Namespace:         rl.Namespace,
			Template:          rl.Template,
			Variables:         variables,
			ProviderNamespace: rl.ProviderNamespace,
			CreatedAt:         timestamppb.New(rl.CreatedAt),
			UpdatedAt:         timestamppb.New(rl.UpdatedAt),
		})
	}

	return &sirenv1beta1.ListRulesResponse{
		Rules: rulesProto,
	}, nil
}

func (s *GRPCServer) UpdateRule(ctx context.Context, req *sirenv1beta1.UpdateRuleRequest) (*sirenv1beta1.UpdateRuleResponse, error) {
	variables := make([]rule.RuleVariable, 0)
	for _, variable := range req.Variables {
		variables = append(variables, rule.RuleVariable{
			Name:        variable.Name,
			Type:        variable.Type,
			Value:       variable.Value,
			Description: variable.Description,
		})
	}

	rl := &rule.Rule{
		Enabled:           req.GetEnabled(),
		GroupName:         req.GetGroupName(),
		Namespace:         req.GetNamespace(),
		Template:          req.GetTemplate(),
		ProviderNamespace: req.GetProviderNamespace(),
		Variables:         variables,
	}

	id, err := s.ruleService.Upsert(ctx, rl)
	if err != nil {
		return nil, s.generateRPCErr(err)
	}

	return &sirenv1beta1.UpdateRuleResponse{
		Id: id,
	}, nil
}
