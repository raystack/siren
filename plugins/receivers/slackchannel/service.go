package slackchannel

import (
	"context"
	"fmt"

	"github.com/goto/siren/core/notification"
	"github.com/goto/siren/pkg/errors"
	"github.com/goto/siren/plugins/receivers/base"
	"github.com/goto/siren/plugins/receivers/slack"
	"github.com/mitchellh/mapstructure"
)

// PluginService is a plugin service layer for slack channel type
type PluginService struct {
	base.UnimplementedService
	slackPluginService *slack.PluginService
}

// NewPluginService returns slack channel plugin service struct. This service implement [receiver.Resolver] and [notification.Notifier] interface.
func NewPluginService(cfg slack.AppConfig, cryptoClient slack.Encryptor, opts ...slack.ServiceOption) *PluginService {
	return &PluginService{
		slackPluginService: slack.NewPluginService(cfg, cryptoClient, opts...),
	}
}

func (s *PluginService) PreHookDBTransformConfigs(ctx context.Context, configurations map[string]any) (map[string]any, error) {
	receiverConfig := &ReceiverConfig{}
	if err := mapstructure.Decode(configurations, receiverConfig); err != nil {
		return nil, fmt.Errorf("failed to transform configurations to receiver config: %w", err)
	}

	if err := receiverConfig.Validate(); err != nil {
		return nil, errors.ErrInvalid.WithMsgf(err.Error())
	}

	return configurations, nil
}

// PostHookTransformConfigs do transformation in post-hook service lifecycle
func (s *PluginService) PostHookDBTransformConfigs(ctx context.Context, configurations map[string]any) (map[string]any, error) {
	transformedConfigs, err := s.slackPluginService.PostHookDBTransformConfigs(ctx, configurations)
	// if slack_channel is not expaneded, it is okay to have slack config empty
	if err != nil && !errors.Is(err, errors.ErrInvalid) {
		return nil, fmt.Errorf("slack channel post hook db failed: %w", err)
	}

	var mergedConfigs = map[string]any{}
	for k, v := range configurations {
		if value, ok := transformedConfigs[k]; ok {
			mergedConfigs[k] = value
		} else {
			mergedConfigs[k] = v
		}
	}

	return mergedConfigs, nil
}

func (s *PluginService) PreHookQueueTransformConfigs(ctx context.Context, notificationConfigMap map[string]any) (map[string]any, error) {
	transformedConfigs, err := s.slackPluginService.PreHookQueueTransformConfigs(ctx, notificationConfigMap)
	if err != nil {
		return nil, err
	}

	var mergedConfigs = map[string]any{}
	for k, v := range notificationConfigMap {
		if value, ok := transformedConfigs[k]; ok {
			mergedConfigs[k] = value
		} else {
			mergedConfigs[k] = v
		}
	}

	return mergedConfigs, nil
}

func (s *PluginService) PostHookQueueTransformConfigs(ctx context.Context, notificationConfigMap map[string]any) (map[string]any, error) {
	transformedConfigs, err := s.slackPluginService.PostHookQueueTransformConfigs(ctx, notificationConfigMap)
	if err != nil {
		return nil, err
	}

	var mergedConfigs = map[string]any{}
	for k, v := range notificationConfigMap {
		if value, ok := transformedConfigs[k]; ok {
			mergedConfigs[k] = value
		} else {
			mergedConfigs[k] = v
		}
	}

	return mergedConfigs, nil
}

func (s *PluginService) Send(ctx context.Context, notificationMessage notification.Message) (bool, error) {
	return s.slackPluginService.Send(ctx, notificationMessage)
}

func (s *PluginService) GetSystemDefaultTemplate() string {
	return s.slackPluginService.GetSystemDefaultTemplate()
}
