package v1beta1

import (
	"context"
	sirenv1beta1 "github.com/odpf/siren/api/proto/odpf/siren/v1beta1"
	"github.com/odpf/siren/domain"
	"github.com/odpf/siren/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *GRPCServer) ListTemplates(_ context.Context, req *sirenv1beta1.ListTemplatesRequest) (*sirenv1beta1.ListTemplatesResponse, error) {
	templates, err := s.container.TemplatesService.Index(req.GetTag())
	if err != nil {
		return nil, utils.GRPCLogError(s.logger, codes.Internal, err)
	}

	res := &sirenv1beta1.ListTemplatesResponse{Templates: make([]*sirenv1beta1.Template, 0)}
	for _, template := range templates {
		variables := make([]*sirenv1beta1.TemplateVariables, 0)
		for _, variable := range template.Variables {
			variables = append(variables, &sirenv1beta1.TemplateVariables{
				Name:        variable.Name,
				Type:        variable.Type,
				Default:     variable.Default,
				Description: variable.Description,
			})
		}
		res.Templates = append(res.Templates, &sirenv1beta1.Template{
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

func (s *GRPCServer) GetTemplateByName(_ context.Context, req *sirenv1beta1.GetTemplateByNameRequest) (*sirenv1beta1.TemplateResponse, error) {
	template, err := s.container.TemplatesService.GetByName(req.GetName())
	if err != nil {
		return nil, utils.GRPCLogError(s.logger, codes.Internal, err)
	}
	if template == nil {
		return nil, status.Errorf(codes.NotFound, "template not found")
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
	res := &sirenv1beta1.TemplateResponse{
		Template: &sirenv1beta1.Template{
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

func (s *GRPCServer) UpsertTemplate(_ context.Context, req *sirenv1beta1.UpsertTemplateRequest) (*sirenv1beta1.TemplateResponse, error) {
	variables := make([]domain.Variable, 0)
	for _, variable := range req.GetVariables() {
		variables = append(variables, domain.Variable{
			Name:        variable.Name,
			Type:        variable.Type,
			Default:     variable.Default,
			Description: variable.Description,
		})
	}
	template := &domain.Template{
		ID:        uint(req.GetId()),
		Name:      req.GetName(),
		Body:      req.GetBody(),
		Tags:      req.GetTags(),
		Variables: variables,
	}
	err := s.container.TemplatesService.Upsert(template)
	if err != nil {
		return nil, utils.GRPCLogError(s.logger, codes.Internal, err)
	}

	templateVariables := make([]*sirenv1beta1.TemplateVariables, 0)
	for _, variable := range template.Variables {
		templateVariables = append(templateVariables, &sirenv1beta1.TemplateVariables{
			Name:        variable.Name,
			Type:        variable.Type,
			Default:     variable.Default,
			Description: variable.Description,
		})
	}
	res := &sirenv1beta1.TemplateResponse{
		Template: &sirenv1beta1.Template{
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

func (s *GRPCServer) DeleteTemplate(_ context.Context, req *sirenv1beta1.DeleteTemplateRequest) (*sirenv1beta1.DeleteTemplateResponse, error) {
	err := s.container.TemplatesService.Delete(req.GetName())
	if err != nil {
		return nil, utils.GRPCLogError(s.logger, codes.Internal, err)
	}
	return &sirenv1beta1.DeleteTemplateResponse{}, nil
}

func (s *GRPCServer) RenderTemplate(_ context.Context, req *sirenv1beta1.RenderTemplateRequest) (*sirenv1beta1.RenderTemplateResponse, error) {
	body, err := s.container.TemplatesService.Render(req.GetName(), req.GetVariables())
	if err != nil {
		return nil, utils.GRPCLogError(s.logger, codes.Internal, err)
	}
	return &sirenv1beta1.RenderTemplateResponse{
		Body: body,
	}, nil
}
