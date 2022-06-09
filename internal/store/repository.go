package store

import (
	"context"

	"github.com/odpf/siren/core/alert"
	"github.com/odpf/siren/core/namespace"
	"github.com/odpf/siren/core/provider"
	"github.com/odpf/siren/core/receiver"
	"github.com/odpf/siren/domain"
	"github.com/odpf/siren/internal/store/postgres"
	"gorm.io/gorm"
)

type TemplatesRepository interface {
	Upsert(*domain.Template) error
	Index(string) ([]domain.Template, error)
	GetByName(string) (*domain.Template, error)
	Delete(string) error
	Render(string, map[string]string) (string, error)
	Migrate() error
}

type SubscriptionRepository interface {
	Transactor
	Migrate() error
	List(context.Context) ([]*domain.Subscription, error)
	Create(context.Context, *domain.Subscription) error
	Get(context.Context, uint64) (*domain.Subscription, error)
	Update(context.Context, *domain.Subscription) error
	Delete(context.Context, uint64) error
}

type RuleRepository interface {
	Transactor
	Upsert(context.Context, *domain.Rule) error
	Get(context.Context, string, string, string, string, uint64) ([]domain.Rule, error)
	Migrate() error
}

type Transactor interface {
	WithTransaction(ctx context.Context) context.Context
	Rollback(ctx context.Context) error
	Commit(ctx context.Context) error
}

type RepositoryContainer struct {
	ProviderRepository     provider.Repository
	NamespaceRepository    namespace.Repository
	TemplatesRepository    TemplatesRepository
	ReceiverRepository     receiver.Repository
	SubscriptionRepository SubscriptionRepository
	AlertRepository        alert.Repository
	RuleRepository         RuleRepository
}

func NewRepositoryContainer(db *gorm.DB) *RepositoryContainer {
	return &RepositoryContainer{
		NamespaceRepository:    postgres.NewNamespaceRepository(db),
		ProviderRepository:     postgres.NewProviderRepository(db),
		ReceiverRepository:     postgres.NewReceiverRepository(db),
		TemplatesRepository:    postgres.NewTemplateRepository(db),
		SubscriptionRepository: postgres.NewSubscriptionRepository(db),
		AlertRepository:        postgres.NewAlertRepository(db),
		RuleRepository:         postgres.NewRuleRepository(db),
	}
}
