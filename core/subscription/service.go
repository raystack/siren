package subscription

import (
	"context"
	"sort"

	"github.com/odpf/siren/core/namespace"
	"github.com/odpf/siren/core/receiver"
	"github.com/odpf/siren/pkg/errors"
)

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
	BuildNotificationConfig(subsConfig map[string]interface{}, rcv *receiver.Receiver) (map[string]interface{}, error)
}

// Service handles business logic
type Service struct {
	repository                   Repository
	namespaceService             NamespaceService
	receiverService              ReceiverService
	subscriptionProviderRegistry map[string]SubscriptionSyncer
}

// NewService returns service struct
func NewService(repository Repository, namespaceService NamespaceService, receiverService ReceiverService, sopts ...ServiceOption) *Service {
	svc := &Service{
		repository:                   repository,
		namespaceService:             namespaceService,
		receiverService:              receiverService,
		subscriptionProviderRegistry: map[string]SubscriptionSyncer{},
	}

	for _, opt := range sopts {
		opt(svc)
	}

	return svc
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

	subscriptionsInNamespace, err := s.FetchEnrichedSubscriptionsByNamespace(ctx, ns)
	if err != nil {
		if err := s.repository.Rollback(ctx, err); err != nil {
			return err
		}
		return err
	}

	if err := pluginService.CreateSubscription(ctx, sub, subscriptionsInNamespace, ns.URN); err != nil {
		if err := s.repository.Rollback(ctx, err); err != nil {
			return err
		}
		return err
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

	subscriptionsInNamespace, err := s.FetchEnrichedSubscriptionsByNamespace(ctx, ns)
	if err != nil {
		if err := s.repository.Rollback(ctx, err); err != nil {
			return err
		}
		return err
	}

	if err := pluginService.UpdateSubscription(ctx, sub, subscriptionsInNamespace, ns.URN); err != nil {
		if err := s.repository.Rollback(ctx, err); err != nil {
			return err
		}
		return err
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

	subscriptionsInNamespace, err := s.FetchEnrichedSubscriptionsByNamespace(ctx, ns)
	if err != nil {
		if err := s.repository.Rollback(ctx, err); err != nil {
			return err
		}
		return err
	}

	if err := pluginService.DeleteSubscription(ctx, sub, subscriptionsInNamespace, ns.URN); err != nil {
		if err := s.repository.Rollback(ctx, err); err != nil {
			return err
		}
		return err
	}

	if err := s.repository.Commit(ctx); err != nil {
		return err
	}

	return nil
}

// TODO we might want to add filter by namespace id too here
// to filter by tenant
func (s *Service) MatchByLabels(ctx context.Context, labels map[string]string) ([]Subscription, error) {
	// fetch all subscriptions by matching labels.
	subscriptionsByLabels, err := s.repository.List(ctx, Filter{
		Labels: labels,
	})
	if err != nil {
		return nil, err
	}

	if len(subscriptionsByLabels) == 0 {
		return nil, nil
	}

	receiversMap, err := CreateReceiversMap(ctx, s.receiverService, subscriptionsByLabels)
	if err != nil {
		return nil, err
	}

	subscriptionsByLabels, err = AssignReceivers(s.receiverService, receiversMap, subscriptionsByLabels)
	if err != nil {
		return nil, err
	}

	return subscriptionsByLabels, nil
}

func (s *Service) getProviderPluginService(providerType string) (SubscriptionSyncer, error) {
	pluginService, exist := s.subscriptionProviderRegistry[providerType]
	if !exist {
		return nil, errors.ErrInvalid.WithMsgf("unsupported provider type: %q", providerType)
	}
	return pluginService, nil
}

func (s *Service) FetchEnrichedSubscriptionsByNamespace(
	ctx context.Context,
	ns *namespace.Namespace) ([]Subscription, error) {

	// fetch all subscriptions in this namespace.
	subscriptionsInNamespace, err := s.repository.List(ctx, Filter{
		NamespaceID: ns.ID,
	})
	if err != nil {
		return nil, err
	}

	if len(subscriptionsInNamespace) == 0 {
		return subscriptionsInNamespace, nil
	}

	receiversMap, err := CreateReceiversMap(ctx, s.receiverService, subscriptionsInNamespace)
	if err != nil {
		return nil, err
	}

	subscriptionsInNamespace, err = AssignReceivers(s.receiverService, receiversMap, subscriptionsInNamespace)
	if err != nil {
		return nil, err
	}

	return subscriptionsInNamespace, nil
}

func CreateReceiversMap(ctx context.Context, receiverService ReceiverService, subscriptions []Subscription) (map[uint64]*receiver.Receiver, error) {
	receiversMap := map[uint64]*receiver.Receiver{}
	for _, subs := range subscriptions {
		for _, rcv := range subs.Receivers {
			if rcv.ID != 0 {
				receiversMap[rcv.ID] = nil
			}
		}
	}

	// empty receivers map
	if len(receiversMap) == 0 {
		return nil, errors.New("no receivers found in subscription")
	}

	listOfReceiverIDs := []uint64{}
	for k := range receiversMap {
		listOfReceiverIDs = append(listOfReceiverIDs, k)
	}

	filteredReceivers, err := receiverService.List(ctx, receiver.Filter{
		ReceiverIDs: listOfReceiverIDs,
	})
	if err != nil {
		return nil, err
	}

	for i, rcv := range filteredReceivers {
		receiversMap[rcv.ID] = &filteredReceivers[i]
	}

	nilReceivers := []uint64{}
	for id, rcv := range receiversMap {
		if rcv == nil {
			nilReceivers = append(nilReceivers, id)
			continue
		}
	}

	if len(nilReceivers) > 0 {
		return nil, errors.ErrInvalid.WithMsgf("receiver id %v don't exist", nilReceivers)
	}

	return receiversMap, nil
}

func AssignReceivers(receiverService ReceiverService, receiversMap map[uint64]*receiver.Receiver, subscriptions []Subscription) ([]Subscription, error) {
	for is := range subscriptions {
		for ir, subsRcv := range subscriptions[is].Receivers {
			if mappedRcv := receiversMap[subsRcv.ID]; mappedRcv == nil {
				return nil, errors.ErrInvalid.WithMsgf("receiver id %d not found", subsRcv.ID)
			}
			subsConfig, err := receiverService.BuildNotificationConfig(subsRcv.Configuration, receiversMap[subsRcv.ID])
			if err != nil {
				return nil, errors.ErrInvalid.WithMsgf(err.Error())
			}
			subscriptions[is].Receivers[ir].ID = receiversMap[subsRcv.ID].ID
			subscriptions[is].Receivers[ir].Type = receiversMap[subsRcv.ID].Type
			subscriptions[is].Receivers[ir].Configuration = subsConfig
		}
	}

	return subscriptions, nil
}

func sortReceivers(sub *Subscription) {
	sort.Slice(sub.Receivers, func(i, j int) bool {
		return sub.Receivers[i].ID < sub.Receivers[j].ID
	})
}
