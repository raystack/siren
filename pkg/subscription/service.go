package subscription

import (
	"time"

	"github.com/odpf/siren/domain"
	"github.com/odpf/siren/pkg/namespace"
	"github.com/odpf/siren/pkg/provider"
	"github.com/odpf/siren/pkg/receiver"
	"github.com/odpf/siren/store"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

// Service handles business logic
type Service struct {
	repository       SubscriptionRepository
	providerService  domain.ProviderService
	namespaceService domain.NamespaceService
	receiverService  domain.ReceiverService
}

// NewService returns service struct
func NewService(providerRepository store.ProviderRepository, namespaceRepository store.NamespaceRepository,
	receiverRepository store.ReceiverRepository, db *gorm.DB, key string) (domain.SubscriptionService, error) {
	repository := NewRepository(db)
	namespaceService, err := namespace.NewService(namespaceRepository, key)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create namespace service")
	}
	receiverService, err := receiver.NewService(receiverRepository, nil, key)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create receiver service")
	}
	return &Service{repository, provider.NewService(providerRepository),
		namespaceService, receiverService}, nil
}

func (s Service) ListSubscriptions() ([]*domain.Subscription, error) {
	subscriptions, err := s.repository.List()
	if err != nil {
		return nil, errors.Wrap(err, "s.repository.List")
	}

	return subscriptions, nil
}

func (s Service) CreateSubscription(sub *domain.Subscription) (*domain.Subscription, error) {
	newSubscription, err := s.repository.Create(sub, s.namespaceService, s.providerService, s.receiverService)
	if err != nil {
		return nil, errors.Wrap(err, "s.repository.Create")
	}
	return newSubscription, nil
}

func (s Service) GetSubscription(id uint64) (*domain.Subscription, error) {
	subscription, err := s.repository.Get(id)
	if err != nil {
		return nil, errors.Wrap(err, "s.repository.Get")
	}
	if subscription == nil {
		return nil, nil
	}
	return subscription, nil
}

func (s Service) UpdateSubscription(sub *domain.Subscription) (*domain.Subscription, error) {
	updatedSubscription, err := s.repository.Update(sub, s.namespaceService, s.providerService, s.receiverService)
	if err != nil {
		return nil, errors.Wrap(err, "s.repository.Update")
	}
	return updatedSubscription, nil
}

func (s Service) DeleteSubscription(id uint64) error {
	return s.repository.Delete(id, s.namespaceService, s.providerService, s.receiverService)
}

func (s Service) Migrate() error {
	return s.repository.Migrate()
}

type AuditLog struct {
	Timestamp time.Time
	Action    string // example: appeal.created, provider.created, provider.updated, etc.
	Actor     string // example: user@example.com or system
	Data      interface{}
	Message   string
}

type AuditRepository interface {
	List(filters map[string]interface{}) ([]*AuditLog, error)
	Create(*AuditLog) error
	BulkCreate([]*AuditLog) error
}

type AuditAction string

var (
	ProviderCreated     AuditAction = "provider.created"
	ResourceBulkCreated AuditAction = "resource.bulkCreated"

	// ...
)

type ProviderCreatedData domain.Provider
type ResourceBulkCreatedData struct {
	CreatedResourceIDs []string
	RemovedResourceIDs []string
}

type AuditService interface {
	Log(...*AuditLog) error
}
