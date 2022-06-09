package v1beta1

import (
	"context"

	"github.com/odpf/siren/core/template"
	"github.com/odpf/siren/utils"
	sirenv1beta1 "go.buf.build/odpf/gw/odpf/proton/odpf/siren/v1beta1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

//go:generate mockery --name=TemplateService -r --case underscore --with-expecter --structname TemplateService --filename template_service.go --output=./mocks
type TemplateService interface {
	Upsert(*template.Template) error
	Index(string) ([]template.Template, error)
	GetByName(string) (*template.Template, error)
	Delete(string) error
	Render(string, map[string]string) (string, error)
	Migrate() error
}

func (s *GRPCServer) ListTemplates(_ context.Context, req *sirenv1beta1.ListTemplatesRequest) (*sirenv1beta1.ListTemplatesResponse, error) {
	templates, err := s.container.TemplateService.Index(req.GetTag())
	if err != nil {
		return nil, utils.GRPCLogError(s.logger, codes.Internal, err)
	}

	res := &sirenv1beta1.ListTemplatesResponse{Templates: make([]*sirenv1beta1.Template, 0)}
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
		res.Templates = append(res.Templates, &sirenv1beta1.Template{
			Id:        uint64(tmpl.ID),
			Name:      tmpl.Name,
			Body:      tmpl.Body,
			Tags:      tmpl.Tags,
			CreatedAt: timestamppb.New(tmpl.CreatedAt),
			UpdatedAt: timestamppb.New(tmpl.UpdatedAt),
			Variables: variables,
		})
	}

	return res, nil
}

func (s *GRPCServer) GetTemplateByName(_ context.Context, req *sirenv1beta1.GetTemplateByNameRequest) (*sirenv1beta1.TemplateResponse, error) {
	template, err := s.container.TemplateService.GetByName(req.GetName())
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
	variables := make([]template.Variable, 0)
	for _, variable := range req.GetVariables() {
		variables = append(variables, template.Variable{
			Name:        variable.Name,
			Type:        variable.Type,
			Default:     variable.Default,
			Description: variable.Description,
		})
	}
	template := &template.Template{
		ID:        uint(req.GetId()),
		Name:      req.GetName(),
		Body:      req.GetBody(),
		Tags:      req.GetTags(),
		Variables: variables,
	}
	err := s.container.TemplateService.Upsert(template)
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
	err := s.container.TemplateService.Delete(req.GetName())
	if err != nil {
		return nil, utils.GRPCLogError(s.logger, codes.Internal, err)
	}
	return &sirenv1beta1.DeleteTemplateResponse{}, nil
}

func (s *GRPCServer) RenderTemplate(_ context.Context, req *sirenv1beta1.RenderTemplateRequest) (*sirenv1beta1.RenderTemplateResponse, error) {
	body, err := s.container.TemplateService.Render(req.GetName(), req.GetVariables())
	if err != nil {
		return nil, utils.GRPCLogError(s.logger, codes.Internal, err)
	}
	return &sirenv1beta1.RenderTemplateResponse{
		Body: body,
	}, nil
}
