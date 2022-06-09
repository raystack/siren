package store

import (
	"context"

	"github.com/odpf/siren/core/alert"
	"github.com/odpf/siren/core/namespace"
	"github.com/odpf/siren/core/provider"
	"github.com/odpf/siren/core/receiver"
	"github.com/odpf/siren/core/rule"
	"github.com/odpf/siren/core/subscription"
	"github.com/odpf/siren/core/template"
	"github.com/odpf/siren/internal/store/postgres"
	"gorm.io/gorm"
)

type Transactor interface {
	WithTransaction(ctx context.Context) context.Context
	Rollback(ctx context.Context) error
	Commit(ctx context.Context) error
}

type RepositoryContainer struct {
	ProviderRepository     provider.Repository
	NamespaceRepository    namespace.Repository
	TemplatesRepository    template.Repository
	ReceiverRepository     receiver.Repository
	SubscriptionRepository subscription.Repository
	AlertRepository        alert.Repository
	RuleRepository         rule.Repository
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
