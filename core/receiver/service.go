package receiver

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/odpf/siren/pkg/slack"
	"github.com/pkg/errors"
	goslack "github.com/slack-go/slack"
)

const (
	Slack string = "slack"
)

// Service handles business logic
type Service struct {
	repository   Repository
	slackClient  SlackClient
	cryptoClient Encryptor
}

// NewService returns service struct
func NewService(repository Repository, slackClient SlackClient, cryptoClient Encryptor) *Service {
	return &Service{
		repository:   repository,
		slackClient:  slackClient,
		cryptoClient: cryptoClient,
	}
}

func (s *Service) ListReceivers() ([]*Receiver, error) {
	receivers, err := s.repository.List()
	if err != nil {
		return nil, fmt.Errorf("secureService.repository.List: %w", err)
	}

	domainReceivers := make([]*Receiver, 0, len(receivers))
	for i := 0; i < len(receivers); i++ {
		rcv := receivers[i]

		if rcv.Type == Slack {
			if err = s.postTransform(rcv); err != nil {
				return nil, err
			}
		}

		domainReceivers = append(domainReceivers, rcv)
	}
	return domainReceivers, nil
}

func (s *Service) CreateReceiver(rcv *Receiver) error {
	if rcv.Type == Slack {
		if err := s.preTransform(rcv); err != nil {
			return err
		}
	}

	if err := s.repository.Create(rcv); err != nil {
		return fmt.Errorf("secureService.repository.Create: %w", err)
	}

	if rcv.Type == Slack {
		if err := s.postTransform(rcv); err != nil {
			return err
		}
	}
	return nil
}

func (s *Service) GetReceiver(id uint64) (*Receiver, error) {
	rcv, err := s.repository.Get(id)
	if err != nil {
		return nil, fmt.Errorf("secureService.repository.Get: %w", err)
	}

	if rcv.Type == Slack {
		if err := s.postTransform(rcv); err != nil {
			return nil, err
		}

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
	}
	return rcv, nil
}

func (s *Service) UpdateReceiver(rcv *Receiver) error {
	if rcv.Type == Slack {
		if err := s.preTransform(rcv); err != nil {
			return err
		}
	}

	if err := s.repository.Update(rcv); err != nil {
		return fmt.Errorf("secureService.repository.Update: %w", err)
	}
	return nil
}

func (s *Service) DeleteReceiver(id uint64) error {
	return s.repository.Delete(id)
}

func (s *Service) Migrate() error {
	return s.repository.Migrate()
}

func (s *Service) NotifyReceiver(rcv *Receiver, payloadMessage string, payloadReceiverName string, payloadReceiverType string, payloadBlock []byte) error {
	switch rcv.Type {
	case Slack:
		blocks := goslack.Blocks{}
		if err := json.Unmarshal(payloadBlock, &blocks); err != nil {
			return fmt.Errorf("unable to parse slack block: %w", ErrInvalid)
		}

		token, ok := rcv.Configurations["token"].(string)
		if !ok {
			return fmt.Errorf("no token found in configuration: %w", ErrInvalid)
		}

		payloadMessage := &slack.Message{
			ReceiverName: payloadReceiverName,
			ReceiverType: payloadReceiverType,
			Token:        rcv.Configurations["token"].(string),
			Message:      payloadMessage,
			Blocks:       blocks,
		}
		if err := s.slackClient.Notify(payloadMessage, slack.CallWithToken(token)); err != nil {
			return fmt.Errorf("failed to notify: %w", err)
		}

	default:
		return errors.New("type not recognized")
	}
	return nil
}

func (s *Service) preTransform(r *Receiver) error {
	var token string
	var ok bool
	if token, ok = r.Configurations["token"].(string); !ok {
		return errors.New("no token field found")
	}
	chiperText, err := s.cryptoClient.Encrypt(token)
	if err != nil {
		return fmt.Errorf("pre transform encrypt failed: %w", err)
	}
	r.Configurations["token"] = chiperText

	return nil
}

func (s *Service) postTransform(r *Receiver) error {
	var cipherText string
	var ok bool
	if cipherText, ok = r.Configurations["token"].(string); !ok {
		return errors.New("no token field found")
	}
	token, err := s.cryptoClient.Decrypt(cipherText)
	if err != nil {
		return fmt.Errorf("post transform decrypt failed: %w", err)
	}
	r.Configurations["token"] = token
	return nil
}
