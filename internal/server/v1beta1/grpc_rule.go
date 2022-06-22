package v1beta1

import (
	"context"

	"github.com/odpf/siren/core/rule"
	sirenv1beta1 "github.com/odpf/siren/internal/server/proto/odpf/siren/v1beta1"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/types/known/timestamppb"
)

//go:generate mockery --name=RuleService -r --case underscore --with-expecter --structname RuleService --filename rule_service.go --output=./mocks
type RuleService interface {
	Upsert(context.Context, *rule.Rule) error
	Get(context.Context, string, string, string, string, uint64) ([]rule.Rule, error)
}

func (s *GRPCServer) ListRules(ctx context.Context, req *sirenv1beta1.ListRulesRequest) (*sirenv1beta1.ListRulesResponse, error) {
	name := req.GetName()
	namespace := req.GetNamespace()
	groupName := req.GetGroupName()
	template := req.GetTemplate()
	providerNamespace := req.GetProviderNamespace()

	rules, err := s.ruleService.Get(ctx, name, namespace, groupName, template, providerNamespace)
	if err != nil {
		return nil, gRPCLogError(s.logger, codes.Internal, err)
	}

	res := &sirenv1beta1.ListRulesResponse{Rules: make([]*sirenv1beta1.Rule, 0)}
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
		res.Rules = append(res.Rules, &sirenv1beta1.Rule{
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

	return res, nil
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

	rule := &rule.Rule{
		Enabled:           req.GetEnabled(),
		GroupName:         req.GetGroupName(),
		Namespace:         req.GetNamespace(),
		Template:          req.GetTemplate(),
		ProviderNamespace: req.GetProviderNamespace(),
		Variables:         variables,
	}

	if err := s.ruleService.Upsert(ctx, rule); err != nil {
		return nil, gRPCLogError(s.logger, codes.Internal, err)
	}

	responseVariables := make([]*sirenv1beta1.Variables, 0)
	for _, variable := range rule.Variables {
		responseVariables = append(responseVariables, &sirenv1beta1.Variables{
			Name:        variable.Name,
			Type:        variable.Type,
			Value:       variable.Value,
			Description: variable.Description,
		})
	}
	res := &sirenv1beta1.UpdateRuleResponse{
		Rule: &sirenv1beta1.Rule{
			Id:                rule.ID,
			Name:              rule.Name,
			Enabled:           rule.Enabled,
			GroupName:         rule.GroupName,
			Namespace:         rule.Namespace,
			Template:          rule.Template,
			Variables:         responseVariables,
			ProviderNamespace: rule.ProviderNamespace,
			CreatedAt:         timestamppb.New(rule.CreatedAt),
			UpdatedAt:         timestamppb.New(rule.UpdatedAt),
		},
	}
	return res, nil
}