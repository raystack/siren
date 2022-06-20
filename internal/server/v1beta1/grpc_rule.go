package v1beta1

import (
	"context"

	"github.com/odpf/siren/core/rule"
	"github.com/odpf/siren/utils"
	sirenv1beta1 "go.buf.build/odpf/gw/odpf/proton/odpf/siren/v1beta1"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/types/known/timestamppb"
)

//go:generate mockery --name=RuleService -r --case underscore --with-expecter --structname RuleService --filename rule_service.go --output=./mocks
type RuleService interface {
	Upsert(context.Context, *rule.Rule) error
	Get(context.Context, string, string, string, string, uint64) ([]rule.Rule, error)
	Migrate() error
}

func (s *GRPCServer) ListRules(ctx context.Context, req *sirenv1beta1.ListRulesRequest) (*sirenv1beta1.ListRulesResponse, error) {
	name := req.GetName()
	namespace := req.GetNamespace()
	groupName := req.GetGroupName()
	template := req.GetTemplate()
	providerNamespace := req.GetProviderNamespace()

	rules, err := s.ruleService.Get(ctx, name, namespace, groupName, template, providerNamespace)
	if err != nil {
		return nil, utils.GRPCLogError(s.logger, codes.Internal, err)
	}

	res := &sirenv1beta1.ListRulesResponse{Rules: make([]*sirenv1beta1.Rule, 0)}
	for _, rule := range rules {
		variables := make([]*sirenv1beta1.Variables, 0)
		for _, variable := range rule.Variables {
			variables = append(variables, &sirenv1beta1.Variables{
				Name:        variable.Name,
				Value:       variable.Value,
				Type:        variable.Type,
				Description: variable.Description,
			})
		}
		res.Rules = append(res.Rules, &sirenv1beta1.Rule{
			Id:                rule.Id,
			Name:              rule.Name,
			Enabled:           rule.Enabled,
			GroupName:         rule.GroupName,
			Namespace:         rule.Namespace,
			Template:          rule.Template,
			Variables:         variables,
			ProviderNamespace: rule.ProviderNamespace,
			CreatedAt:         timestamppb.New(rule.CreatedAt),
			UpdatedAt:         timestamppb.New(rule.UpdatedAt),
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
		return nil, utils.GRPCLogError(s.logger, codes.Internal, err)
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
			Id:                rule.Id,
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
