package v1

import (
	"context"
	"fmt"
	"strings"

	"github.com/newrelic/go-agent/v3/newrelic"
	pb "github.com/odpf/siren/api/proto/odpf/siren"
	"github.com/odpf/siren/domain"
	"github.com/odpf/siren/service"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GRPCServer struct {
	container *service.Container
	newrelic  *newrelic.Application
	logger    *zap.Logger
	pb.UnimplementedSirenServiceServer
}

func NewGRPCServer(container *service.Container, nr *newrelic.Application, logger *zap.Logger) *GRPCServer {
	return &GRPCServer{
		container: container,
		newrelic:  nr,
		logger:    logger,
	}
}

func (s *GRPCServer) Ping(ctx context.Context, in *pb.PingRequest) (*pb.PingResponse, error) {
	return &pb.PingResponse{Message: "Pong"}, nil
}

func (s *GRPCServer) GetAlertHistory(_ context.Context, req *pb.GetAlertHistoryRequest) (*pb.GetAlertHistoryResponse, error) {
	name := req.GetResource()
	startTime := req.GetStartTime()
	endTime := req.GetEndTime()
	if name == "" {
		return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("resource name cannot be empty"))
	}
	alerts, err := s.container.AlertHistoryService.Get(name, startTime, endTime)
	if err != nil {
		s.logger.Error("handler", zap.Error(err))
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	res := &pb.GetAlertHistoryResponse{
		Alerts: make([]*pb.AlertHistory, 0),
	}
	for _, alert := range alerts {
		item := &pb.AlertHistory{
			Name:        alert.Name,
			Id:          alert.ID,
			MetricName:  alert.MetricName,
			MetricValue: alert.MetricValue,
			TemplateId:  alert.TemplateID,
			Level:       alert.Level,
			CreatedAt:   alert.CreatedAt.String(),
			UpdatedAt:   alert.UpdatedAt.String(),
		}
		res.Alerts = append(res.Alerts, item)
	}
	return res, nil
}

func (s *GRPCServer) CreateAlertHistory(_ context.Context, req *pb.CreateAlertHistoryRequest) (*pb.CreateAlertHistoryResponse, error) {
	alerts := domain.Alerts{Alerts: make([]domain.Alert, 0)}
	for _, item := range req.GetAlerts() {
		labels := domain.Labels{
			Severity: item.Labels.Severity,
		}
		annotations := domain.Annotations{
			Resource:    item.GetAnnotations().GetResource(),
			Template:    item.GetAnnotations().GetTemplate(),
			MetricName:  item.GetAnnotations().GetMetricName(),
			MetricValue: item.GetAnnotations().GetMetricValue(),
		}
		alert := domain.Alert{
			Labels:      labels,
			Annotations: annotations,
			Status:      item.Status,
		}
		alerts.Alerts = append(alerts.Alerts, alert)
	}
	createdAlerts, err := s.container.AlertHistoryService.Create(&alerts)
	result := &pb.CreateAlertHistoryResponse{Alerts: make([]*pb.AlertHistory, 0)}
	for _, item := range createdAlerts {
		alertHistoryItem := &pb.AlertHistory{
			Name:        item.Name,
			Id:          item.ID,
			MetricName:  item.MetricName,
			MetricValue: item.MetricValue,
			TemplateId:  item.TemplateID,
			Level:       item.Level,
			CreatedAt:   item.CreatedAt.String(),
			UpdatedAt:   item.UpdatedAt.String(),
		}
		result.Alerts = append(result.Alerts, alertHistoryItem)
	}
	if err != nil {
		if strings.Contains(err.Error(), "alert history parameters missing") {
			s.logger.Error(err.Error())
			return result, nil
		}
		return nil, err
	}
	return result, nil
}

func (s *GRPCServer) GetWorkspaceChannels(_ context.Context, req *pb.GetWorkspaceChannelsRequest) (*pb.GetWorkspaceChannelsResponse, error) {
	workspace := req.GetWorkspaceName()
	workspaces, err := s.container.WorkspaceService.GetChannels(workspace)
	if err != nil {
		s.logger.Error("handler", zap.Error(err))
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	res := &pb.GetWorkspaceChannelsResponse{
		Data: make([]*pb.Workspace, 0),
	}
	for _, workspace := range workspaces {
		item := &pb.Workspace{
			Id:   workspace.ID,
			Name: workspace.Name,
		}
		res.Data = append(res.Data, item)
	}
	return res, nil
}

func (s *GRPCServer) ExchangeCode(_ context.Context, req *pb.ExchangeCodeRequest) (*pb.ExchangeCodeResponse, error) {
	code := req.GetCode()
	workspace := req.GetWorkspace()
	result, err := s.container.CodeExchangeService.Exchange(domain.OAuthPayload{
		Code:      code,
		Workspace: workspace,
	})
	if err != nil {
		s.logger.Error("handler", zap.Error(err))
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	res := &pb.ExchangeCodeResponse{
		Ok: result.OK,
	}
	return res, nil
}

func (s *GRPCServer) GetAlertCredentials(_ context.Context, req *pb.GetAlertCredentialsRequest) (*pb.GetAlertCredentialsResponse, error) {
	teamName := req.GetTeamName()
	alertCredential, err := s.container.AlertmanagerService.Get(teamName)
	if err != nil {
		s.logger.Error("handler", zap.Error(err))
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	res := &pb.GetAlertCredentialsResponse{
		Entity:               alertCredential.Entity,
		TeamName:             alertCredential.TeamName,
		PagerdutyCredentials: alertCredential.PagerdutyCredentials,
		SlackConfig: &pb.SlackConfig{
			Critical: &pb.Critical{Channel: alertCredential.SlackConfig.Critical.Channel},
			Warning:  &pb.Warning{Channel: alertCredential.SlackConfig.Warning.Channel},
		},
	}
	return res, nil
}

func (s *GRPCServer) UpdateAlertCredentials(_ context.Context, req *pb.UpdateAlertCredentialsRequest) (*pb.UpdateAlertCredentialsResponse, error) {
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
		s.logger.Error("handler", zap.Error(err))
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return &pb.UpdateAlertCredentialsResponse{}, nil
}
