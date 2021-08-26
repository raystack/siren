package v1

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/newrelic/go-agent/v3/newrelic"
	pb "github.com/odpf/siren/api/proto/odpf/siren"
	"github.com/odpf/siren/domain"
	"github.com/odpf/siren/service"
	"github.com/slack-go/slack"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"strings"
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
	if err != nil {
		if strings.Contains(err.Error(), "alert history parameters missing") {
			s.logger.Error(err.Error())
			return nil, status.Errorf(codes.InvalidArgument, err.Error())
		}
		return nil, status.Errorf(codes.Internal, err.Error())
	}
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

func (s *GRPCServer) SendSlackNotification(_ context.Context, req *pb.SendSlackNotificationRequest) (*pb.SendSlackNotificationResponse, error) {
	var payload *domain.SlackMessage
	provider := req.GetProvider()

	b, err := json.Marshal(req.GetBlocks())
	if err != nil {
		s.logger.Error("handler", zap.Error(err))
		return nil, status.Errorf(codes.InvalidArgument, "Invalid block")
	}

	blocks := slack.Blocks{}
	err = json.Unmarshal(b, &blocks)
	if err != nil {
		s.logger.Error("handler", zap.Error(err))
		return nil, status.Errorf(codes.InvalidArgument, "unable to parse block")
	}

	if provider == "slack" {
		payload = &domain.SlackMessage{
			ReceiverName: req.GetReceiverName(),
			ReceiverType: req.GetReceiverType(),
			Entity:       req.GetEntity(),
			Message:      req.GetMessage(),
			Blocks:       blocks,
		}
	} else {
		return nil, status.Errorf(codes.InvalidArgument, "provider not supported")
	}

	result, err := s.container.NotifierServices.Slack.Notify(payload)
	if err != nil {
		s.logger.Error("handler", zap.Error(err))
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	res := &pb.SendSlackNotificationResponse{
		Ok: result.OK,
	}
	return res, nil
}

func (s *GRPCServer) GetRules(_ context.Context, req *pb.GetRulesRequest) (*pb.GetRulesResponse, error) {
	namespace := req.GetNamespace()
	entity := req.GetEntity()
	groupName := req.GetGroupName()
	ruleStatus := req.GetStatus()
	template := req.GetTemplate()

	rules, err := s.container.RulesService.Get(namespace, entity, groupName, ruleStatus, template)
	if err != nil {
		s.logger.Error("handler", zap.Error(err))
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	res := &pb.GetRulesResponse{Rules: make([]*pb.Rule, 0)}
	for _, rule := range rules {

		variables := make([]*pb.Variables, 0)
		for _, variable := range rule.Variables {
			variables = append(variables, &pb.Variables{
				Name:        variable.Name,
				Type:        variable.Type,
				Value:       variable.Value,
				Description: variable.Description,
			})
		}
		res.Rules = append(res.Rules, &pb.Rule{
			Id:        uint64(rule.ID),
			Name:      rule.Name,
			Entity:    rule.Entity,
			Namespace: rule.Namespace,
			GroupName: rule.GroupName,
			Template:  rule.Template,
			Status:    rule.Status,
			CreatedAt: timestamppb.New(rule.CreatedAt),
			UpdatedAt: timestamppb.New(rule.UpdatedAt),
			Variables: variables,
		})
	}

	return res, nil
}

func (s *GRPCServer) UpdateRule(_ context.Context, req *pb.UpdateRuleRequest) (*pb.Rule, error) {
	variables := make([]domain.RuleVariable, 0)
	for _, variable := range req.Variables {
		variables = append(variables, domain.RuleVariable{
			Name:        variable.Name,
			Type:        variable.Type,
			Value:       variable.Value,
			Description: variable.Description,
		})
	}

	payload := &domain.Rule{
		ID:        uint(req.GetId()),
		Name:      req.GetName(),
		Entity:    req.GetEntity(),
		Namespace: req.GetNamespace(),
		GroupName: req.GetGroupName(),
		Template:  req.GetTemplate(),
		Status:    req.GetStatus(),
		Variables: variables,
	}

	rule, err := s.container.RulesService.Upsert(payload)
	if err != nil {
		s.logger.Error("handler", zap.Error(err))
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	responseVariables := make([]*pb.Variables, 0)
	for _, variable := range rule.Variables {
		responseVariables = append(responseVariables, &pb.Variables{
			Name:        variable.Name,
			Type:        variable.Type,
			Value:       variable.Value,
			Description: variable.Description,
		})
	}
	res := &pb.Rule{
		Id:        uint64(rule.ID),
		Name:      rule.Name,
		Entity:    rule.Entity,
		Namespace: rule.Namespace,
		GroupName: rule.GroupName,
		Template:  rule.Template,
		Status:    rule.Status,
		CreatedAt: timestamppb.New(rule.CreatedAt),
		UpdatedAt: timestamppb.New(rule.UpdatedAt),
		Variables: responseVariables,
	}
	return res, nil
}

func (s *GRPCServer) GetTemplates(_ context.Context, req *pb.GetTemplatesRequest) (*pb.GetTemplatesResponse, error) {
	templates, err := s.container.TemplatesService.Index(req.GetTag())
	if err != nil {
		s.logger.Error("handler", zap.Error(err))
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	res := &pb.GetTemplatesResponse{Templates: make([]*pb.Template, 0)}
	for _, template := range templates {

		variables := make([]*pb.TemplateVariables, 0)
		for _, variable := range template.Variables {
			variables = append(variables, &pb.TemplateVariables{
				Name:        variable.Name,
				Type:        variable.Type,
				Default:     variable.Default,
				Description: variable.Description,
			})
		}
		res.Templates = append(res.Templates, &pb.Template{
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

func (s *GRPCServer) GetTemplateByName(_ context.Context, req *pb.GetTemplateByNameRequest) (*pb.Template, error) {
	template, err := s.container.TemplatesService.GetByName(req.GetName())
	if err != nil {
		s.logger.Error("handler", zap.Error(err))
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	variables := make([]*pb.TemplateVariables, 0)
	for _, variable := range template.Variables {
		variables = append(variables, &pb.TemplateVariables{
			Name:        variable.Name,
			Type:        variable.Type,
			Default:     variable.Default,
			Description: variable.Description,
		})
	}
	res := &pb.Template{
		Id:        uint64(template.ID),
		Name:      template.Name,
		Body:      template.Body,
		Tags:      template.Tags,
		CreatedAt: timestamppb.New(template.CreatedAt),
		UpdatedAt: timestamppb.New(template.UpdatedAt),
		Variables: variables,
	}
	return res, nil
}

func (s *GRPCServer) UpsertTemplate(_ context.Context, req *pb.UpsertTemplateRequest) (*pb.Template, error) {
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
		s.logger.Error("handler", zap.Error(err))
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	templateVariables := make([]*pb.TemplateVariables, 0)
	for _, variable := range template.Variables {
		templateVariables = append(templateVariables, &pb.TemplateVariables{
			Name:        variable.Name,
			Type:        variable.Type,
			Default:     variable.Default,
			Description: variable.Description,
		})
	}
	res := &pb.Template{
		Id:        uint64(template.ID),
		Name:      template.Name,
		Body:      template.Body,
		Tags:      template.Tags,
		CreatedAt: timestamppb.New(template.CreatedAt),
		UpdatedAt: timestamppb.New(template.UpdatedAt),
		Variables: templateVariables,
	}
	return res, nil
}

func (s *GRPCServer) DeleteTemplate(_ context.Context, req *pb.DeleteTemplateRequest) (*pb.DeleteTemplateResponse, error) {
	err := s.container.TemplatesService.Delete(req.GetName())
	if err != nil {
		s.logger.Error("handler", zap.Error(err))
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return &pb.DeleteTemplateResponse{}, nil
}

func (s *GRPCServer) RenderTemplate(_ context.Context, req *pb.RenderTemplateRequest) (*pb.RenderTemplateResponse, error) {
	body, err := s.container.TemplatesService.Render(req.GetName(), req.GetVariables())
	if err != nil {
		s.logger.Error("handler", zap.Error(err))
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return &pb.RenderTemplateResponse{
		Body: body,
	}, nil
}
