package cortex

import (
	"context"
	"fmt"

	"github.com/grafana/cortex-tools/pkg/rules/rwrulefmt"
	"github.com/odpf/siren/core/subscription"
)

//go:generate mockery --name=CortexClient -r --case underscore --with-expecter --structname CortexClient --filename cortex_client.go --output=./mocks
type CortexClient interface {
	CreateAlertmanagerConfig(AlertManagerConfig, string) error
	CreateRuleGroup(context.Context, string, rwrulefmt.RuleGroup) error
	DeleteRuleGroup(context.Context, string, string) error
	GetRuleGroup(context.Context, string, string) (*rwrulefmt.RuleGroup, error)
	ListRules(context.Context, string) (map[string][]rwrulefmt.RuleGroup, error)
}

type CortexService struct {
	cortexClient CortexClient
}

// NewProviderService returns cortex service struct
func NewProviderService(cortexClient CortexClient) *CortexService {
	return &CortexService{
		cortexClient: cortexClient,
	}
}

func (s *CortexService) SyncSubscriptions(_ context.Context, subscriptions []subscription.Subscription, namespaceURN string) error {
	amConfig := make([]ReceiverConfig, 0)
	for _, item := range subscriptions {
		amConfig = append(amConfig, GetAlertManagerReceiverConfig(&item)...)
	}

	if err := s.cortexClient.CreateAlertmanagerConfig(AlertManagerConfig{
		Receivers: amConfig,
	}, namespaceURN); err != nil {
		return fmt.Errorf("error calling cortex: %w", err)
	}

	return nil
}

func GetAlertManagerReceiverConfig(subs *subscription.Subscription) []ReceiverConfig {
	if subs == nil {
		return nil
	}
	amReceiverConfig := make([]ReceiverConfig, 0)
	for idx, item := range subs.Receivers {
		configMapString := make(map[string]string)
		for key, value := range item.Configuration {
			strKey := fmt.Sprintf("%v", key)
			strValue := fmt.Sprintf("%v", value)

			configMapString[strKey] = strValue
		}
		newAMReceiver := ReceiverConfig{
			Name:           fmt.Sprintf("%s_receiverId_%d_idx_%d", subs.URN, item.ID, idx),
			Match:          subs.Match,
			Configurations: configMapString,
			Type:           item.Type,
		}
		amReceiverConfig = append(amReceiverConfig, newAMReceiver)
	}
	return amReceiverConfig
}

// func (s *CortexService) MergeReceiverConfigs(sub *subscription.Subscription, cortexReceivers []*promconfig.Receiver) []ReceiverConfig {
// 	var receiverConfigs []ReceiverConfig

// 	subscriptionReceivers := GetAlertManagerReceiverConfig(sub)

// 	for _, sam := range subscriptionReceivers {
// 		for i := 0; i < len(cortexReceivers); i++ {
// 			if sam.Name == cortexReceivers[0].Name {

// 			}
// 		}
// 	}

// 	return receiverConfigs
// }

// func (s *CortexService) UpsertSubscription(ctx context.Context, sub *subscription.Subscription, namespaceURN string) error {
// 	return plugins.ErrProviderSyncMethodNotSupported
// 	// cortexAMConfigYaml, _, err := s.cortexClient.GetAlertmanagerConfig(ctx, namespaceURN)
// 	// if err != nil {
// 	// 	return err
// 	// }

// 	// var cortexAMConfig promconfig.Config
// 	// if err := yaml.Unmarshal([]byte(cortexAMConfigYaml), &cortexAMConfig); err != nil {
// 	// 	return errors.New("cannot parse cortex alertmanager config string from remote alertmanager")
// 	// }

// 	// amReceivers := cortexAMConfig.Receivers
// 	// subReceivers := sub.Receivers

// 	// resultReceivers := make([]ReceiverConfig, 0)
// 	// for _, sr := range subReceivers {
// 	// 	for _, amr := range amReceivers {
// 	// 		if sr.ID == amr.Name {

// 	// 		}
// 	// 		switch sr.Type {
// 	// 		case receiver.TypeHTTP:
// 	// 			for _, cfg := range amr.WebhookConfigs {

// 	// 			}
// 	// 		case receiver.TypePagerDuty:
// 	// 		case receiver.TypeSlack:
// 	// 		}
// 	// 	}
// 	// }
// 	// amReceiverConfig := make([]ReceiverConfig, 0)
// 	// for idx, item := range subs.Receivers {
// 	// 	configMapString := make(map[string]string)
// 	// 	for key, value := range item.Configuration {
// 	// 		strKey := fmt.Sprintf("%v", key)
// 	// 		strValue := fmt.Sprintf("%v", value)

// 	// 		configMapString[strKey] = strValue
// 	// 	}
// 	// 	newAMReceiver := ReceiverConfig{
// 	// 		Receiver:      fmt.Sprintf("%s_receiverId_%d_idx_%d", subs.URN, item.ID, idx),
// 	// 		Match:         subs.Match,
// 	// 		Configuration: configMapString,
// 	// 		Type:          item.Type,
// 	// 	}
// 	// 	amReceiverConfig = append(amReceiverConfig, newAMReceiver)
// 	// }
// 	// return amReceiverConfig

// 	// amConfig := make([]ReceiverConfig, 0)
// 	// for _, item := range subscriptions {
// 	// 	amConfig = append(amConfig, GetAlertManagerReceiverConfig(&item)...)
// 	// }

// 	// if err := s.cortexClient.CreateAlertmanagerConfig(AlertManagerConfig{
// 	// 	Receivers: amConfig,
// 	// }, namespaceURN); err != nil {
// 	// 	return fmt.Errorf("error calling cortex: %w", err)
// 	// }

// 	return nil
// }
