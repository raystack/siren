package receiver

import (
	"context"

	"github.com/odpf/siren/pkg/errors"
	"github.com/odpf/siren/pkg/telemetry"
	"go.opencensus.io/tag"
)

// Service handles business logic
type Service struct {
	receiverPlugins map[string]ConfigResolver
	repository      Repository
}

func NewService(repository Repository, receiverPlugins map[string]ConfigResolver) *Service {
	return &Service{
		repository:      repository,
		receiverPlugins: receiverPlugins,
	}
}

func (s *Service) getReceiverPlugin(receiverType string) (ConfigResolver, error) {
	receiverPlugin, exist := s.receiverPlugins[receiverType]
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
		transformedConfigs, err := receiverPlugin.PostHookDBTransformConfigs(ctx, rcv.Configurations)
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

	rcv.Configurations, err = receiverPlugin.PreHookDBTransformConfigs(ctx, rcv.Configurations)
	if err != nil {
		telemetry.IncrementInt64Counter(ctx, telemetry.MetricReceiverPreHookDBFailed,
			tag.Upsert(telemetry.TagReceiverType, rcv.Type))

		return err
	}

	err = s.repository.Create(ctx, rcv)
	if err != nil {
		return err
	}

	return nil
}

type getOpts struct {
	withData bool
}

type GetOption func(*getOpts)

func GetWithData(withData bool) GetOption {
	return func(g *getOpts) {
		g.withData = withData
	}
}

func (s *Service) Get(ctx context.Context, id uint64, gopts ...GetOption) (*Receiver, error) {
	opt := &getOpts{}

	for _, g := range gopts {
		g(opt)
	}

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

	transformedConfigs, err := receiverPlugin.PostHookDBTransformConfigs(ctx, rcv.Configurations)
	if err != nil {
		telemetry.IncrementInt64Counter(ctx, telemetry.MetricReceiverPostHookDBFailed,
			tag.Upsert(telemetry.TagReceiverType, rcv.Type))

		return nil, err
	}
	rcv.Configurations = transformedConfigs

	if opt.withData {
		populatedData, err := receiverPlugin.BuildData(ctx, rcv.Configurations)
		if err != nil {
			return nil, err
		}

		rcv.Data = populatedData
	}

	return rcv, nil
}

func (s *Service) Update(ctx context.Context, rcv *Receiver) error {
	rcv, err := s.repository.Get(ctx, rcv.ID)
	if err != nil {
		if errors.As(err, new(NotFoundError)) {
			return errors.ErrNotFound.WithMsgf(err.Error())
		}
		return err
	}

	receiverPlugin, err := s.getReceiverPlugin(rcv.Type)
	if err != nil {
		return err
	}

	rcv.Configurations, err = receiverPlugin.PreHookDBTransformConfigs(ctx, rcv.Configurations)
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
