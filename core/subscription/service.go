package subscription

import (
	"context"
	"sort"

	"github.com/odpf/siren/core/namespace"
	"github.com/odpf/siren/core/provider"
	"github.com/odpf/siren/core/receiver"
	"github.com/odpf/siren/pkg/errors"
)

//go:generate mockery --name=ProviderPlugin -r --case underscore --with-expecter --structname ProviderPlugin --filename plugin_provider.go --output=./mocks
type ProviderPlugin interface {
	CreateSubscription(ctx context.Context, sub *Subscription, namespaceURN string) error
	UpdateSubscription(ctx context.Context, sub *Subscription, namespaceURN string) error
	DeleteSubscription(ctx context.Context, sub *Subscription, namespaceURN string) error
	SyncSubscriptions(ctx context.Context, subs []Subscription, namespaceURN string) error
	SyncMethod() provider.SyncMethod
}

//go:generate mockery --name=NamespaceService -r --case underscore --with-expecter --structname NamespaceService --filename namespace_service.go --output=./mocks
type NamespaceService interface {
	List(context.Context) ([]namespace.Namespace, error)
	Create(context.Context, *namespace.Namespace) error
	Get(context.Context, uint64) (*namespace.Namespace, error)
	Update(context.Context, *namespace.Namespace) error
	Delete(context.Context, uint64) error
}

//go:generate mockery --name=ReceiverService -r --case underscore --with-expecter --structname ReceiverService --filename receiver_service.go --output=./mocks
type ReceiverService interface {
	List(ctx context.Context, flt receiver.Filter) ([]receiver.Receiver, error)
	Create(ctx context.Context, rcv *receiver.Receiver) error
	Get(ctx context.Context, id uint64) (*receiver.Receiver, error)
	Update(ctx context.Context, rcv *receiver.Receiver) error
	Delete(ctx context.Context, id uint64) error
	Notify(ctx context.Context, id uint64, payloadMessage map[string]interface{}) error
	EnrichSubscriptionConfig(subsConfig map[string]string, rcv *receiver.Receiver) (map[string]string, error)
}

// Service handles business logic
type Service struct {
	repository                   Repository
	namespaceService             NamespaceService
	receiverService              ReceiverService
	subscriptionProviderRegistry map[string]ProviderPlugin
}

// NewService returns service struct
func NewService(repository Repository, namespaceService NamespaceService, receiverService ReceiverService, subscriptionProviderRegistry map[string]ProviderPlugin) *Service {
	return &Service{
		repository:                   repository,
		namespaceService:             namespaceService,
		receiverService:              receiverService,
		subscriptionProviderRegistry: subscriptionProviderRegistry,
	}
}

func (s *Service) List(ctx context.Context, flt Filter) ([]Subscription, error) {
	subscriptions, err := s.repository.List(ctx, flt)
	if err != nil {
		return nil, err
	}

	return subscriptions, nil
}

func (s *Service) Create(ctx context.Context, sub *Subscription) error {
	// check provider type of the namespace
	ns, err := s.namespaceService.Get(ctx, sub.Namespace)
	if err != nil {
		return err
	}

	sortReceivers(sub)

	pluginService, err := s.getProviderPluginService(ns.Provider.Type)
	if err != nil {
		return err
	}

	ctx = s.repository.WithTransaction(ctx)
	if err = s.repository.Create(ctx, sub); err != nil {
		if err := s.repository.Rollback(ctx, err); err != nil {
			return err
		}
		if errors.Is(err, ErrDuplicate) {
			return errors.ErrConflict.WithMsgf(err.Error())
		}
		if errors.Is(err, ErrRelation) {
			return errors.ErrNotFound.WithMsgf(err.Error())
		}
		return err
	}

	switch pluginService.SyncMethod() {
	case provider.TypeSyncSingle:
		if err := pluginService.CreateSubscription(ctx, sub, ns.URN); err != nil {
			if err := s.repository.Rollback(ctx, err); err != nil {
				return err
			}
			return err
		}
	case provider.TypeSyncBatch:
	default:
		if err := s.SyncBatchToUpstream(ctx, ns, pluginService); err != nil {
			if err := s.repository.Rollback(ctx, err); err != nil {
				return err
			}
			return err
		}
	}

	if err := s.repository.Commit(ctx); err != nil {
		return err
	}
	return nil
}

