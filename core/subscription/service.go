package subscription

import (
	"context"

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
}

// Service handles business logic
type Service struct {
	repository       Repository
	namespaceService NamespaceService
	receiverService  ReceiverService
}

// NewService returns service struct
func NewService(repository Repository, namespaceService NamespaceService, receiverService ReceiverService) *Service {
	svc := &Service{
		repository:       repository,
		namespaceService: namespaceService,
		receiverService:  receiverService,
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
	if err := s.repository.Create(ctx, sub); err != nil {
		if errors.Is(err, ErrDuplicate) {
			return errors.ErrConflict.WithMsgf(err.Error())
		}
		if errors.Is(err, ErrRelation) {
			return errors.ErrNotFound.WithMsgf(err.Error())
		}
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
	if err := s.repository.Update(ctx, sub); err != nil {
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

	return nil
}

func (s *Service) Delete(ctx context.Context, id uint64) error {
	if err := s.repository.Delete(ctx, id); err != nil {
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

	subscriptionsByLabels, err = AssignReceivers(receiversMap, subscriptionsByLabels)
	if err != nil {
		return nil, err
	}

	return subscriptionsByLabels, nil
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

func AssignReceivers(receiversMap map[uint64]*receiver.Receiver, subscriptions []Subscription) ([]Subscription, error) {
	for is := range subscriptions {
		for ir, subsRcv := range subscriptions[is].Receivers {
			if mappedRcv := receiversMap[subsRcv.ID]; mappedRcv == nil {
				return nil, errors.ErrInvalid.WithMsgf("receiver id %d not found", subsRcv.ID)
			}
			mergedConfigMap := MergeConfigsMap(subsRcv.Configuration, receiversMap[subsRcv.ID].Configurations)
			subscriptions[is].Receivers[ir].ID = receiversMap[subsRcv.ID].ID
			subscriptions[is].Receivers[ir].Type = receiversMap[subsRcv.ID].Type
			subscriptions[is].Receivers[ir].Configuration = mergedConfigMap
		}
	}

	return subscriptions, nil
}

func MergeConfigsMap(subscriptionConfigMap map[string]interface{}, receiverConfigsMap map[string]interface{}) map[string]interface{} {
	var newConfigMap = make(map[string]interface{})
	for k, v := range subscriptionConfigMap {
		newConfigMap[k] = v
	}
	for k, v := range receiverConfigsMap {
		newConfigMap[k] = v
	}
	return newConfigMap
}
