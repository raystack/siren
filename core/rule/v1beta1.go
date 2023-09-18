package rule

import (
	sirenv1beta1 "github.com/goto/siren/proto/gotocompany/siren/v1beta1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (r *Rule) ToV1beta1Proto() *sirenv1beta1.Rule {
	variables := make([]*sirenv1beta1.Variables, 0)
	for _, variable := range r.Variables {
		variables = append(variables, &sirenv1beta1.Variables{
			Name:        variable.Name,
			Value:       variable.Value,
			Type:        variable.Type,
			Description: variable.Description,
		})
	}
	return &sirenv1beta1.Rule{
		Id:                r.ID,
		Name:              r.Name,
		Enabled:           r.Enabled,
		GroupName:         r.GroupName,
		Namespace:         r.Namespace,
		Template:          r.Template,
		Variables:         variables,
		ProviderNamespace: r.ProviderNamespace,
		CreatedAt:         timestamppb.New(r.CreatedAt),
		UpdatedAt:         timestamppb.New(r.UpdatedAt),
	}
}

func FromV1beta1Proto(proto *sirenv1beta1.Rule) *Rule {
	variables := make([]RuleVariable, 0)
	for _, variable := range proto.GetVariables() {
		variables = append(variables, RuleVariable{
			Name:        variable.Name,
			Type:        variable.Type,
			Value:       variable.Value,
			Description: variable.Description,
		})
	}

	return &Rule{
		ID:                proto.GetId(),
		Name:              proto.GetName(),
		Enabled:           proto.GetEnabled(),
		GroupName:         proto.GetGroupName(),
		Namespace:         proto.GetNamespace(),
		Template:          proto.GetTemplate(),
		Variables:         variables,
		ProviderNamespace: proto.GetProviderNamespace(),
		CreatedAt:         proto.GetCreatedAt().AsTime(),
		UpdatedAt:         proto.GetUpdatedAt().AsTime(),
	}
}
