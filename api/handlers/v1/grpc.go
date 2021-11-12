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

