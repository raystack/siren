package v1beta1

import (
	"context"

	"github.com/goto/siren/core/template"
	sirenv1beta1 "github.com/goto/siren/proto/gotocompany/siren/v1beta1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *GRPCServer) ListTemplates(ctx context.Context, req *sirenv1beta1.ListTemplatesRequest) (*sirenv1beta1.ListTemplatesResponse, error) {
	templates, err := s.templateService.List(ctx, template.Filter{
		Tag: req.GetTag(),
	})
	if err != nil {
		return nil, s.generateRPCErr(err)
	}

	items := []*sirenv1beta1.Template{}
	for _, tmpl := range templates {
		variables := make([]*sirenv1beta1.TemplateVariables, 0)
		for _, variable := range tmpl.Variables {
			variables = append(variables, &sirenv1beta1.TemplateVariables{
				Name:        variable.Name,
				Type:        variable.Type,
				Default:     variable.Default,
				Description: variable.Description,
			})
		}

		items = append(items, &sirenv1beta1.Template{
			Id:        uint64(tmpl.ID),
			Name:      tmpl.Name,
			Body:      tmpl.Body,
			Tags:      tmpl.Tags,
			CreatedAt: timestamppb.New(tmpl.CreatedAt),
			UpdatedAt: timestamppb.New(tmpl.UpdatedAt),
			Variables: variables,
		})
	}

	return &sirenv1beta1.ListTemplatesResponse{
		Templates: items,
	}, nil
}

func (s *GRPCServer) GetTemplate(ctx context.Context, req *sirenv1beta1.GetTemplateRequest) (*sirenv1beta1.GetTemplateResponse, error) {
	template, err := s.templateService.GetByName(ctx, req.GetName())
	if err != nil {
		return nil, s.generateRPCErr(err)
	}

	variables := make([]*sirenv1beta1.TemplateVariables, 0)
	for _, variable := range template.Variables {
		variables = append(variables, &sirenv1beta1.TemplateVariables{
			Name:        variable.Name,
			Type:        variable.Type,
			Default:     variable.Default,
			Description: variable.Description,
		})
	}
	return &sirenv1beta1.GetTemplateResponse{
		Template: &sirenv1beta1.Template{
			Id:        uint64(template.ID),
			Name:      template.Name,
			Body:      template.Body,
			Tags:      template.Tags,
			CreatedAt: timestamppb.New(template.CreatedAt),
			UpdatedAt: timestamppb.New(template.UpdatedAt),
			Variables: variables,
		},
	}, nil
}

func (s *GRPCServer) UpsertTemplate(ctx context.Context, req *sirenv1beta1.UpsertTemplateRequest) (*sirenv1beta1.UpsertTemplateResponse, error) {
	variables := make([]template.Variable, 0)
	for _, variable := range req.GetVariables() {
		variables = append(variables, template.Variable{
			Name:        variable.Name,
			Type:        variable.Type,
			Default:     variable.Default,
			Description: variable.Description,
		})
	}

	tmpl := &template.Template{
		ID:        req.GetId(),
		Name:      req.GetName(),
		Body:      req.GetBody(),
		Tags:      req.GetTags(),
		Variables: variables,
	}

	if err := s.templateService.Upsert(ctx, tmpl); err != nil {
		return nil, s.generateRPCErr(err)
	}

	return &sirenv1beta1.UpsertTemplateResponse{
		Id: tmpl.ID,
	}, nil
}

func (s *GRPCServer) DeleteTemplate(ctx context.Context, req *sirenv1beta1.DeleteTemplateRequest) (*sirenv1beta1.DeleteTemplateResponse, error) {
	err := s.templateService.Delete(ctx, req.GetName())
	if err != nil {
		return nil, s.generateRPCErr(err)
	}
	return &sirenv1beta1.DeleteTemplateResponse{}, nil
}

func (s *GRPCServer) RenderTemplate(ctx context.Context, req *sirenv1beta1.RenderTemplateRequest) (*sirenv1beta1.RenderTemplateResponse, error) {
	body, err := s.templateService.Render(ctx, req.GetName(), req.GetVariables())
	if err != nil {
		return nil, s.generateRPCErr(err)
	}
	return &sirenv1beta1.RenderTemplateResponse{
		Body: body,
	}, nil
}
