package subscription

import (
	"context"
)

// SubscriptionSyncer is an interface for the provider to upload subscription(s).
// Provider plugin needs to implement this interface in order to support subscription
// synchronization from siren to provider. All methods in the interface pass
// 2 kind of data: a subscription and all subscriptions in the tenant.
// This is to accomodate some provider (e.g. cortex) that only support bulk upload.
//
//go:generate mockery --name=SubscriptionSyncer -r --case underscore --with-expecter --structname SubscriptionSyncer --filename subscription_syncer.go --output=./mocks
type SubscriptionSyncer interface {
	CreateSubscription(ctx context.Context, sub *Subscription, subscriptionsInNamespace []Subscription, namespaceURN string) error
	UpdateSubscription(ctx context.Context, sub *Subscription, subscriptionsInNamespace []Subscription, namespaceURN string) error
	DeleteSubscription(ctx context.Context, sub *Subscription, subscriptionsInNamespace []Subscription, namespaceURN string) error
}
