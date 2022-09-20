package receiver

import (
	"context"

	"github.com/odpf/siren/pkg/errors"
)

// Service handles business logic
type Service struct {
	registry   map[string]Resolver
	repository Repository
}

func NewService(repository Repository, registry map[string]Resolver) *Service {
	return &Service{
		repository: repository,
		registry:   registry,
	}
}

func (s *Service) getReceiverPlugin(receiverType string) (Resolver, error) {
	receiverPlugin, exist := s.registry[receiverType]
	if !exist {
		return nil, errors.ErrInvalid.WithMsgf("unsupported receiver type: %q", receiverType)
	}
	return receiverPlugin, nil
}

func (s *Service) List(ctx context.Context, flt Filter) ([]Receiver, error) {
	receivers, err := s.repository.List(ctx, flt)
	if err != nil {
		return nil, err
	}

	domainReceivers := make([]Receiver, 0, len(receivers))
	for i := 0; i < len(receivers); i++ {
		rcv := receivers[i]

		receiverPlugin, err := s.getReceiverPlugin(rcv.Type)
		if err != nil {
			return nil, err
		}
		transformedConfigs, err := receiverPlugin.PostHookTransformConfigs(ctx, rcv.Configurations)
		if err != nil {
			return nil, err
		}
		rcv.Configurations = transformedConfigs

		domainReceivers = append(domainReceivers, rcv)
	}
	return domainReceivers, nil
}

func (s *Service) Create(ctx context.Context, rcv *Receiver) error {
	receiverPlugin, err := s.getReceiverPlugin(rcv.Type)
	if err != nil {
		return err
	}

	if err := receiverPlugin.ValidateConfigurations(rcv.Configurations); err != nil {
		return errors.ErrInvalid.WithMsgf(err.Error())
	}

	rcv.Configurations, err = receiverPlugin.PreHookTransformConfigs(ctx, rcv.Configurations)
	if err != nil {
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

	receiverPlugin, err := s.getReceiverPlugin(rcv.Type)
	if err != nil {
		return nil, err
	}

	transformedConfigs, err := receiverPlugin.PostHookTransformConfigs(ctx, rcv.Configurations)
	if err != nil {
		return nil, err
	}
	rcv.Configurations = transformedConfigs

	populatedData, err := receiverPlugin.PopulateDataFromConfigs(ctx, rcv.Configurations)
	if err != nil {
		return nil, err
	}

	rcv.Data = populatedData

	return rcv, nil
}

func (s *Service) Update(ctx context.Context, rcv *Receiver) error {
	receiverPlugin, err := s.getReceiverPlugin(rcv.Type)
	if err != nil {
		return err
	}

	if err := receiverPlugin.ValidateConfigurations(rcv.Configurations); err != nil {
		return err
	}

	rcv.Configurations, err = receiverPlugin.PreHookTransformConfigs(ctx, rcv.Configurations)
	if err != nil {
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

func (s *Service) Notify(ctx context.Context, id uint64, payloadMessage map[string]interface{}) error {
	rcv, err := s.Get(ctx, id)
	if err != nil {
		return errors.ErrInvalid.WithMsgf("error getting receiver with id %d", id).WithCausef(err.Error())
	}

	receiverPlugin, err := s.getReceiverPlugin(rcv.Type)
	if err != nil {
		return err
	}

	return receiverPlugin.Notify(ctx, rcv.Configurations, payloadMessage)
}

func (s *Service) EnrichSubscriptionConfig(subsConfs map[string]string, rcv *Receiver) (map[string]string, error) {
	if rcv == nil {
		return nil, errors.ErrInvalid.WithCausef("receiver is nil")
	}

	receiverPlugin, err := s.getReceiverPlugin(rcv.Type)
	if err != nil {
		return nil, err
	}

	return receiverPlugin.EnrichSubscriptionConfig(subsConfs, rcv.Configurations)
}
