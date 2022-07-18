package subscription

import (
	"context"
	"fmt"
	"sort"

	"github.com/odpf/siren/core/namespace"
	"github.com/odpf/siren/core/provider"
	"github.com/odpf/siren/core/receiver"
	"github.com/odpf/siren/pkg/cortex"
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
	Notify(ctx context.Context, id uint64, payloadMessage receiver.NotificationMessage) error
	GetSubscriptionConfig(subsConfs map[string]string, rcv *receiver.Receiver) (map[string]string, error)
}

//go:generate mockery --name=ProviderService -r --case underscore --with-expecter --structname ProviderService --filename provider_service.go --output=./mocks
type ProviderService interface {
	List(context.Context, provider.Filter) ([]provider.Provider, error)
	Create(context.Context, *provider.Provider) error
	Get(context.Context, uint64) (*provider.Provider, error)
	Update(context.Context, *provider.Provider) error
	Delete(context.Context, uint64) error
}

// Service handles business logic
type Service struct {
	repository       Repository
	providerService  ProviderService
	namespaceService NamespaceService
	receiverService  ReceiverService
	cortexClient     CortexClient
}

// NewService returns service struct
func NewService(repository Repository, providerService ProviderService, namespaceService NamespaceService,
	receiverService ReceiverService, cortexClient CortexClient) *Service {

	return &Service{
		repository:       repository,
		providerService:  providerService,
		namespaceService: namespaceService,
		receiverService:  receiverService,
		cortexClient:     cortexClient,
	}
}

func (s Service) List(ctx context.Context, flt Filter) ([]Subscription, error) {
	subscriptions, err := s.repository.List(ctx, flt)
	if err != nil {
		return nil, err
	}

	return subscriptions, nil
}

func (s Service) Create(ctx context.Context, sub *Subscription) error {
	// check provider type of the namespace
	ns, err := s.namespaceService.Get(ctx, sub.Namespace)
	if err != nil {
		return err
	}

	prov, err := s.providerService.Get(ctx, ns.Provider)
	if err != nil {
		return err
	}

	sortReceivers(sub)

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

	if err = s.SyncToUpstream(ctx, ns, prov); err != nil {
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

func (s Service) Get(ctx context.Context, id uint64) (*Subscription, error) {
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
	prov, err := s.providerService.Get(ctx, ns.Provider)
	if err != nil {
		return err
	}

	sortReceivers(sub)

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

	if err = s.SyncToUpstream(ctx, ns, prov); err != nil {
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

func (s Service) Delete(ctx context.Context, id uint64) error {
	sub, err := s.repository.Get(ctx, id)
	if err != nil {
		return err
	}
	// check provider type of the namespace
	ns, err := s.namespaceService.Get(ctx, sub.Namespace)
	if err != nil {
		return err
	}
	prov, err := s.providerService.Get(ctx, ns.Provider)
	if err != nil {
		return err
	}

	ctx = s.repository.WithTransaction(ctx)
	if err := s.repository.Delete(ctx, id, sub.Namespace); err != nil {
		if err := s.repository.Rollback(ctx, err); err != nil {
			return err
		}
		return err
	}

	if err = s.SyncToUpstream(ctx, ns, prov); err != nil {
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

func (s *Service) SyncToUpstream(
	ctx context.Context,
	ns *namespace.Namespace,
	prov *provider.Provider) error {

	// fetch all subscriptions in this namespace.
	subscriptionsInNamespace, err := s.repository.List(ctx, Filter{
		NamespaceID: ns.ID,
	})
	if err != nil {
		return err
	}

	receiversMap, err := CreateReceiversMap(ctx, s.receiverService, subscriptionsInNamespace)
	if err != nil {
		return err
	}

	subscriptions, err := AssignReceivers(s.receiverService, receiversMap, subscriptionsInNamespace)
	if err != nil {
		return err
	}

	// do upstream call to create subscriptions as per provider type
	switch prov.Type {
	case "cortex":
		amConfig := make([]cortex.ReceiverConfig, 0)
		for _, item := range subscriptions {
			amConfig = append(amConfig, item.ToAlertManagerReceiverConfig()...)
		}

		err = s.cortexClient.CreateAlertmanagerConfig(cortex.AlertManagerConfig{
			Receivers: amConfig,
		}, ns.URN)
		if err != nil {
			return fmt.Errorf("error calling cortex: %w", err)
		}
	default:
		return errors.New(fmt.Sprintf("subscriptions for provider type '%s' not supported", prov.Type))
	}
	return nil
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
			subsConfig, err := receiverService.GetSubscriptionConfig(subsRcv.Configuration, receiversMap[subsRcv.ID])
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
