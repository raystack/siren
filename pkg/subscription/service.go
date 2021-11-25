package subscription

import (
	"github.com/odpf/siren/domain"
	"github.com/odpf/siren/pkg/namespace"
	"github.com/odpf/siren/pkg/provider"
	"github.com/odpf/siren/pkg/receiver"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

// Service handles business logic
type Service struct {
	repository      SubscriptionRepository
	providerService domain.ProviderService
	namespaceService domain.NamespaceService
	receiverService  domain.ReceiverService
}

// NewService returns service struct
func NewService(db *gorm.DB, key string) (domain.SubscriptionService, error) {
	repository := NewRepository(db)
	namespaceService, err := namespace.NewService(db, key)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create namespace service")
	}
	receiverService, err := receiver.NewService(db, nil, key)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create receiver service")
	}
	return &Service{repository, provider.NewService(db),
		namespaceService, receiverService}, nil
}

func (s Service) ListSubscriptions() ([]*domain.Subscription, error) {
	subscriptions, err := s.repository.List()
	if err != nil {
		return nil, errors.Wrap(err, "s.repository.List")
	}
	domainSubscriptions := make([]*domain.Subscription, 0, len(subscriptions))
	for i := 0; i < len(subscriptions); i++ {
		domainSubscriptions = append(domainSubscriptions, subscriptions[i].toDomain())
	}

	return domainSubscriptions, nil
}

func (s Service) CreateSubscription(domainSubscription *domain.Subscription) (*domain.Subscription, error) {
	sub := &Subscription{}
	sub.fromDomain(domainSubscription)
	newSubscription, err := s.repository.Create(sub, s.namespaceService, s.providerService, s.receiverService)
	if err != nil {
		return nil, errors.Wrap(err, "s.repository.Create")
	}
	return newSubscription.toDomain(), nil
}

func (s Service) GetSubscription(id uint64) (*domain.Subscription, error) {
	subscription, err := s.repository.Get(id)
	if err != nil {
		return nil, errors.Wrap(err, "s.repository.Get")
	}
	if subscription == nil {
		return nil, nil
	}
	return subscription.toDomain(), nil
}

func (s Service) UpdateSubscription(domainSubscription *domain.Subscription) (*domain.Subscription, error) {
	subscription := &Subscription{}
	subscription.fromDomain(domainSubscription)
	updatedSubscription, err := s.repository.Update(subscription, s.namespaceService, s.providerService, s.receiverService)
	if err != nil {
		return nil, errors.Wrap(err, "s.repository.Update")
	}
	return updatedSubscription.toDomain(), nil
}

func (s Service) DeleteSubscription(id uint64) error {
	return s.repository.Delete(id, s.namespaceService, s.providerService, s.receiverService)
}

func (s Service) Migrate() error {
	return s.repository.Migrate()
}
