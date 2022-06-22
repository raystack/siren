package receiver

import (
	"github.com/odpf/siren/pkg/errors"
)

// Service handles business logic
type Service struct {
	registry   map[string]TypeService
	repository Repository
}

func NewService(repository Repository, registry map[string]TypeService) *Service {
	return &Service{
		repository: repository,
		registry:   registry,
	}
}

func (s *Service) getTypeService(receiverType string) (TypeService, error) {
	strategyService, exist := s.registry[receiverType]
	if !exist {
		return nil, errors.ErrInvalid.WithMsgf("unsupported receiver type: %q", receiverType)
	}
	return strategyService, nil
}

func (s *Service) ListReceivers() ([]*Receiver, error) {
	receivers, err := s.repository.List()
	if err != nil {
		return nil, err
	}

	domainReceivers := make([]*Receiver, 0, len(receivers))
	for i := 0; i < len(receivers); i++ {
		rcv := receivers[i]

		strategyService, err := s.getTypeService(rcv.Type)
		if err != nil {
			return nil, err
		}
		if err = strategyService.Decrypt(rcv); err != nil {
			return nil, err
		}

		domainReceivers = append(domainReceivers, rcv)
	}
	return domainReceivers, nil
}

func (s *Service) CreateReceiver(rcv *Receiver) error {
	strategyService, err := s.getTypeService(rcv.Type)
	if err != nil {
		return err
	}

	if err := strategyService.ValidateConfiguration(rcv.Configurations); err != nil {
		return err
	}

	if err := strategyService.Encrypt(rcv); err != nil {
		return err
	}

	if err := s.repository.Create(rcv); err != nil {
		return err
	}

	if err := strategyService.Decrypt(rcv); err != nil {
		return err
	}
	return nil
}

func (s *Service) GetReceiver(id uint64) (*Receiver, error) {
	rcv, err := s.repository.Get(id)
	if err != nil {
		if errors.As(err, new(NotFoundError)) {
			return nil, errors.ErrNotFound.WithMsgf(err.Error())
		}
		return nil, err
	}

	strategyService, err := s.getTypeService(rcv.Type)
	if err != nil {
		return nil, err
	}

	if err := strategyService.Decrypt(rcv); err != nil {
		return nil, err
	}

	return strategyService.PopulateReceiver(rcv)
}

func (s *Service) UpdateReceiver(rcv *Receiver) error {
	strategyService, err := s.getTypeService(rcv.Type)
	if err != nil {
		return err
	}

	if err := strategyService.ValidateConfiguration(rcv.Configurations); err != nil {
		return err
	}

	if err := strategyService.Encrypt(rcv); err != nil {
		return err
	}

	if err := s.repository.Update(rcv); err != nil {
		if errors.As(err, new(NotFoundError)) {
			return errors.ErrNotFound.WithMsgf(err.Error())
		}
		return err
	}
	return nil
}

func (s *Service) NotifyReceiver(id uint64, payloadMessage NotificationMessage) error {
	rcv, err := s.GetReceiver(id)
	if err != nil {
		return errors.ErrInvalid.WithMsgf("error getting receiver with id %d", id).WithCausef(err.Error())
	}

	strategyService, err := s.getTypeService(rcv.Type)
	if err != nil {
		return err
	}

	return strategyService.Notify(rcv, payloadMessage)
}

func (s *Service) DeleteReceiver(id uint64) error {
	return s.repository.Delete(id)
}
