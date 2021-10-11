package receiver

import (
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
	repository ReceiverRepository
	exchanger  Exchanger
}

// NewService returns repository struct
func NewService(db *gorm.DB, httpClient Doer, encryptionKey string) (domain.ReceiverService, error) {
	repository, err := NewRepository(db, encryptionKey)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create receiver repository")
	}

	return &Service{
		repository: repository,
		exchanger:  NewSlackClient(httpClient),
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
	p := &Receiver{}
	payload := p.fromDomain(receiver)

	if payload.Type == Slack {
		configurations := payload.Configurations
		clientId := configurations["client_id"].(string)
		clientSecret := configurations["client_secret"].(string)
		code := configurations["auth_code"].(string)

		response, err := service.exchanger.Exchange(code, clientId, clientSecret)
		if err != nil {
			return nil, errors.Wrap(err, "failed to exchange code with slack OAuth server")
		}

		newConfigurations := map[string]interface{}{}
		newConfigurations["workspace"] = response.Team.Name
		newConfigurations["token"] = response.AccessToken
		payload.Configurations = newConfigurations
		//transormedReceiver := service.slackTransformer.transform()
	}

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

	return receiver.toDomain(), nil
}

func (service Service) UpdateReceiver(receiver *domain.Receiver) (*domain.Receiver, error) {
	p := &Receiver{}
	payload := p.fromDomain(receiver)

	if payload.Type == Slack {
		configurations := payload.Configurations
		clientId := configurations["client_id"].(string)
		clientSecret := configurations["client_secret"].(string)
		code := configurations["auth_code"].(string)

		response, err := service.exchanger.Exchange(code, clientId, clientSecret)
		if err != nil {
			return nil, errors.Wrap(err, "failed to exchange code with slack OAuth server")
		}

		newConfigurations := map[string]interface{}{}
		newConfigurations["workspace"] = response.Team.Name
		newConfigurations["token"] = response.AccessToken
		payload.Configurations = newConfigurations
	}

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
