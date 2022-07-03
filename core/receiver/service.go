package receiver

import (
	"context"

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
	typeService, exist := s.registry[receiverType]
	if !exist {
		return nil, errors.ErrInvalid.WithMsgf("unsupported receiver type: %q", receiverType)
	}
	return typeService, nil
}

func (s *Service) List(ctx context.Context, flt Filter) ([]Receiver, error) {
	receivers, err := s.repository.List(ctx, flt)
	if err != nil {
		return nil, err
	}

	domainReceivers := make([]Receiver, 0, len(receivers))
	for i := 0; i < len(receivers); i++ {
		rcv := receivers[i]

		typeService, err := s.getTypeService(rcv.Type)
		if err != nil {
			return nil, err
		}
		if err = typeService.Decrypt(&rcv); err != nil {
			return nil, err
		}

		domainReceivers = append(domainReceivers, rcv)
	}
	return domainReceivers, nil
}

func (s *Service) Create(ctx context.Context, rcv *Receiver) error {
	typeService, err := s.getTypeService(rcv.Type)
	if err != nil {
		return err
	}

	if err := typeService.ValidateConfiguration(rcv); err != nil {
		return errors.ErrInvalid.WithMsgf(err.Error())
	}

	if err := typeService.Encrypt(rcv); err != nil {
		return err
	}

	err = s.repository.Create(ctx, rcv)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) Get(ctx context.Context, id uint64) (*Receiver, error) {
	rcv, err := s.repository.Get(ctx, id)
	if err != nil {
		if errors.As(err, new(NotFoundError)) {
			return nil, errors.ErrNotFound.WithMsgf(err.Error())
		}
		return nil, err
	}

	typeService, err := s.getTypeService(rcv.Type)
	if err != nil {
		return nil, err
	}

	if err := typeService.Decrypt(rcv); err != nil {
		return nil, err
	}

	return typeService.PopulateReceiver(ctx, rcv)
}

func (s *Service) Update(ctx context.Context, rcv *Receiver) error {
	typeService, err := s.getTypeService(rcv.Type)
	if err != nil {
		return err
	}

	if err := typeService.ValidateConfiguration(rcv); err != nil {
		return err
	}

	if err := typeService.Encrypt(rcv); err != nil {
		return err
	}

	err = s.repository.Update(ctx, rcv)
	if err != nil {
		if errors.As(err, new(NotFoundError)) {
			return errors.ErrNotFound.WithMsgf(err.Error())
		}
		return err
	}
	return nil
}

func (s *Service) Delete(ctx context.Context, id uint64) error {
	return s.repository.Delete(ctx, id)
}

func (s *Service) Notify(ctx context.Context, id uint64, payloadMessage NotificationMessage) error {
	rcv, err := s.Get(ctx, id)
	if err != nil {
		return errors.ErrInvalid.WithMsgf("error getting receiver with id %d", id).WithCausef(err.Error())
	}

	typeService, err := s.getTypeService(rcv.Type)
	if err != nil {
		return err
	}

	return typeService.Notify(ctx, rcv, payloadMessage)
}

func (s *Service) GetSubscriptionConfig(subsConfs map[string]string, rcv *Receiver) (map[string]string, error) {
	if rcv == nil {
		return nil, errors.ErrInvalid.WithCausef("receiver is nil")
	}

	typeService, err := s.getTypeService(rcv.Type)
	if err != nil {
		return nil, err
	}

	return typeService.GetSubscriptionConfig(subsConfs, rcv.Configurations)
}
