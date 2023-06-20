package slack

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/goto/siren/core/notification"
	"github.com/goto/siren/pkg/errors"
	"github.com/goto/siren/pkg/httpclient"
	"github.com/goto/siren/pkg/retry"
	"github.com/goto/siren/plugins/receivers/base"
	"github.com/mitchellh/mapstructure"
)

const (
	TypeChannelChannel = "channel"
	TypeChannelUser    = "user"

	defaultChannelType = TypeChannelChannel
)

// PluginService is a plugin service layer for slack
type PluginService struct {
	base.UnimplementedService
	client       SlackCaller
	cryptoClient Encryptor
	httpClient   *httpclient.Client
	retrier      retry.Runner
}

// NewPluginService returns slack plugin service struct. This service implement [receiver.Resolver] and [notification.Notifier] interface.
func NewPluginService(cfg AppConfig, cryptoClient Encryptor, opts ...ServiceOption) *PluginService {
	s := &PluginService{}

	for _, opt := range opts {
		opt(s)
	}

	s.cryptoClient = cryptoClient

	if s.httpClient == nil {
		s.httpClient = httpclient.New(cfg.HTTPClient)
	}

	if s.client == nil {
		s.client = NewClient(cfg, ClientWithHTTPClient(s.httpClient), ClientWithRetrier(s.retrier))
	}

	return s
}

func (s *PluginService) PreHookDBTransformConfigs(ctx context.Context, configurations map[string]any) (map[string]any, error) {
	slackCredentialConfig := &SlackCredentialConfig{}
	if err := mapstructure.Decode(configurations, slackCredentialConfig); err != nil {
		return nil, fmt.Errorf("failed to transform configurations to pre transform config: %w", err)
	}

	if err := slackCredentialConfig.Validate(); err != nil {
		return nil, errors.ErrInvalid.WithMsgf(err.Error())
	}

	creds, err := s.client.ExchangeAuth(ctx, slackCredentialConfig.AuthCode, slackCredentialConfig.ClientID, slackCredentialConfig.ClientSecret)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code with slack OAuth server: %w", err)
	}

	cipherText, err := s.cryptoClient.Encrypt(creds.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("slack token encryption failed: %w", err)
	}

	receiverConfig := ReceiverConfig{
		Workspace: creds.TeamName,
		Token:     cipherText,
	}

	return receiverConfig.AsMap(), nil
}

// PostHookTransformConfigs do transformation in post-hook service lifecycle
func (s *PluginService) PostHookDBTransformConfigs(ctx context.Context, configurations map[string]any) (map[string]any, error) {
	receiverConfig := &ReceiverConfig{}
	if err := mapstructure.Decode(configurations, receiverConfig); err != nil {
		return nil, fmt.Errorf("failed to transform configurations to receiver config: %w", err)
	}

	if err := receiverConfig.Validate(); err != nil {
		return nil, err
	}

	token, err := s.cryptoClient.Decrypt(receiverConfig.Token)
	if err != nil {
		return nil, fmt.Errorf("slack token decryption failed: %w", err)
	}

	receiverConfig.Token = token

	return receiverConfig.AsMap(), nil
}

// BuildData populates receiver data field based on config
func (s *PluginService) BuildData(ctx context.Context, configurations map[string]any) (map[string]any, error) {
	receiverConfig := &ReceiverConfig{}
	if err := mapstructure.Decode(configurations, receiverConfig); err != nil {
		return nil, fmt.Errorf("failed to transform configurations to receiver config: %w", err)
	}

	if err := receiverConfig.Validate(); err != nil {
		return nil, err
	}

	channels, err := s.client.GetWorkspaceChannels(
		ctx,
		receiverConfig.Token,
	)
	if err != nil {
		return nil, fmt.Errorf("could not get channels: %w", err)
	}

	data, err := json.Marshal(channels)
	if err != nil {
		// this is very unlikely to return error since we have an explicitly defined type of channels
		return nil, fmt.Errorf("invalid channels: %w", err)
	}

	receiverData := ReceiverData{
		Channels: string(data),
	}

	return receiverData.AsMap(), nil
}

func (s *PluginService) PreHookQueueTransformConfigs(ctx context.Context, notificationConfigMap map[string]any) (map[string]any, error) {
	notificationConfig := &NotificationConfig{}
	if err := mapstructure.Decode(notificationConfigMap, notificationConfig); err != nil {
		return nil, fmt.Errorf("failed to transform configurations to slack notification config: %w", err)
	}

	if err := notificationConfig.Validate(); err != nil {
		return nil, err
	}

	cipher, err := s.cryptoClient.Encrypt(notificationConfig.Token)
	if err != nil {
		return nil, fmt.Errorf("slack token encryption failed: %w", err)
	}

	notificationConfig.Token = cipher

	return notificationConfig.AsMap(), nil
}

func (s *PluginService) PostHookQueueTransformConfigs(ctx context.Context, notificationConfigMap map[string]any) (map[string]any, error) {
	notificationConfig := &NotificationConfig{}
	if err := mapstructure.Decode(notificationConfigMap, notificationConfig); err != nil {
		return nil, fmt.Errorf("failed to transform configurations to notification config: %w", err)
	}

	if err := notificationConfig.Validate(); err != nil {
		return nil, err
	}

	token, err := s.cryptoClient.Decrypt(notificationConfig.Token)
	if err != nil {
		return nil, fmt.Errorf("slack token decryption failed: %w", err)
	}

	notificationConfig.Token = token

	return notificationConfig.AsMap(), nil
}

func (s *PluginService) Send(ctx context.Context, notificationMessage notification.Message) (bool, error) {
	notificationConfig := &NotificationConfig{}
	if err := mapstructure.Decode(notificationMessage.Configs, notificationConfig); err != nil {
		return false, err
	}

	slackMessage := &Message{}
	if err := mapstructure.Decode(notificationMessage.Details, &slackMessage); err != nil {
		return false, err
	}

	if notificationConfig.ChannelType == "" {
		notificationConfig.ChannelType = defaultChannelType
	}
	if notificationConfig.ChannelName != "" {
		slackMessage.Channel = notificationConfig.ChannelName
	}

	if err := s.client.Notify(ctx, *notificationConfig, *slackMessage); err != nil {
		if errors.As(err, new(retry.RetryableError)) {
			return true, err
		} else {
			return false, err
		}
	}

	return false, nil
}

func (s *PluginService) GetSystemDefaultTemplate() string {
	return defaultAlertTemplateBody
}