func (s *Service) Get(ctx context.Context, id uint64) (*Subscription, error) {
	subscription, err := s.repository.Get(ctx, id)
	if err != nil {
		if errors.As(err, new(NotFoundError)) {
			return nil, errors.ErrNotFound.WithMsgf(err.Error())
		}
		return nil, err
	}

	return subscription, nil
}

func (s *Service) Update(ctx context.Context, sub *Subscription) error {
	// check provider type of the namespace
	ns, err := s.namespaceService.Get(ctx, sub.Namespace)
	if err != nil {
		return err
	}

	sortReceivers(sub)

	pluginService, err := s.getProviderPluginService(ns.Provider.Type)
	if err != nil {
		return err
	}

	ctx = s.repository.WithTransaction(ctx)
	if err = s.repository.Update(ctx, sub); err != nil {
		if err := s.repository.Rollback(ctx, err); err != nil {
			return err
		}
		if errors.Is(err, ErrDuplicate) {
			return errors.ErrConflict.WithMsgf(err.Error())
		}
		if errors.Is(err, ErrRelation) {
			return errors.ErrNotFound.WithMsgf(err.Error())
		}
		if errors.As(err, new(NotFoundError)) {
			return errors.ErrNotFound.WithMsgf(err.Error())
		}
		return err
	}

	switch pluginService.SyncMethod() {
	case provider.TypeSyncSingle:
		if err := pluginService.UpdateSubscription(ctx, sub, ns.URN); err != nil {
			if err := s.repository.Rollback(ctx, err); err != nil {
				return err
			}
			return err
		}
	case provider.TypeSyncBatch:
	default:
		if err := s.SyncBatchToUpstream(ctx, ns, pluginService); err != nil {
			if err := s.repository.Rollback(ctx, err); err != nil {
				return err
			}
			return err
		}
	}
	if err := s.repository.Commit(ctx); err != nil {
		return err
	}
	return nil
}

func (s *Service) Delete(ctx context.Context, id uint64) error {
	sub, err := s.repository.Get(ctx, id)
	if err != nil {
		return err
	}

	// check provider type of the namespace
	ns, err := s.namespaceService.Get(ctx, sub.Namespace)
	if err != nil {
		return err
	}

	pluginService, err := s.getProviderPluginService(ns.Provider.Type)
	if err != nil {
		return err
	}

	ctx = s.repository.WithTransaction(ctx)
	if err := s.repository.Delete(ctx, id); err != nil {
		if err := s.repository.Rollback(ctx, err); err != nil {
			return err
		}
		return err
	}

	switch pluginService.SyncMethod() {
	case provider.TypeSyncSingle:
		if err := pluginService.DeleteSubscription(ctx, sub, ns.URN); err != nil {
			if err := s.repository.Rollback(ctx, err); err != nil {
				return err
			}
			return err
		}
	case provider.TypeSyncBatch:
	default:
		if err := s.SyncBatchToUpstream(ctx, ns, pluginService); err != nil {
			if err := s.repository.Rollback(ctx, err); err != nil {
				return err
			}
			return err
		}
	}

	if err := s.repository.Commit(ctx); err != nil {
		return err
	}

	return nil
}

func (s *Service) getProviderPluginService(providerType string) (ProviderPlugin, error) {
	pluginService, exist := s.subscriptionProviderRegistry[providerType]
	if !exist {
		return nil, errors.ErrInvalid.WithMsgf("unsupported provider type: %q", providerType)
	}
	return pluginService, nil
}

func sortReceivers(sub *Subscription) {
	sort.Slice(sub.Receivers, func(i, j int) bool {
		return sub.Receivers[i].ID < sub.Receivers[j].ID
	})
}
