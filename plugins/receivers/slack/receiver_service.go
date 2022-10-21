package slack

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mitchellh/mapstructure"
	"github.com/odpf/siren/pkg/errors"
	"github.com/odpf/siren/plugins/receivers/base"
)

// ReceiverService is a receiver plugin service layer for slack
type ReceiverService struct {
	base.UnimplementedReceiverService
	client       SlackCaller
	cryptoClient Encryptor
}

// NewReceiverService returns slack service struct. This service implement [receiver.Resolver] interface.
func NewReceiverService(client SlackCaller, cryptoClient Encryptor) *ReceiverService {
	return &ReceiverService{
		client:       client,
		cryptoClient: cryptoClient,
	}
}

func (s *ReceiverService) PreHookTransformConfigs(ctx context.Context, configurations map[string]interface{}) (map[string]interface{}, error) {
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
func (s *ReceiverService) PostHookTransformConfigs(ctx context.Context, configurations map[string]interface{}) (map[string]interface{}, error) {
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
func (s *ReceiverService) BuildData(ctx context.Context, configurations map[string]interface{}) (map[string]interface{}, error) {
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

func (s *ReceiverService) BuildNotificationConfig(subsConfs map[string]interface{}, receiverConfs map[string]interface{}) (map[string]interface{}, error) {
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
