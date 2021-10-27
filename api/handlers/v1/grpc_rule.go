package v1

import (
	"context"
	sirenv1 "github.com/odpf/siren/api/proto/odpf/siren/v1"
	"github.com/odpf/siren/domain"
	"github.com/odpf/siren/helper"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *GRPCServer) ListRules(_ context.Context, req *sirenv1.ListRulesRequest) (*sirenv1.ListRulesResponse, error) {
	name := req.GetName()
	namespace := req.GetNamespace()
	groupName := req.GetGroupName()
	template := req.GetTemplate()
	providerNamespace := req.GetProviderNamespace()

	rules, err := s.container.RulesService.Get(name, namespace, groupName, template, providerNamespace)
	if err != nil {
		return nil, helper.GRPCLogError(s.logger, codes.Internal, err)
	}

	res := &sirenv1.ListRulesResponse{Rules: make([]*sirenv1.Rule, 0)}
	for _, rule := range rules {
		variables := make([]*sirenv1.Variables, 0)
		for _, variable := range rule.Variables {
			variables = append(variables, &sirenv1.Variables{
				Name:        variable.Name,
				Value:       variable.Value,
				Type:        variable.Type,
				Description: variable.Description,
			})
		}
		res.Rules = append(res.Rules, &sirenv1.Rule{
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

func (s *GRPCServer) UpdateRule(_ context.Context, req *sirenv1.UpdateRuleRequest) (*sirenv1.UpdateRuleResponse, error) {
	variables := make([]domain.RuleVariable, 0)
	for _, variable := range req.Variables {
		variables = append(variables, domain.RuleVariable{
			Name:        variable.Name,
			Type:        variable.Type,
			Value:       variable.Value,
			Description: variable.Description,
		})
	}

	payload := &domain.Rule{
		Enabled:           req.GetEnabled(),
		GroupName:         req.GetGroupName(),
		Namespace:         req.GetNamespace(),
		Template:          req.GetTemplate(),
		ProviderNamespace: req.GetProviderNamespace(),
		Variables:         variables,
	}

	rule, err := s.container.RulesService.Upsert(payload)
	if err != nil {
		return nil, helper.GRPCLogError(s.logger, codes.Internal, err)
	}

	responseVariables := make([]*sirenv1.Variables, 0)
	for _, variable := range rule.Variables {
		responseVariables = append(responseVariables, &sirenv1.Variables{
			Name:        variable.Name,
			Type:        variable.Type,
			Value:       variable.Value,
			Description: variable.Description,
		})
	}
	res := &sirenv1.UpdateRuleResponse{
		Rule: &sirenv1.Rule{
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
