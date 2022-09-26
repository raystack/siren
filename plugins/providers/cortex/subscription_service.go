package cortex

import (
	"context"

	"github.com/odpf/siren/core/subscription"
)

// CreateSubscription is the abstraction to create a subscription called in core service
func (s *CortexService) CreateSubscription(ctx context.Context, sub *subscription.Subscription, subscriptionsInNamespace []subscription.Subscription, namespaceURN string) error {
	return s.SyncSubscriptions(ctx, subscriptionsInNamespace, namespaceURN)
}

// UpdateSubscription is the abstraction to update a subscription called in core service
func (s *CortexService) UpdateSubscription(ctx context.Context, sub *subscription.Subscription, subscriptionsInNamespace []subscription.Subscription, namespaceURN string) error {
	return s.SyncSubscriptions(ctx, subscriptionsInNamespace, namespaceURN)
}

// DeleteSubscription is the abstraction to remove a subscription called in core service
func (s *CortexService) DeleteSubscription(ctx context.Context, sub *subscription.Subscription, subscriptionsInNamespace []subscription.Subscription, namespaceURN string) error {
	return s.SyncSubscriptions(ctx, subscriptionsInNamespace, namespaceURN)
}
