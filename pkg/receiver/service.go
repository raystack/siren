package receiver

import (
	"github.com/odpf/siren/domain"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

// Service handles business logic
type Service struct {
	repository ReceiverRepository
}

// NewService returns repository struct
func NewService(db *gorm.DB) domain.ReceiverService {
	return &Service{NewRepository(db)}
}

func (service Service) ListReceivers() ([]*domain.Receiver, error) {
	receivers, err := service.repository.List()
	if err != nil {
		return nil, errors.Wrap(err, "service.repository.List")
	}

	domainReceivers := make([]*domain.Receiver, 0, len(receivers))
	for i := 0; i < len(receivers); i++ {
		provider := receivers[i].toDomain()
		domainReceivers = append(domainReceivers, provider)
	}

	return domainReceivers, nil

}

func (service Service) CreateReceiver(receiver *domain.Receiver) (*domain.Receiver, error) {
	p := &Receiver{}
	newReceiver, err := service.repository.Create(p.fromDomain(receiver))
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
	w := &Receiver{}
	newReceiver, err := service.repository.Update(w.fromDomain(receiver))
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
