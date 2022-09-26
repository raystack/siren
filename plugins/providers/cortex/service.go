package cortex

import (
	"context"
	"fmt"

	"github.com/grafana/cortex-tools/pkg/rules/rwrulefmt"
	"github.com/odpf/siren/core/subscription"
)

//go:generate mockery --name=CortexClient -r --case underscore --with-expecter --structname CortexClient --filename cortex_client.go --output=./mocks
type CortexClient interface {
	CreateAlertmanagerConfig(context.Context, AlertManagerConfig, string) error
	CreateRuleGroup(context.Context, string, string, rwrulefmt.RuleGroup) error
	DeleteRuleGroup(context.Context, string, string, string) error
	GetRuleGroup(context.Context, string, string, string) (*rwrulefmt.RuleGroup, error)
}

// CortexService is a service layer of cortex provider plugin
type CortexService struct {
	cortexClient CortexClient
}

// NewProviderService returns cortex service provider plugin struct
func NewProviderService(cortexClient CortexClient) *CortexService {
	return &CortexService{
		cortexClient: cortexClient,
	}
}

// SyncSubscriptions dumps all subscriptions of a tenant as an alertmanager config to cortex
// namespaceURN is the tenant ID
func (s *CortexService) SyncSubscriptions(ctx context.Context, subscriptions []subscription.Subscription, namespaceURN string) error {
	amConfig := make([]ReceiverConfig, 0)
	for _, item := range subscriptions {
		amConfig = append(amConfig, GetAlertManagerReceiverConfig(&item)...)
	}

	if err := s.cortexClient.CreateAlertmanagerConfig(ctx, AlertManagerConfig{
		Receivers: amConfig,
	}, namespaceURN); err != nil {
		return err
	}

	return nil
}

// GetAlertManagerReceiverConfig transforms subscription to list of receiver in cortex receiver config format
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
