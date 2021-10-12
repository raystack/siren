package v1

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/newrelic/go-agent/v3/newrelic"
	sirenv1 "github.com/odpf/siren/api/proto/odpf/siren/v1"
	"github.com/odpf/siren/domain"
	"github.com/odpf/siren/service"
	"github.com/slack-go/slack"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"strings"
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

func (s *GRPCServer) ListProviders(_ context.Context, _ *emptypb.Empty) (*sirenv1.ListProvidersResponse, error) {
	providers, err := s.container.ProviderService.ListProviders()
	if err != nil {
		s.logger.Error("handler", zap.Error(err))
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	res := &sirenv1.ListProvidersResponse{
		Providers: make([]*sirenv1.Provider, 0),
	}
	for _, provider := range providers {
		credentials, err := structpb.NewStruct(provider.Credentials)
		if err != nil {
			s.logger.Error("handler", zap.Error(err))
			return nil, status.Errorf(codes.Internal, err.Error())
		}

		item := &sirenv1.Provider{
			Id:          provider.Id,
			Host:        provider.Host,
			Type:        provider.Type,
			Name:        provider.Name,
			Credentials: credentials,
			Labels:      provider.Labels,
			CreatedAt:   timestamppb.New(provider.CreatedAt),
			UpdatedAt:   timestamppb.New(provider.UpdatedAt),
		}
		res.Providers = append(res.Providers, item)
	}
	return res, nil
}

func (s *GRPCServer) CreateProvider(_ context.Context, req *sirenv1.CreateProviderRequest) (*sirenv1.Provider, error) {
	provider, err := s.container.ProviderService.CreateProvider(&domain.Provider{
		Host:        req.GetHost(),
		Name:        req.GetName(),
		Type:        req.GetType(),
		Credentials: req.GetCredentials().AsMap(),
		Labels:      req.GetLabels(),
	})
	if err != nil {
		s.logger.Error("handler", zap.Error(err))
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	grpcCredentials, err := structpb.NewStruct(provider.Credentials)
	if err != nil {
		s.logger.Error("handler", zap.Error(err))
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &sirenv1.Provider{
		Id:          provider.Id,
		Host:        provider.Host,
		Name:        provider.Name,
		Type:        provider.Type,
		Credentials: grpcCredentials,
		Labels:      provider.Labels,
		CreatedAt:   timestamppb.New(provider.CreatedAt),
		UpdatedAt:   timestamppb.New(provider.UpdatedAt),
	}, nil
}

func (s *GRPCServer) GetProvider(_ context.Context, req *sirenv1.GetProviderRequest) (*sirenv1.Provider, error) {
	provider, err := s.container.ProviderService.GetProvider(req.GetId())
	if provider == nil {
		return nil, status.Errorf(codes.NotFound, "provider not found")
	}
	if err != nil {
		s.logger.Error("handler", zap.Error(err))
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	grpcCredentials, err := structpb.NewStruct(provider.Credentials)
	if err != nil {
		s.logger.Error("handler", zap.Error(err))
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &sirenv1.Provider{
		Id:          provider.Id,
		Host:        provider.Host,
		Name:        provider.Name,
		Type:        provider.Type,
		Credentials: grpcCredentials,
		Labels:      provider.Labels,
		CreatedAt:   timestamppb.New(provider.CreatedAt),
		UpdatedAt:   timestamppb.New(provider.UpdatedAt),
	}, nil
}

func (s *GRPCServer) UpdateProvider(_ context.Context, req *sirenv1.UpdateProviderRequest) (*sirenv1.Provider, error) {
	provider, err := s.container.ProviderService.UpdateProvider(&domain.Provider{
		Id:          req.GetId(),
		Host:        req.GetHost(),
		Name:        req.GetName(),
		Type:        req.GetType(),
		Credentials: req.GetCredentials().AsMap(),
		Labels:      req.GetLabels(),
	})
	if err != nil {
		s.logger.Error("handler", zap.Error(err))
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	grpcCredentials, err := structpb.NewStruct(provider.Credentials)
	if err != nil {
		s.logger.Error("handler", zap.Error(err))
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &sirenv1.Provider{
		Id:          provider.Id,
		Host:        provider.Host,
		Name:        provider.Name,
		Type:        provider.Type,
		Credentials: grpcCredentials,
		Labels:      provider.Labels,
		CreatedAt:   timestamppb.New(provider.CreatedAt),
		UpdatedAt:   timestamppb.New(provider.UpdatedAt),
	}, nil
}

func (s *GRPCServer) DeleteProvider(_ context.Context, req *sirenv1.DeleteProviderRequest) (*emptypb.Empty, error) {
	err := s.container.ProviderService.DeleteProvider(uint64(req.GetId()))
	if err != nil {
		s.logger.Error("handler", zap.Error(err))
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	
	return &emptypb.Empty{}, nil
}

func (s *GRPCServer) ListAlertHistory(_ context.Context, req *sirenv1.ListAlertHistoryRequest) (*sirenv1.ListAlertHistoryResponse, error) {
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
	res := &sirenv1.ListAlertHistoryResponse{
		Alerts: make([]*sirenv1.AlertHistory, 0),
	}
	for _, alert := range alerts {
		item := &sirenv1.AlertHistory{
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

func (s *GRPCServer) CreateAlertHistory(_ context.Context, req *sirenv1.CreateAlertHistoryRequest) (*sirenv1.CreateAlertHistoryResponse, error) {
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
	result := &sirenv1.CreateAlertHistoryResponse{Alerts: make([]*sirenv1.AlertHistory, 0)}
	for _, item := range createdAlerts {
		alertHistoryItem := &sirenv1.AlertHistory{
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
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return result, nil
}

func (s *GRPCServer) ListWorkspaceChannels(_ context.Context, req *sirenv1.ListWorkspaceChannelsRequest) (*sirenv1.ListWorkspaceChannelsResponse, error) {
	workspace := req.GetWorkspaceName()
	workspaces, err := s.container.SlackWorkspaceService.GetChannels(workspace)
	if err != nil {
		s.logger.Error("handler", zap.Error(err))
		return nil, status.Errorf(codes.Internal, err.Error())
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
		s.logger.Error("handler", zap.Error(err))
		return nil, status.Errorf(codes.Internal, err.Error())
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
		s.logger.Error("handler", zap.Error(err))
		return nil, status.Errorf(codes.Internal, err.Error())
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
		s.logger.Error("handler", zap.Error(err))
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return &sirenv1.UpdateAlertCredentialsResponse{}, nil
}

func (s *GRPCServer) SendSlackNotification(_ context.Context, req *sirenv1.SendSlackNotificationRequest) (*sirenv1.SendSlackNotificationResponse, error) {
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
	res := &sirenv1.SendSlackNotificationResponse{
		Ok: result.OK,
	}
	return res, nil
}

func (s *GRPCServer) ListRules(_ context.Context, req *sirenv1.ListRulesRequest) (*sirenv1.ListRulesResponse, error) {
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

	res := &sirenv1.ListRulesResponse{Rules: make([]*sirenv1.Rule, 0)}
	for _, rule := range rules {
		variables := make([]*sirenv1.Variables, 0)
		for _, variable := range rule.Variables {
			variables = append(variables, &sirenv1.Variables{
				Name:        variable.Name,
				Type:        variable.Type,
				Value:       variable.Value,
				Description: variable.Description,
			})
		}
		res.Rules = append(res.Rules, &sirenv1.Rule{
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

func (s *GRPCServer) UpdateRule(_ context.Context, req *sirenv1.UpdateRuleRequest) (*sirenv1.UpdateRuleResponse, error) {
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

	responseVariables := make([]*sirenv1.Variables, 0)
	for _, variable := range rule.Variables {
		responseVariables = append(responseVariables, &sirenv1.Variables{
			Name:        variable.Name,
			Type:        variable.Type,
			Value:       variable.Value,
			Description: variable.Description,
		})
	}
	res := &sirenv1.UpdateRuleResponse{
		Rule: &sirenv1.Rule{
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
		},
	}
	return res, nil
}

func (s *GRPCServer) ListTemplates(_ context.Context, req *sirenv1.ListTemplatesRequest) (*sirenv1.ListTemplatesResponse, error) {
	templates, err := s.container.TemplatesService.Index(req.GetTag())
	if err != nil {
		s.logger.Error("handler", zap.Error(err))
		return nil, status.Errorf(codes.Internal, err.Error())
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
		s.logger.Error("handler", zap.Error(err))
		return nil, status.Errorf(codes.Internal, err.Error())
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
		s.logger.Error("handler", zap.Error(err))
		return nil, status.Errorf(codes.Internal, err.Error())
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
		s.logger.Error("handler", zap.Error(err))
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return &sirenv1.DeleteTemplateResponse{}, nil
}

func (s *GRPCServer) RenderTemplate(_ context.Context, req *sirenv1.RenderTemplateRequest) (*sirenv1.RenderTemplateResponse, error) {
	body, err := s.container.TemplatesService.Render(req.GetName(), req.GetVariables())
	if err != nil {
		s.logger.Error("handler", zap.Error(err))
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return &sirenv1.RenderTemplateResponse{
		Body: body,
	}, nil
}
