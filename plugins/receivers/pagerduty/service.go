package pagerduty

import (
	"context"
	"fmt"

	"github.com/mitchellh/mapstructure"
	"github.com/raystack/siren/core/notification"
	"github.com/raystack/siren/pkg/errors"
	"github.com/raystack/siren/pkg/httpclient"
	"github.com/raystack/siren/pkg/retry"
	"github.com/raystack/siren/plugins/receivers/base"
)

type PluginService struct {
	base.UnimplementedService
	client     PagerDutyCaller
	httpClient *httpclient.Client
	retrier    retry.Runner
}

func NewPluginService(cfg AppConfig, opts ...ServiceOption) *PluginService {
	s := &PluginService{}

	for _, opt := range opts {
		opt(s)
	}

	if s.client == nil {
		s.client = NewClient(cfg, ClientWithHTTPClient(s.httpClient), ClientWithRetrier(s.retrier))
	}

	return s
}

// TODO validation could be done by default and using validator library but we need to use generic for it
func (s *PluginService) PreHookDBTransformConfigs(ctx context.Context, configurations map[string]interface{}) (map[string]interface{}, error) {
	receiverConfig := &ReceiverConfig{}
	if err := mapstructure.Decode(configurations, receiverConfig); err != nil {
		return nil, fmt.Errorf("failed to transform configurations to receiver config: %w", err)
	}

	if err := receiverConfig.Validate(); err != nil {
		return nil, errors.ErrInvalid.WithMsgf(err.Error())
	}

	return configurations, nil
}

func (s *PluginService) PreHookQueueTransformConfigs(ctx context.Context, notificationConfigMap map[string]interface{}) (map[string]interface{}, error) {
	notificationConfig := &NotificationConfig{}
	if err := mapstructure.Decode(notificationConfigMap, notificationConfig); err != nil {
		return nil, fmt.Errorf("failed to transform configurations to slack notification config: %w", err)
	}

	if err := notificationConfig.Validate(); err != nil {
		return nil, err
	}

	return notificationConfig.AsMap(), nil
}

func (s *PluginService) Send(ctx context.Context, notificationMessage notification.Message) (bool, error) {
	notificationConfig := &NotificationConfig{}
	if err := mapstructure.Decode(notificationMessage.Configs, notificationConfig); err != nil {
		return false, err
	}

	pgMessageV1 := &MessageV1{}
	if err := mapstructure.Decode(notificationMessage.Details, &pgMessageV1); err != nil {
		return false, err
	}
	pgMessageV1.ServiceKey = notificationConfig.ServiceKey

	if err := s.client.NotifyV1(ctx, *pgMessageV1); err != nil {
		if errors.As(err, new(retry.RetryableError)) {
			return true, err
		} else {
			return false, err
		}
	}

	return false, nil
}

func (s *PluginService) GetSystemDefaultTemplate() string {
	return defaultAlertTemplateBodyV1
}
