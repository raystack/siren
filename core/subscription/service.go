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
	List(context.Context) ([]*namespace.Namespace, error)
	Create(context.Context, *namespace.Namespace) (uint64, error)
	Get(context.Context, uint64) (*namespace.Namespace, error)
	Update(context.Context, *namespace.Namespace) (uint64, error)
	Delete(context.Context, uint64) error
}

//go:generate mockery --name=ReceiverService -r --case underscore --with-expecter --structname ReceiverService --filename receiver_service.go --output=./mocks
type ReceiverService interface {
	List(context.Context) ([]*receiver.Receiver, error)
	Create(context.Context, *receiver.Receiver) (uint64, error)
	Get(context.Context, uint64) (*receiver.Receiver, error)
	Update(context.Context, *receiver.Receiver) (uint64, error)
	Delete(context.Context, uint64) error
	Notify(context.Context, uint64, receiver.NotificationMessage) error
}

//go:generate mockery --name=ProviderService -r --case underscore --with-expecter --structname ProviderService --filename provider_service.go --output=./mocks
type ProviderService interface {
	List(context.Context, provider.Filter) ([]*provider.Provider, error)
	Create(context.Context, *provider.Provider) (uint64, error)
	Get(context.Context, uint64) (*provider.Provider, error)
	Update(context.Context, *provider.Provider) (uint64, error)
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

func (s Service) ListSubscriptions(ctx context.Context) ([]*Subscription, error) {
	subscriptions, err := s.repository.List(ctx)
	if err != nil {
		return nil, err
	}

	return subscriptions, nil
}

func (s Service) CreateSubscription(ctx context.Context, sub *Subscription) error {
	ctx = s.repository.WithTransaction(ctx)
	sortReceivers(sub)
	if err := s.repository.Create(ctx, sub); err != nil {
		if err := s.repository.Rollback(ctx); err != nil {
			return err
		}
		if errors.Is(err, ErrDuplicate) {
			return errors.ErrConflict.WithMsgf(err.Error())
		}
		return err
	}

	if err := s.syncInUpstreamCurrentSubscriptionsOfNamespace(ctx, sub.Namespace); err != nil {
		if err := s.repository.Rollback(ctx); err != nil {
			return err
		}
		return err
	}

	if err := s.repository.Commit(ctx); err != nil {
		return err
	}
	return nil
}

func (s Service) GetSubscription(ctx context.Context, id uint64) (*Subscription, error) {
	subscription, err := s.repository.Get(ctx, id)
	if err != nil {
		if errors.As(err, new(NotFoundError)) {
			return nil, errors.ErrNotFound.WithMsgf(err.Error())
		}
		return nil, err
	}

	return subscription, nil
}

func (s Service) UpdateSubscription(ctx context.Context, sub *Subscription) error {
	ctx = s.repository.WithTransaction(ctx)
	sortReceivers(sub)
	if err := s.repository.Update(ctx, sub); err != nil {
		if err := s.repository.Rollback(ctx); err != nil {
			return err
		}
		if errors.Is(err, ErrDuplicate) {
			return errors.ErrConflict.WithMsgf(err.Error())
		}
		if errors.As(err, new(NotFoundError)) {
			return errors.ErrNotFound.WithMsgf(err.Error())
		}
		return err
	}

	if err := s.syncInUpstreamCurrentSubscriptionsOfNamespace(ctx, sub.Namespace); err != nil {
		fmt.Printf("err: %v\n", err)
		if err := s.repository.Rollback(ctx); err != nil {
			return err
		}
		return err
	}

	if err := s.repository.Commit(ctx); err != nil {
		return err
	}
	return nil
}

func (s Service) DeleteSubscription(ctx context.Context, id uint64) error {
	sub, err := s.repository.Get(ctx, id)
	if err != nil {
		return err
	}

	ctx = s.repository.WithTransaction(ctx)
	if err := s.repository.Delete(ctx, id); err != nil {
		if err := s.repository.Rollback(ctx); err != nil {
			return err
		}
		return err
	}

	if err := s.syncInUpstreamCurrentSubscriptionsOfNamespace(ctx, sub.Namespace); err != nil {
		if err := s.repository.Rollback(ctx); err != nil {
			return err
		}
		return err
	}
	return nil
}

func (s Service) syncInUpstreamCurrentSubscriptionsOfNamespace(ctx context.Context, namespaceId uint64) error {
	// fetch all subscriptions in this namespace.
	subscriptionsInNamespace, err := s.getAllSubscriptionsWithinNamespace(ctx, namespaceId)
	if err != nil {
		return err
	}
	// check provider type of the namespace
	providerInfo, namespaceInfo, err := s.getProviderAndNamespaceInfoFromNamespaceId(ctx, namespaceId)
	if err != nil {
		return err
	}
	subscriptionsInNamespaceEnrichedWithReceivers, err := s.addReceiversConfiguration(ctx, subscriptionsInNamespace)
	if err != nil {
		return err
	}
	// do upstream call to create subscriptions as per provider type
	switch providerInfo.Type {
	case "cortex":
		amConfig := getAmConfigFromSubscriptions(subscriptionsInNamespaceEnrichedWithReceivers)
		err = s.cortexClient.CreateAlertmanagerConfig(amConfig, namespaceInfo.URN)
		if err != nil {
			return err
		}
	default:
		return errors.New(fmt.Sprintf("subscriptions for provider type '%s' not supported", providerInfo.Type))
	}
	return nil
}

//TODO this can use repository filter by namespace id
func (s Service) getAllSubscriptionsWithinNamespace(ctx context.Context, id uint64) ([]*Subscription, error) {
	subscriptions, err := s.repository.List(ctx)
	if err != nil {
		return nil, err
	}
	var subscriptionsWithinNamespace []*Subscription
	for _, sub := range subscriptions {
		if sub.Namespace == id {
			subscriptionsWithinNamespace = append(subscriptionsWithinNamespace, sub)
		}
	}
	return subscriptionsWithinNamespace, nil
}

func (s Service) getProviderAndNamespaceInfoFromNamespaceId(ctx context.Context, id uint64) (*provider.Provider, *namespace.Namespace, error) {
	namespaceInfo, err := s.namespaceService.Get(ctx, id)
	if err != nil {
		return nil, nil, err
	}
	providerInfo, err := s.providerService.Get(ctx, namespaceInfo.Provider)
	if err != nil {
		return nil, nil, err
	}
	return providerInfo, namespaceInfo, nil
}

func (s Service) addReceiversConfiguration(ctx context.Context, subscriptions []*Subscription) ([]SubscriptionEnrichedWithReceivers, error) {
	res := make([]SubscriptionEnrichedWithReceivers, 0)
	allReceivers, err := s.receiverService.List(ctx)
	if err != nil {
		return nil, err
	}
	for _, item := range subscriptions {
		enrichedReceivers := make([]EnrichedReceiverMetadata, 0)
		for _, receiverItem := range item.Receivers {
			var receiverInfo *receiver.Receiver
			for idx := range allReceivers {
				if allReceivers[idx].ID == receiverItem.ID {
					receiverInfo = allReceivers[idx]
					break
				}
			}
			if receiverInfo == nil {
				return nil, errors.New(fmt.Sprintf("receiver id %d does not exist", receiverItem.ID))
			}
			//initialize the nil map using the make function
			//to avoid panics while adding elements in future
			if receiverItem.Configuration == nil {
				receiverItem.Configuration = make(map[string]string)
			}
			switch receiverInfo.Type {
			case "slack":
				if _, ok := receiverItem.Configuration["channel_name"]; !ok {
					return nil, errors.New(fmt.Sprintf(
						"configuration.channel_name missing from receiver with id %d", receiverItem.ID))
				}
				if val, ok := receiverInfo.Configurations["token"]; ok {
					receiverItem.Configuration["token"] = val.(string)
				}
			case "pagerduty":
				if val, ok := receiverInfo.Configurations["service_key"]; ok {
					receiverItem.Configuration["service_key"] = val.(string)
				}
			case "http":
				if val, ok := receiverInfo.Configurations["url"]; ok {
					receiverItem.Configuration["url"] = val.(string)
				}
			default:
				return nil, errors.New(fmt.Sprintf(`subscriptions for receiver type %s not supported via Siren inside Cortex`, receiverInfo.Type))
			}
			enrichedReceiver := EnrichedReceiverMetadata{
				ID:            receiverItem.ID,
				Configuration: receiverItem.Configuration,
				Type:          receiverInfo.Type,
			}
			enrichedReceivers = append(enrichedReceivers, enrichedReceiver)
		}
		enrichedSubscription := SubscriptionEnrichedWithReceivers{
			ID:          item.ID,
			NamespaceId: item.Namespace,
			URN:         item.URN,
			Receiver:    enrichedReceivers,
			Match:       item.Match,
		}
		res = append(res, enrichedSubscription)
	}
	return res, nil
}

func sortReceivers(sub *Subscription) {
	sort.Slice(sub.Receivers, func(i, j int) bool {
		return sub.Receivers[i].ID < sub.Receivers[j].ID
	})
}

func getAMReceiverConfigPerSubscription(sub SubscriptionEnrichedWithReceivers) []cortex.ReceiverConfig {
	amReceiverConfig := make([]cortex.ReceiverConfig, 0)
	for idx, item := range sub.Receiver {
		newAMReceiver := cortex.ReceiverConfig{
			Receiver:      fmt.Sprintf("%s_receiverId_%d_idx_%d", sub.URN, item.ID, idx),
			Match:         sub.Match,
			Configuration: item.Configuration,
			Type:          item.Type,
		}
		amReceiverConfig = append(amReceiverConfig, newAMReceiver)
	}
	return amReceiverConfig
}

func getAmConfigFromSubscriptions(subscriptions []SubscriptionEnrichedWithReceivers) cortex.AlertManagerConfig {
	amConfig := make([]cortex.ReceiverConfig, 0)
	for _, item := range subscriptions {
		amConfig = append(amConfig, getAMReceiverConfigPerSubscription(item)...)
	}
	return cortex.AlertManagerConfig{
		Receivers: amConfig,
	}
}
