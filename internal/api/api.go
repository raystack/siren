package api

import (
	"context"
	"time"

	"github.com/goto/siren/core/alert"
	"github.com/goto/siren/core/namespace"
	"github.com/goto/siren/core/notification"
	"github.com/goto/siren/core/provider"
	"github.com/goto/siren/core/receiver"
	"github.com/goto/siren/core/rule"
	"github.com/goto/siren/core/silence"
	"github.com/goto/siren/core/subscription"
	"github.com/goto/siren/core/template"
)

type AlertService interface {
	CreateAlerts(ctx context.Context, providerType string, providerID uint64, namespaceID uint64, body map[string]any) ([]alert.Alert, int, error)
	List(context.Context, alert.Filter) ([]alert.Alert, error)
}

type NamespaceService interface {
	List(context.Context) ([]namespace.Namespace, error)
	Create(context.Context, *namespace.Namespace) error
	Get(context.Context, uint64) (*namespace.Namespace, error)
	Update(context.Context, *namespace.Namespace) error
	Delete(context.Context, uint64) error
}

type ProviderService interface {
	List(context.Context, provider.Filter) ([]provider.Provider, error)
	Create(context.Context, *provider.Provider) error
	Get(context.Context, uint64) (*provider.Provider, error)
	Update(context.Context, *provider.Provider) error
	Delete(context.Context, uint64) error
}

type ReceiverService interface {
	List(ctx context.Context, flt receiver.Filter) ([]receiver.Receiver, error)
	Create(ctx context.Context, rcv *receiver.Receiver) error
	Get(ctx context.Context, id uint64, gopts ...receiver.GetOption) (*receiver.Receiver, error)
	Update(ctx context.Context, rcv *receiver.Receiver) error
	Delete(ctx context.Context, id uint64) error
}

type RuleService interface {
	Upsert(context.Context, *rule.Rule) error
	List(context.Context, rule.Filter) ([]rule.Rule, error)
}

type SubscriptionService interface {
	List(context.Context, subscription.Filter) ([]subscription.Subscription, error)
	Create(context.Context, *subscription.Subscription) error
	Get(context.Context, uint64) (*subscription.Subscription, error)
	Update(context.Context, *subscription.Subscription) error
	Delete(context.Context, uint64) error
}

type TemplateService interface {
	Upsert(context.Context, *template.Template) error
	List(context.Context, template.Filter) ([]template.Template, error)
	GetByName(context.Context, string) (*template.Template, error)
	Delete(context.Context, string) error
	Render(context.Context, string, map[string]string) (string, error)
}

type NotificationService interface {
	Dispatch(ctx context.Context, n notification.Notification) error
	CheckAndInsertIdempotency(ctx context.Context, scope, key string) (uint64, error)
	MarkIdempotencyAsSuccess(ctx context.Context, id uint64) error
	RemoveIdempotencies(ctx context.Context, TTL time.Duration) error
	BuildFromAlerts(alerts []alert.Alert, firingLen int, createdTime time.Time) ([]notification.Notification, error)
}

type SilenceService interface {
	Create(ctx context.Context, sil silence.Silence) (string, error)
	List(ctx context.Context, filter silence.Filter) ([]silence.Silence, error)
	Get(ctx context.Context, id string) (silence.Silence, error)
	Delete(ctx context.Context, id string) error
}

type Deps struct {
	TemplateService     TemplateService
	RuleService         RuleService
	AlertService        AlertService
	ProviderService     ProviderService
	NamespaceService    NamespaceService
	ReceiverService     ReceiverService
	SubscriptionService SubscriptionService
	NotificationService NotificationService
	SilenceService      SilenceService
}
