package subscription

import (
	"fmt"

	"github.com/odpf/siren/core/subscription/alertmanager"
)

func getAMReceiverConfigPerSubscription(subscription SubscriptionEnrichedWithReceivers) []alertmanager.AMReceiverConfig {
	amReceiverConfig := make([]alertmanager.AMReceiverConfig, 0)
	for idx, item := range subscription.Receiver {
		newAMReceiver := alertmanager.AMReceiverConfig{
			Receiver:      fmt.Sprintf("%s_receiverId_%d_idx_%d", subscription.Urn, item.Id, idx),
			Match:         subscription.Match,
			Configuration: item.Configuration,
			Type:          item.Type,
		}
		amReceiverConfig = append(amReceiverConfig, newAMReceiver)
	}
	return amReceiverConfig
}

func getAmConfigFromSubscriptions(subscriptions []SubscriptionEnrichedWithReceivers) alertmanager.AMConfig {
	amConfig := make([]alertmanager.AMReceiverConfig, 0)
	for _, item := range subscriptions {
		amConfig = append(amConfig, getAMReceiverConfigPerSubscription(item)...)
	}
	return alertmanager.AMConfig{
		Receivers: amConfig,
	}
}
