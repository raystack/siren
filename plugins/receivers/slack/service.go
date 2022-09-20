package slack

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/odpf/siren/core/receiver"
	"github.com/odpf/siren/pkg/errors"
)

//go:generate mockery --name=Encryptor -r --case underscore --with-expecter --structname Encryptor --filename encryptor.go --output=./mocks
type Encryptor interface {
	Encrypt(str string) (string, error)
	Decrypt(str string) (string, error)
}

//go:generate mockery --name=SlackClient -r --case underscore --with-expecter --structname SlackClient --filename slack_client.go --output=./mocks
type SlackClient interface {
	ExchangeAuth(ctx context.Context, authCode, clientID, clientSecret string) (Credential, error)
	GetWorkspaceChannels(ctx context.Context, opts ...ClientCallOption) ([]Channel, error)
	Notify(ctx context.Context, message *Message, opts ...ClientCallOption) error
}

// SlackService is a receiver plugin service layer for slack
type SlackService struct {
	slackClient  SlackClient
	cryptoClient Encryptor
}

// NewService returns slack service struct. This service implement [receiver.Resolver] interface.
func NewReceiverService(slackClient SlackClient, cryptoClient Encryptor) *SlackService {
	return &SlackService{
		slackClient:  slackClient,
		cryptoClient: cryptoClient,
	}
}

func (s *SlackService) Notify(ctx context.Context, configurations receiver.Configurations, payloadMessage map[string]interface{}) error {
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

func (s *SlackService) PreHookTransformConfigs(ctx context.Context, configurations receiver.Configurations) (receiver.Configurations, error) {
	clientID := configurations["client_id"].(string)
	clientSecret := configurations["client_secret"].(string)
	authCode := configurations["auth_code"].(string)

	creds, err := s.slackClient.ExchangeAuth(ctx, authCode, clientID, clientSecret)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code with slack OAuth server: %w", err)
	}

	cipherText, err := s.cryptoClient.Encrypt(creds.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("slack token encryption failed: %w", err)
	}

	newConfigurations := receiver.Configurations{}
	newConfigurations["workspace"] = creds.TeamName
	newConfigurations["token"] = cipherText

	return newConfigurations, nil
}

func (s *SlackService) PostHookTransformConfigs(ctx context.Context, configurations receiver.Configurations) (receiver.Configurations, error) {
	cipherText, err := configurations.GetString("token")
	if err != nil {
		return nil, err
	}

	token, err := s.cryptoClient.Decrypt(cipherText)
	if err != nil {
		return nil, fmt.Errorf("slack token decryption failed: %w", err)
	}

	configurations["token"] = token

	return configurations, nil
}

func (s *SlackService) PopulateDataFromConfigs(ctx context.Context, configurations receiver.Configurations) (map[string]interface{}, error) {
	token, ok := configurations["token"].(string)
	if !ok {
		return nil, errors.ErrInvalid.WithMsgf("no token in configurations found")
	}

	channels, err := s.slackClient.GetWorkspaceChannels(
		ctx,
		CallWithToken(token),
	)
	if err != nil {
		return nil, fmt.Errorf("could not get channels: %w", err)
	}

	data, err := json.Marshal(channels)
	if err != nil {
		// this is very unlikely to return error since we have an explicitly defined type of channels
		return nil, fmt.Errorf("invalid channels: %w", err)
	}

	var receiverData = make(map[string]interface{})
	receiverData["channels"] = string(data)

	return receiverData, nil
}

func (s *SlackService) ValidateConfigurations(configurations receiver.Configurations) error {
	_, err := configurations.GetString("client_id")
	if err != nil {
		return err
	}

	_, err = configurations.GetString("client_secret")
	if err != nil {
		return err
	}

	_, err = configurations.GetString("auth_code")
	if err != nil {
		return err
	}

	return nil
}

func (s *SlackService) EnrichSubscriptionConfig(subsConfs map[string]string, receiverConfs receiver.Configurations) (map[string]string, error) {
	mapConf := make(map[string]string)
	if _, ok := subsConfs["channel_name"]; !ok {
		return nil, errors.New("subscription receiver config 'channel_name' was missing")
	}
	mapConf["channel_name"] = subsConfs["channel_name"]
	if val, ok := receiverConfs["token"]; ok {
		if mapConf["token"], ok = val.(string); !ok {
			return nil, errors.New("token config from receiver should be in string")
		}
	}
	return mapConf, nil
}
