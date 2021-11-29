package subscription

import (
	"fmt"
	"github.com/odpf/siren/domain"
	"github.com/odpf/siren/pkg/subscription/alertmanager"
	"github.com/pkg/errors"
)

func (r Repository) addReceiversConfiguration(subscriptions []Subscription, receiverService domain.ReceiverService) ([]SubscriptionEnrichedWithReceivers, error) {
	res := make([]SubscriptionEnrichedWithReceivers, 0)
	allReceivers, err := receiverService.ListReceivers()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get receivers")
	}
	for _, item := range subscriptions {
		enrichedReceivers := make([]EnrichedReceiverMetadata, 0)
		for _, receiverItem := range item.Receiver {
			var receiverInfo *domain.Receiver
			found := false
			for idx := range allReceivers {
				if allReceivers[idx].Id == receiverItem.Id {
					found = true
					receiverInfo = allReceivers[idx]
					break
				}
			}
			if found != true {
				return nil, errors.New(fmt.Sprintf("receiver id %d does not exist", receiverItem.Id))
			}
			//initialize the nil map using the make function
			//to avoid panics while adding elements in future
			if receiverItem.Configuration == nil {
				receiverItem.Configuration = make(map[string]string)
			}
			if receiverInfo.Type == "slack" {
				if _, ok := receiverItem.Configuration["channel_name"]; !ok {
					return nil, errors.New(fmt.Sprintf(
						"configuration.channel_name missing from receiver with id %d", receiverItem.Id))
				}
				if val, ok := receiverInfo.Configurations["token"]; ok {
					receiverItem.Configuration["token"] = val.(string)
				}
			} else if receiverInfo.Type == "pagerduty" {
				if val, ok := receiverInfo.Configurations["service_key"]; ok {
					receiverItem.Configuration["service_key"] = val.(string)
				}
			} else if receiverInfo.Type == "http" {
				if val, ok := receiverInfo.Configurations["url"]; ok {
					receiverItem.Configuration["url"] = val.(string)
				}
			} else {
				return nil, errors.New(fmt.Sprintf(`subscriptions for receiver type %s not supported via Siren inside Cortex`, receiverInfo.Type))
			}
			enrichedReceiver := EnrichedReceiverMetadata{
				Id:            receiverItem.Id,
				Configuration: receiverItem.Configuration,
				Type:          receiverInfo.Type,
			}
			enrichedReceivers = append(enrichedReceivers, enrichedReceiver)
		}
		enrichedSubscription := SubscriptionEnrichedWithReceivers{
			Id:          item.Id,
			NamespaceId: item.NamespaceId,
			Urn:         item.Urn,
			Receiver:    enrichedReceivers,
			Match:       item.Match,
		}
		res = append(res, enrichedSubscription)
	}
	return res, nil
}

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
