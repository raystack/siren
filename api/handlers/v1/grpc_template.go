package v1

import (
	"context"
	sirenv1 "github.com/odpf/siren/api/proto/odpf/siren/v1"
	"github.com/odpf/siren/domain"
	"github.com/odpf/siren/helper"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *GRPCServer) ListTemplates(_ context.Context, req *sirenv1.ListTemplatesRequest) (*sirenv1.ListTemplatesResponse, error) {
	templates, err := s.container.TemplatesService.Index(req.GetTag())
	if err != nil {
		return nil, helper.GRPCLogError(s.logger, codes.Internal, err)
	}

	res := &sirenv1.ListTemplatesResponse{Templates: make([]*sirenv1.Template, 0)}
	for _, template := range templates {
		variables := make([]*sirenv1.TemplateVariables, 0)
		for _, variable := range template.Variables {
			variables = append(variables, &sirenv1.TemplateVariables{
				Name:        variable.Name,
				Type:        variable.Type,
				Default:     variable.Default,
				Description: variable.Description,
			})
		}
		res.Templates = append(res.Templates, &sirenv1.Template{
			Id:        uint64(template.ID),
			Name:      template.Name,
			Body:      template.Body,
			Tags:      template.Tags,
			CreatedAt: timestamppb.New(template.CreatedAt),
			UpdatedAt: timestamppb.New(template.UpdatedAt),
			Variables: variables,
		})
	}

	return res, nil
}

func (s *GRPCServer) GetTemplateByName(_ context.Context, req *sirenv1.GetTemplateByNameRequest) (*sirenv1.TemplateResponse, error) {
	template, err := s.container.TemplatesService.GetByName(req.GetName())
	if err != nil {
		return nil, helper.GRPCLogError(s.logger, codes.Internal, err)
	}

	variables := make([]*sirenv1.TemplateVariables, 0)
	for _, variable := range template.Variables {
		variables = append(variables, &sirenv1.TemplateVariables{
			Name:        variable.Name,
			Type:        variable.Type,
			Default:     variable.Default,
			Description: variable.Description,
		})
	}
	res := &sirenv1.TemplateResponse{
		Template: &sirenv1.Template{
			Id:        uint64(template.ID),
			Name:      template.Name,
			Body:      template.Body,
			Tags:      template.Tags,
			CreatedAt: timestamppb.New(template.CreatedAt),
			UpdatedAt: timestamppb.New(template.UpdatedAt),
			Variables: variables,
		},
	}
	return res, nil
}

func (s *GRPCServer) UpsertTemplate(_ context.Context, req *sirenv1.UpsertTemplateRequest) (*sirenv1.TemplateResponse, error) {
	variables := make([]domain.Variable, 0)
	for _, variable := range req.GetVariables() {
		variables = append(variables, domain.Variable{
			Name:        variable.Name,
			Type:        variable.Type,
			Default:     variable.Default,
			Description: variable.Description,
		})
	}
	payload := &domain.Template{
		ID:        uint(req.GetId()),
		Name:      req.GetName(),
		Body:      req.GetBody(),
		Tags:      req.GetTags(),
		Variables: variables,
	}
	template, err := s.container.TemplatesService.Upsert(payload)
	if err != nil {
		return nil, helper.GRPCLogError(s.logger, codes.Internal, err)
	}

	templateVariables := make([]*sirenv1.TemplateVariables, 0)
	for _, variable := range template.Variables {
		templateVariables = append(templateVariables, &sirenv1.TemplateVariables{
			Name:        variable.Name,
			Type:        variable.Type,
			Default:     variable.Default,
			Description: variable.Description,
		})
	}
	res := &sirenv1.TemplateResponse{
		Template: &sirenv1.Template{
			Id:        uint64(template.ID),
			Name:      template.Name,
			Body:      template.Body,
			Tags:      template.Tags,
			CreatedAt: timestamppb.New(template.CreatedAt),
			UpdatedAt: timestamppb.New(template.UpdatedAt),
			Variables: templateVariables,
		},
	}
	return res, nil
}

func (s *GRPCServer) DeleteTemplate(_ context.Context, req *sirenv1.DeleteTemplateRequest) (*sirenv1.DeleteTemplateResponse, error) {
	err := s.container.TemplatesService.Delete(req.GetName())
	if err != nil {
		return nil, helper.GRPCLogError(s.logger, codes.Internal, err)
	}
	return &sirenv1.DeleteTemplateResponse{}, nil
}

func (s *GRPCServer) RenderTemplate(_ context.Context, req *sirenv1.RenderTemplateRequest) (*sirenv1.RenderTemplateResponse, error) {
	body, err := s.container.TemplatesService.Render(req.GetName(), req.GetVariables())
	if err != nil {
		return nil, helper.GRPCLogError(s.logger, codes.Internal, err)
	}
	return &sirenv1.RenderTemplateResponse{
		Body: body,
	}, nil
}
