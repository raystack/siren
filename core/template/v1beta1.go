package template

import (
	sirenv1beta1 "github.com/goto/siren/proto/gotocompany/siren/v1beta1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (t *Template) ToV1beta1Proto() *sirenv1beta1.Template {
	variables := make([]*sirenv1beta1.TemplateVariables, 0)
	for _, variable := range t.Variables {
		variables = append(variables, &sirenv1beta1.TemplateVariables{
			Name:        variable.Name,
			Type:        variable.Type,
			Default:     variable.Default,
			Description: variable.Description,
		})
	}

	return &sirenv1beta1.Template{
		Id:        uint64(t.ID),
		Name:      t.Name,
		Body:      t.Body,
		Tags:      t.Tags,
		CreatedAt: timestamppb.New(t.CreatedAt),
		UpdatedAt: timestamppb.New(t.UpdatedAt),
		Variables: variables,
	}
}

func FromV1beta1Proto(proto *sirenv1beta1.Template) *Template {
	variables := make([]Variable, 0)
	for _, variable := range proto.GetVariables() {
		variables = append(variables, Variable{
			Name:        variable.Name,
			Type:        variable.Type,
			Default:     variable.Default,
			Description: variable.Description,
		})
	}

	return &Template{
		ID:        proto.GetId(),
		Name:      proto.GetName(),
		Body:      proto.GetBody(),
		Tags:      proto.GetTags(),
		Variables: variables,
		CreatedAt: proto.GetCreatedAt().AsTime(),
		UpdatedAt: proto.GetUpdatedAt().AsTime(),
	}
}
