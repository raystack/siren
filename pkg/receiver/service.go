package receiver

import (
	"encoding/json"
	"fmt"
	"github.com/odpf/siren/domain"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

const (
	Slack string = "slack"
)

type CodeExchangeHTTPResponse struct {
	AccessToken string `json:"access_token"`
	Team        struct {
		Name string `json:"name"`
	} `json:"team"`
	Ok bool `json:"ok"`
}

// Service handles business logic
type Service struct {
	repository      ReceiverRepository
	slackRepository SlackRepository
	slackHelper     SlackHelper
}

// NewService returns repository struct
func NewService(db *gorm.DB, httpClient Doer, encryptionKey string) (domain.ReceiverService, error) {
	repository := NewRepository(db)
	slackHelper, err := NewSlackHelper(httpClient, encryptionKey)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create slack helper")
	}

	return &Service{
		repository:      repository,
		slackHelper:     slackHelper,
		slackRepository: NewSlackRepository(),
	}, nil
}

func (service Service) ListReceivers() ([]*domain.Receiver, error) {
	receivers, err := service.repository.List()
	if err != nil {
		return nil, errors.Wrap(err, "service.repository.List")
	}

	domainReceivers := make([]*domain.Receiver, 0, len(receivers))
	for i := 0; i < len(receivers); i++ {
		receiver := receivers[i].toDomain()
		domainReceivers = append(domainReceivers, receiver)
	}

	return domainReceivers, nil

}

func (service Service) CreateReceiver(receiver *domain.Receiver) (*domain.Receiver, error) {
	var data *domain.Receiver
	var err error
	p := &Receiver{}

	if receiver.Type == Slack {
		data, err = service.slackHelper.Transform(receiver)
		if err != nil {
			return nil, errors.Wrap(err, "slackHelper.Transform")
		}
	}

	payload := p.fromDomain(data)
	newReceiver, err := service.repository.Create(payload)
	if err != nil {
		return nil, errors.Wrap(err, "service.repository.Create")
	}

	return newReceiver.toDomain(), nil
}

func (service Service) GetReceiver(id uint64) (*domain.Receiver, error) {
	receiver, err := service.repository.Get(id)
	if err != nil {
		return nil, err
	}

	if receiver.Type == Slack {
		encryptedToken := receiver.Configurations["token"].(string)
		token, err := service.slackHelper.Decrypt(encryptedToken)
		if err != nil {
			return nil, errors.Wrap(err, "slackHelper.Decrypt")
		}

		receiver.Configurations["token"] = token

		channels, err := service.slackRepository.GetWorkspaceChannels(token)
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("could not get channels"))
		}

		data, err := json.Marshal(channels)
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("invalid channels"))
		}

		receiver.Data = make(map[string]interface{})
		receiver.Data["channels"] = string(data)
	}

	return receiver.toDomain(), nil
}

func (service Service) UpdateReceiver(receiver *domain.Receiver) (*domain.Receiver, error) {
	var data *domain.Receiver
	var err error
	p := &Receiver{}

	if receiver.Type == Slack {
		data, err = service.slackHelper.Transform(receiver)
		if err != nil {
			return nil, errors.Wrap(err, "slackHelper.Transform")
		}
	}

	payload := p.fromDomain(data)
	newReceiver, err := service.repository.Update(payload)
	if err != nil {
		return nil, err
	}

	return newReceiver.toDomain(), nil
}

func (service Service) DeleteReceiver(id uint64) error {
	return service.repository.Delete(id)
}

func (service Service) Migrate() error {
	return service.repository.Migrate()
}