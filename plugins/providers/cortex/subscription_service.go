package cortex

import (
	"context"

	"github.com/odpf/siren/core/subscription"
)

func (s *CortexService) CreateSubscription(ctx context.Context, sub *subscription.Subscription, subscriptionsInNamespace []subscription.Subscription, namespaceURN string) error {
	return s.SyncSubscriptions(ctx, subscriptionsInNamespace, namespaceURN)
}

func (s *CortexService) UpdateSubscription(ctx context.Context, sub *subscription.Subscription, subscriptionsInNamespace []subscription.Subscription, namespaceURN string) error {
	return s.SyncSubscriptions(ctx, subscriptionsInNamespace, namespaceURN)
}

func (s *CortexService) DeleteSubscription(ctx context.Context, sub *subscription.Subscription, subscriptionsInNamespace []subscription.Subscription, namespaceURN string) error {
	return s.SyncSubscriptions(ctx, subscriptionsInNamespace, namespaceURN)
}
