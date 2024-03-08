package receiver

import (
	"context"

	"github.com/goto/siren/pkg/errors"
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
	if err := rcv.Validate(); err != nil {
		return errors.ErrInvalid.WithMsgf("%s", err.Error())
	}

	if err := s.validateParent(ctx, rcv); err != nil {
		return err
	}

	receiverPlugin, err := s.getReceiverPlugin(rcv.Type)
	if err != nil {
		return err
	}

	rcv.Configurations, err = receiverPlugin.PreHookDBTransformConfigs(ctx, rcv.Configurations)
	if err != nil {
		return errors.ErrInvalid.WithMsgf("%s", err.Error())
	}

	err = s.repository.Create(ctx, rcv)
	if err != nil {
		return err
	}

	return nil
}

type getOpts struct {
	withData   bool
	withExpand bool
}

type GetOption func(*getOpts)

func GetWithData() GetOption {
	return func(g *getOpts) {
		g.withData = true
	}
}

func GetWithExpand() GetOption {
	return func(g *getOpts) {
		g.withExpand = true
	}
}

func (s *Service) Get(ctx context.Context, id uint64, gopts ...GetOption) (*Receiver, error) {
	opt := &getOpts{}

	for _, g := range gopts {
		g(opt)
	}

	filter := Filter{
		Expanded: opt.withExpand,
	}

	rcv, err := s.repository.Get(ctx, id, filter)
	if err != nil {
		if errors.As(err, new(NotFoundError)) {
			return nil, errors.ErrNotFound.WithMsgf(err.Error())
		}
		return nil, err
	}

	if !opt.withExpand {
		return rcv, nil
	}

	receiverPlugin, err := s.getReceiverPlugin(rcv.Type)
	if err != nil {
		return nil, err
	}

	transformedConfigs, err := receiverPlugin.PostHookDBTransformConfigs(ctx, rcv.Configurations)
	if err != nil {
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
	filter := Filter{
		Expanded: false,
	}
	oldReceiver, err := s.repository.Get(ctx, rcv.ID, filter)
	if err != nil {
		if errors.As(err, new(NotFoundError)) {
			return errors.ErrNotFound.WithMsgf(err.Error())
		}
		return err
	}

	receiverPlugin, err := s.getReceiverPlugin(oldReceiver.Type)
	if err != nil {
		return err
	}

	rcv.Configurations, err = receiverPlugin.PreHookDBTransformConfigs(ctx, rcv.Configurations)
	if err != nil {
		return errors.ErrInvalid.WithMsgf("%s", err.Error())
	}

	if err = s.repository.Update(ctx, rcv); err != nil {
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

func (s *Service) validateParent(ctx context.Context, rcv *Receiver) error {
	if rcv.ParentID == 0 {
		return nil
	}

	filter := Filter{
		Expanded: false,
	}

	parentRcv, err := s.repository.Get(ctx, rcv.ParentID, filter)
	if err != nil {
		return errors.ErrInvalid.WithMsgf("failed to check parent id %d", rcv.ParentID).WithCausef(err.Error())
	}

	switch rcv.Type {
	case TypeSlackChannel:
		if parentRcv.Type != TypeSlack {
			return errors.ErrInvalid.WithMsgf("parent of slack_channel type should be slack but found %s", parentRcv.Type)
		}
	default:
		return errors.ErrInvalid.WithMsgf("type %s should not have parent receiver", rcv.Type)
	}

	return nil
}
