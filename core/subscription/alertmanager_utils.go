package subscription

import (
	"fmt"

	"github.com/odpf/siren/pkg/cortex"
)

func getAMReceiverConfigPerSubscription(subscription SubscriptionEnrichedWithReceivers) []cortex.ReceiverConfig {
	amReceiverConfig := make([]cortex.ReceiverConfig, 0)
	for idx, item := range subscription.Receiver {
		newAMReceiver := cortex.ReceiverConfig{
			Receiver:      fmt.Sprintf("%s_receiverId_%d_idx_%d", subscription.Urn, item.Id, idx),
			Match:         subscription.Match,
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
