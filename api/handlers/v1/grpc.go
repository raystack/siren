package v1

import (
	"context"
	"github.com/newrelic/go-agent/v3/newrelic"
	sirenv1 "github.com/odpf/siren/api/proto/odpf/siren/v1"
	"github.com/odpf/siren/domain"
	"github.com/odpf/siren/helper"
	"github.com/odpf/siren/service"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type GRPCServer struct {
	container *service.Container
	newrelic  *newrelic.Application
	logger    *zap.Logger
	sirenv1.UnimplementedSirenServiceServer
}

func NewGRPCServer(container *service.Container, nr *newrelic.Application, logger *zap.Logger) *GRPCServer {
	return &GRPCServer{
		container: container,
		newrelic:  nr,
		logger:    logger,
	}
}

func (s *GRPCServer) Ping(ctx context.Context, in *sirenv1.PingRequest) (*sirenv1.PingResponse, error) {
	return &sirenv1.PingResponse{Message: "Pong"}, nil
}

func (s *GRPCServer) ListWorkspaceChannels(_ context.Context, req *sirenv1.ListWorkspaceChannelsRequest) (*sirenv1.ListWorkspaceChannelsResponse, error) {
	workspace := req.GetWorkspaceName()
	workspaces, err := s.container.SlackWorkspaceService.GetChannels(workspace)
	if err != nil {
		return nil, helper.GRPCLogError(s.logger, codes.Internal, err)
	}
	res := &sirenv1.ListWorkspaceChannelsResponse{
		Data: make([]*sirenv1.SlackWorkspace, 0),
	}
	for _, workspace := range workspaces {
		item := &sirenv1.SlackWorkspace{
			Id:   workspace.ID,
			Name: workspace.Name,
		}
		res.Data = append(res.Data, item)
	}
	return res, nil
}

func (s *GRPCServer) ExchangeCode(_ context.Context, req *sirenv1.ExchangeCodeRequest) (*sirenv1.ExchangeCodeResponse, error) {
	code := req.GetCode()
	workspace := req.GetWorkspace()
	result, err := s.container.CodeExchangeService.Exchange(domain.OAuthPayload{
		Code:      code,
		Workspace: workspace,
	})
	if err != nil {
		return nil, helper.GRPCLogError(s.logger, codes.Internal, err)
	}
	res := &sirenv1.ExchangeCodeResponse{
		Ok: result.OK,
	}
	return res, nil
}

func (s *GRPCServer) GetAlertCredentials(_ context.Context, req *sirenv1.GetAlertCredentialsRequest) (*sirenv1.GetAlertCredentialsResponse, error) {
	teamName := req.GetTeamName()
	alertCredential, err := s.container.AlertmanagerService.Get(teamName)
	if err != nil {
		return nil, helper.GRPCLogError(s.logger, codes.Internal, err)
	}
	res := &sirenv1.GetAlertCredentialsResponse{
		Entity:               alertCredential.Entity,
		TeamName:             alertCredential.TeamName,
		PagerdutyCredentials: alertCredential.PagerdutyCredentials,
		SlackConfig: &sirenv1.SlackConfig{
			Critical: &sirenv1.Critical{Channel: alertCredential.SlackConfig.Critical.Channel},
			Warning:  &sirenv1.Warning{Channel: alertCredential.SlackConfig.Warning.Channel},
		},
	}
	return res, nil
}

func (s *GRPCServer) UpdateAlertCredentials(_ context.Context, req *sirenv1.UpdateAlertCredentialsRequest) (*sirenv1.UpdateAlertCredentialsResponse, error) {
	entity := req.GetEntity()
	teamName := req.GetTeamName()
	pagerdutyCredential := req.GetPagerdutyCredentials()
	criticalChannel := req.GetSlackConfig().GetCritical().GetChannel()
	warningChannel := req.GetSlackConfig().GetWarning().GetChannel()

	if entity == "" {
		return nil, status.Errorf(codes.InvalidArgument, "entity cannot be empty")
	}
	if pagerdutyCredential == "" {
		return nil, status.Errorf(codes.InvalidArgument, "pagerduty credential cannot be empty")
	}

	payload := domain.AlertCredential{
		Entity:               entity,
		TeamName:             teamName,
		PagerdutyCredentials: pagerdutyCredential,
		SlackConfig: domain.SlackConfig{
			Critical: domain.SlackCredential{
				Channel: criticalChannel,
			},
			Warning: domain.SlackCredential{
				Channel: warningChannel,
			},
		},
	}

	err := s.container.AlertmanagerService.Upsert(payload)
	if err != nil {
		return nil, helper.GRPCLogError(s.logger, codes.Internal, err)
	}
	return &sirenv1.UpdateAlertCredentialsResponse{}, nil
}

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
