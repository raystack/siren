package receiver

import (
	"encoding/json"

	"github.com/odpf/siren/domain"
	"github.com/odpf/siren/plugins/receivers/http"
	"github.com/odpf/siren/plugins/receivers/slack"
	"github.com/odpf/siren/store"
	"github.com/pkg/errors"
)

const (
	Slack string = "slack"
)

var (
	jsonMarshal = json.Marshal
)

type slackService interface {
	GetWorkspaceChannels(string) ([]slack.Channel, error)
}

// Service handles business logic
type Service struct {
	repository   store.ReceiverRepository
	slackService slackService
	slackHelper  SlackHelper
}

// NewService returns service struct
func NewService(repository store.ReceiverRepository, httpClient http.Doer, encryptionKey string) (domain.ReceiverService, error) {
	slackHelper, err := NewSlackHelper(httpClient, encryptionKey)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create slack helper")
	}

	return &Service{
		repository:   repository,
		slackHelper:  slackHelper,
		slackService: slack.NewService(),
	}, nil
}

func (service Service) ListReceivers() ([]*domain.Receiver, error) {
	receivers, err := service.repository.List()
	if err != nil {
		return nil, errors.Wrap(err, "service.repository.List")
	}

	domainReceivers := make([]*domain.Receiver, 0, len(receivers))
	for i := 0; i < len(receivers); i++ {
		receiver := receivers[i]

		if receiver.Type == Slack {
			if err = service.slackHelper.PostTransform(receiver); err != nil {
				return nil, errors.Wrap(err, "slackHelper.PostTransform")
			}
		}

		domainReceivers = append(domainReceivers, receiver)
	}

	return domainReceivers, nil

}

func (service Service) CreateReceiver(receiver *domain.Receiver) error {
	if receiver.Type == Slack {
		if err := service.slackHelper.PreTransform(receiver); err != nil {
			return errors.Wrap(err, "slackHelper.PreTransform")
		}
	}

	if err := service.repository.Create(receiver); err != nil {
		return errors.Wrap(err, "service.repository.Create")
	}

	if receiver.Type == Slack {
		if err := service.slackHelper.PostTransform(receiver); err != nil {
			return errors.Wrap(err, "slackHelper.PostTransform")
		}
	}

	return nil
}

func (service Service) GetReceiver(id uint64) (*domain.Receiver, error) {
	receiver, err := service.repository.Get(id)
	if err != nil {
		return nil, err
	}

	if receiver.Type == Slack {
		if err := service.slackHelper.PostTransform(receiver); err != nil {
			return nil, errors.Wrap(err, "slackHelper.PostTransform")
		}

		token := receiver.Configurations["token"].(string)
		channels, err := service.slackService.GetWorkspaceChannels(token)
		if err != nil {
			return nil, errors.Wrap(err, "could not get channels")
		}

		data, err := jsonMarshal(channels)
		if err != nil {
			return nil, errors.Wrap(err, "invalid channels")
		}

		receiver.Data = make(map[string]interface{})
		receiver.Data["channels"] = string(data)
	}

	return receiver, nil
}

func (service Service) UpdateReceiver(receiver *domain.Receiver) error {
	if receiver.Type == Slack {
		if err := service.slackHelper.PreTransform(receiver); err != nil {
			return errors.Wrap(err, "slackHelper.PreTransform")
		}
	}

	if err := service.repository.Update(receiver); err != nil {
		return err
	}

	return nil
}

func (service Service) DeleteReceiver(id uint64) error {
	return service.repository.Delete(id)
}

func (service Service) Migrate() error {
	return service.repository.Migrate()
}
