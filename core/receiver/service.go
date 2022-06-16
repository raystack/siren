package receiver

import (
	"errors"
	"fmt"
)

// Service handles business logic
type Service struct {
	registry   map[string]StrategyService
	repository Repository
}

func NewService(repository Repository, registry map[string]StrategyService) *Service {
	return &Service{
		repository: repository,
		registry:   registry,
	}
}

func (s *Service) getStrategy(receiverType string) StrategyService {
	return s.registry[receiverType]
}

func (s *Service) ListReceivers() ([]*Receiver, error) {
	receivers, err := s.repository.List()
	if err != nil {
		return nil, fmt.Errorf("secureService.repository.List: %w", err)
	}

	domainReceivers := make([]*Receiver, 0, len(receivers))
	for i := 0; i < len(receivers); i++ {
		rcv := receivers[i]

		strategyService := s.getStrategy(rcv.Type)
		if strategyService == nil {
			//TODO log here
			continue
		}
		if err = strategyService.Decrypt(rcv); err != nil {
			return nil, err
		}

		domainReceivers = append(domainReceivers, rcv)
	}
	return domainReceivers, nil
}

func (s *Service) CreateReceiver(rcv *Receiver) error {
	strategyService := s.getStrategy(rcv.Type)
	if strategyService == nil {
		//TODO log here, adjust error
		return errors.New("unsupported receiver type")
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

	strategyService := s.getStrategy(rcv.Type)
	if strategyService == nil {
		//TODO log here, adjust error
		return nil, errors.New("unsupported receiver type")
	}

	if err := strategyService.Decrypt(rcv); err != nil {
		return nil, err
	}

	return strategyService.PopulateReceiver(rcv)
}

func (s *Service) UpdateReceiver(rcv *Receiver) error {
	strategyService := s.getStrategy(rcv.Type)
	if strategyService == nil {
		//TODO log here, adjust error
		return errors.New("unsupported receiver type")
	}

	if err := strategyService.Encrypt(rcv); err != nil {
		return err
	}

	if err := s.repository.Update(rcv); err != nil {
		return fmt.Errorf("secureService.repository.Update: %w", err)
	}
	return nil
}

func (s *Service) NotifyReceiver(rcv *Receiver, payloadMessage string, payloadReceiverName string, payloadReceiverType string, payloadBlock []byte) error {
	strategyService := s.getStrategy(rcv.Type)
	if strategyService == nil {
		//TODO log here, adjust error
		return errors.New("unsupported receiver type")
	}

	return strategyService.Notify(rcv, payloadMessage, payloadReceiverName, payloadReceiverType, payloadBlock)
}

func (s *Service) DeleteReceiver(id uint64) error {
	return s.repository.Delete(id)
}

func (s *Service) Migrate() error {
	return s.repository.Migrate()
}
