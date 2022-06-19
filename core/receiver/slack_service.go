package receiver

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/odpf/siren/pkg/slack"
	"github.com/pkg/errors"
)

type SlackService struct {
	slackClient  SlackClient
	cryptoClient Encryptor
}

// NewService returns slack service struct
func NewSlackService(slackClient SlackClient, cryptoClient Encryptor) *SlackService {
	return &SlackService{
		slackClient:  slackClient,
		cryptoClient: cryptoClient,
	}
}

func (s *SlackService) Notify(rcv *Receiver, payloadMessage NotificationMessage) error {
	token, ok := rcv.Configurations["token"].(string)
	if !ok {
		return errors.New("no token in receiver configurations found")
	}

	sm, err := payloadMessage.ToSlackMessage()
	if err != nil {
		return err
	}

	if err := s.slackClient.Notify(sm, slack.CallWithToken(token)); err != nil {
		return fmt.Errorf("failed to notify: %w", err)
	}

	return nil
}

func (s *SlackService) Encrypt(r *Receiver) error {
	var token string
	var ok bool
	if token, ok = r.Configurations["token"].(string); !ok {
		return errors.New("no token field found")
	}
	chiperText, err := s.cryptoClient.Encrypt(token)
	if err != nil {
		return err
	}
	r.Configurations["token"] = chiperText

	return nil
}

func (s *SlackService) Decrypt(r *Receiver) error {
	var cipherText string
	var ok bool
	if cipherText, ok = r.Configurations["token"].(string); !ok {
		return errors.New("no token field found")
	}
	token, err := s.cryptoClient.Decrypt(cipherText)
	if err != nil {
		return err
	}
	r.Configurations["token"] = token
	return nil
}

func (s *SlackService) PopulateReceiver(rcv *Receiver) (*Receiver, error) {
	token, ok := rcv.Configurations["token"].(string)
	if !ok {
		return nil, errors.New("no token found in configurations")
	}

	channels, err := s.slackClient.GetWorkspaceChannels(
		slack.CallWithContext(context.Background()),
		slack.CallWithToken(token),
	)
	if err != nil {
		return nil, fmt.Errorf("could not get channels: %w", err)
	}

	data, err := json.Marshal(channels)
	if err != nil {
		// this is very unlikely to return error since we have an explicitly defined type of channels
		return nil, fmt.Errorf("invalid channels: %w", err)
	}

	rcv.Data = make(map[string]interface{})
	rcv.Data["channels"] = string(data)

	return rcv, nil
}

func (s *SlackService) ValidateConfiguration(configurations Configurations) error {
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
