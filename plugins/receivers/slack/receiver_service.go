package slack

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mitchellh/mapstructure"
	"github.com/odpf/siren/pkg/errors"
)

// SlackReceiverService is a receiver plugin service layer for slack
type SlackReceiverService struct {
	slackClient  SlackClient
	cryptoClient Encryptor
}

// NewReceiverService returns slack service struct. This service implement [receiver.Resolver] interface.
func NewReceiverService(slackClient SlackClient, cryptoClient Encryptor) *SlackReceiverService {
	return &SlackReceiverService{
		slackClient:  slackClient,
		cryptoClient: cryptoClient,
	}
}

// TODO to be removed
// Deprecated: use Publish and SlackNotificationService instead
func (s *SlackReceiverService) Notify(ctx context.Context, configurations map[string]interface{}, payloadMessage map[string]interface{}) error {
	token, ok := configurations["token"].(string)
	if !ok {
		return errors.ErrInvalid.WithMsgf("no token in configurations found")
	}

	sm, err := GetSlackMessage(payloadMessage)
	if err != nil {
		return err
	}

	if err := s.slackClient.Notify(ctx, sm, CallWithToken(token)); err != nil {
		return fmt.Errorf("error calling slack notify: %w", err)
	}

	return nil
}

func (s *SlackReceiverService) PreHookTransformConfigs(ctx context.Context, configurations map[string]interface{}) (map[string]interface{}, error) {
	slackCredentialConfig := &SlackCredentialConfig{}
	if err := mapstructure.Decode(configurations, slackCredentialConfig); err != nil {
		return nil, fmt.Errorf("failed to transform configurations to pre transform config: %w", err)
	}

	if err := slackCredentialConfig.Validate(); err != nil {
		return nil, errors.ErrInvalid.WithMsgf(err.Error())
	}

	creds, err := s.slackClient.ExchangeAuth(ctx, slackCredentialConfig.AuthCode, slackCredentialConfig.ClientID, slackCredentialConfig.ClientSecret)
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
func (s *SlackReceiverService) PostHookTransformConfigs(ctx context.Context, configurations map[string]interface{}) (map[string]interface{}, error) {
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
func (s *SlackReceiverService) BuildData(ctx context.Context, configurations map[string]interface{}) (map[string]interface{}, error) {
	receiverConfig := &ReceiverConfig{}
	if err := mapstructure.Decode(configurations, receiverConfig); err != nil {
		return nil, fmt.Errorf("failed to transform configurations to receiver config: %w", err)
	}

	if err := receiverConfig.Validate(); err != nil {
		return nil, err
	}

	channels, err := s.slackClient.GetWorkspaceChannels(
		ctx,
		CallWithToken(receiverConfig.Token),
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

func (s *SlackReceiverService) BuildNotificationConfig(subsConfs map[string]interface{}, receiverConfs map[string]interface{}) (map[string]interface{}, error) {
	subscriptionConfig := &SubscriptionConfig{}
	if err := mapstructure.Decode(subsConfs, subscriptionConfig); err != nil {
		return nil, fmt.Errorf("failed to transform configurations to subscription config: %w", err)
	}

	receiverConfig := &ReceiverConfig{}
	if err := mapstructure.Decode(receiverConfs, receiverConfig); err != nil {
		return nil, fmt.Errorf("failed to transform configurations to subscription config: %w", err)
	}

	notificationConfig := NotificationConfig{
		SubscriptionConfig: *subscriptionConfig,
		ReceiverConfig:     *receiverConfig,
	}

	return notificationConfig.AsMap(), nil
}
