package receiver

import (
	"fmt"
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
		return nil, fmt.Errorf("%w: unsupported receiver type", ErrInvalid)
	}
	return strategyService, nil
}

func (s *Service) ListReceivers() ([]*Receiver, error) {
	receivers, err := s.repository.List()
	if err != nil {
		return nil, fmt.Errorf("secureService.repository.List: %w", err)
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
		return fmt.Errorf("%w: invalid receiver configurations", ErrInvalid)
	}

	if err := strategyService.Encrypt(rcv); err != nil {
		return err
	}

	if err := s.repository.Create(rcv); err != nil {
		return fmt.Errorf("secureService.repository.Create: %w", err)
	}

	if err := strategyService.Decrypt(rcv); err != nil {
		return err
	}
	return nil
}

func (s *Service) GetReceiver(id uint64) (*Receiver, error) {
	rcv, err := s.repository.Get(id)
	if err != nil {
		return nil, fmt.Errorf("secureService.repository.Get: %w", err)
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
		return fmt.Errorf("%w: invalid receiver configurations", ErrInvalid)
	}

	if err := strategyService.Encrypt(rcv); err != nil {
		return err
	}

	if err := s.repository.Update(rcv); err != nil {
		return fmt.Errorf("secureService.repository.Update: %w", err)
	}
	return nil
}

func (s *Service) NotifyReceiver(id uint64, payloadMessage NotificationMessage) error {
	rcv, err := s.GetReceiver(id)
	if err != nil {
		return fmt.Errorf("%w: error getting receiver with id %d", ErrInvalid, id)
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

func (s *Service) Migrate() error {
	return s.repository.Migrate()
}
