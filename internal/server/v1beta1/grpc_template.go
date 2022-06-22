package v1beta1

import (
	"context"

	"github.com/odpf/siren/core/template"
	sirenv1beta1 "github.com/odpf/siren/internal/server/proto/odpf/siren/v1beta1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

//go:generate mockery --name=TemplateService -r --case underscore --with-expecter --structname TemplateService --filename template_service.go --output=./mocks
type TemplateService interface {
	Upsert(*template.Template) error
	Index(string) ([]template.Template, error)
	GetByName(string) (*template.Template, error)
	Delete(string) error
	Render(string, map[string]string) (string, error)
}

func (s *GRPCServer) ListTemplates(_ context.Context, req *sirenv1beta1.ListTemplatesRequest) (*sirenv1beta1.ListTemplatesResponse, error) {
	templates, err := s.templateService.Index(req.GetTag())
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

func (s *GRPCServer) GetTemplate(_ context.Context, req *sirenv1beta1.GetTemplateRequest) (*sirenv1beta1.GetTemplateResponse, error) {
	template, err := s.templateService.GetByName(req.GetName())
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

func (s *GRPCServer) UpsertTemplate(_ context.Context, req *sirenv1beta1.UpsertTemplateRequest) (*sirenv1beta1.UpsertTemplateResponse, error) {
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
	err := s.templateService.Upsert(template)
	if err != nil {
		return nil, s.generateRPCErr(err)
	}

	return &sirenv1beta1.UpsertTemplateResponse{
		Id: uint64(template.ID),
	}, nil
}

func (s *GRPCServer) DeleteTemplate(_ context.Context, req *sirenv1beta1.DeleteTemplateRequest) (*sirenv1beta1.DeleteTemplateResponse, error) {
	err := s.templateService.Delete(req.GetName())
	if err != nil {
		return nil, s.generateRPCErr(err)
	}
	return &sirenv1beta1.DeleteTemplateResponse{}, nil
}

func (s *GRPCServer) RenderTemplate(_ context.Context, req *sirenv1beta1.RenderTemplateRequest) (*sirenv1beta1.RenderTemplateResponse, error) {
	body, err := s.templateService.Render(req.GetName(), req.GetVariables())
	if err != nil {
		return nil, s.generateRPCErr(err)
	}
	return &sirenv1beta1.RenderTemplateResponse{
		Body: body,
	}, nil
}
