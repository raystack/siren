package cortex

import (
	"context"
	"fmt"

	"github.com/grafana/cortex-tools/pkg/rules/rwrulefmt"
	"github.com/odpf/siren/core/provider"
	"github.com/odpf/siren/core/rule"
	"github.com/odpf/siren/core/subscription"
	"github.com/odpf/siren/core/template"
	"github.com/odpf/siren/pkg/errors"
	"github.com/odpf/siren/plugins"
	"github.com/prometheus/prometheus/pkg/rulefmt"
	"gopkg.in/yaml.v3"
)

//go:generate mockery --name=CortexClient -r --case underscore --with-expecter --structname CortexClient --filename cortex_client.go --output=./mocks
type CortexClient interface {
	CreateAlertmanagerConfig(AlertManagerConfig, string) error
	CreateRuleGroup(context.Context, string, rwrulefmt.RuleGroup) error
	DeleteRuleGroup(context.Context, string, string) error
	GetRuleGroup(context.Context, string, string) (*rwrulefmt.RuleGroup, error)
	ListRules(context.Context, string) (map[string][]rwrulefmt.RuleGroup, error)
	// GetAlertmanagerConfig(ctx context.Context, tenantID string) (string, map[string]string, error)
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

func (s *CortexService) SyncMethod() provider.SyncMethod {
	return provider.TypeSyncBatch
}

func (s *CortexService) UpsertRule(ctx context.Context, rl *rule.Rule, templateToUpdate *template.Template, namespaceURN string) error {
	inputValues := make(map[string]string)
	for _, v := range rl.Variables {
		inputValues[v.Name] = v.Value
	}

	renderedRule, err := template.RenderWithTemplate(ctx, templateToUpdate, inputValues)
	if err != nil {
		return err
	}

	var upsertedRuleNodes []rulefmt.RuleNode
	if err := yaml.Unmarshal([]byte(renderedRule), &upsertedRuleNodes); err != nil {
		return errors.ErrInvalid.WithMsgf("cannot parse upserted rule").WithCausef(err.Error())
	}

	cortexRuleGroup, err := s.cortexClient.GetRuleGroup(ctx, rl.Namespace, rl.GroupName)
	if err != nil {
		return errors.ErrInvalid.WithMsgf("cannot get rule group from cortex when upserting rules").WithCausef(err.Error())
	}

	newRuleNodes, err := MergeRuleNodes(cortexRuleGroup.Rules, upsertedRuleNodes, rl.Enabled)
	if err != nil {
		return err
	}

	if len(newRuleNodes) == 0 {
		if err := s.cortexClient.DeleteRuleGroup(NewContext(ctx, namespaceURN), rl.Namespace, rl.GroupName); err != nil {
			if err.Error() == "requested resource not found" {
				return nil
			}
			return fmt.Errorf("error calling cortex: %w", err)
		}
	}

	cortexRuleGroup.RuleGroup.Rules = newRuleNodes
	if err := s.cortexClient.CreateRuleGroup(ctx, rl.Namespace, *cortexRuleGroup); err != nil {
		return fmt.Errorf("error calling cortex: %w", err)
	}
	return nil
}

func MergeRuleNodes(ruleNodes []rulefmt.RuleNode, newRuleNodes []rulefmt.RuleNode, enabled bool) ([]rulefmt.RuleNode, error) {
	for _, nrn := range newRuleNodes {
		var status string = "insert"
		var idxCount = 0
		for _, ruleNode := range ruleNodes {
			if ruleNode.Alert.Value == nrn.Alert.Value {
				if !enabled {
					status = "delete"
					break
				}
				status = "update"
				break
			}
			idxCount++
		}

		switch status {
		case "delete":
			if idxCount >= len(ruleNodes) || idxCount < 0 {
				return nil, errors.New("something wrong when comparing rule node")
			}
			ruleNodes = append(ruleNodes[:idxCount], ruleNodes[idxCount+1:]...)
		case "update":
			ruleNodes[idxCount] = nrn
		default:
			if !enabled {
				return ruleNodes, nil
			}
			ruleNodes = append(ruleNodes, nrn)
		}
	}

	return ruleNodes, nil
}

// func (s *CortexService) UploadRuleGroup(ctx context.Context, templateService TemplateService, nspace, groupName, tenantName string, rulesWithinGroup []rule.Rule) error {
// 	renderedBodyForThisGroup := ""

// 	for _, ruleWithinGroup := range rulesWithinGroup {
// 		if !ruleWithinGroup.Enabled {
// 			continue
// 		}
// 		inputValue := make(map[string]string)

// 		for _, v := range ruleWithinGroup.Variables {
// 			inputValue[v.Name] = v.Value
// 		}

// 		renderedBody, err := templateService.Render(ctx, ruleWithinGroup.Template, inputValue)
// 		if err != nil {
// 			return err
// 		}
// 		renderedBodyForThisGroup += renderedBody
// 	}

// 	if renderedBodyForThisGroup == "" {
// 		if err := s.cortexClient.DeleteRuleGroup(NewContext(ctx, tenantName), nspace, groupName); err != nil {
// 			if err.Error() == "requested resource not found" {
// 				return nil
// 			}
// 			return fmt.Errorf("error calling cortex: %w", err)
// 		}
// 		return nil
// 	}

// 	var ruleNodes []rulefmt.RuleNode
// 	err := yaml.Unmarshal([]byte(renderedBodyForThisGroup), &ruleNodes)
// 	if err != nil {
// 		return errors.ErrInvalid.WithMsgf("cannot parse rules to alert manage rule nodes format, check your rule or template").WithCausef(err.Error())
// 	}
// 	y := rwrulefmt.RuleGroup{
// 		RuleGroup: rulefmt.RuleGroup{
// 			Name:  groupName,
// 			Rules: ruleNodes,
// 		},
// 	}
// 	if err := s.cortexClient.CreateRuleGroup(ctx, nspace, y); err != nil {
// 		return fmt.Errorf("error calling cortex: %w", err)
// 	}
// 	return nil

// }

// subscriptions
func (s *CortexService) CreateSubscription(ctx context.Context, sub *subscription.Subscription, namespaceURN string) error {
	return plugins.ErrProviderSyncMethodNotSupported
}

func (s *CortexService) UpdateSubscription(ctx context.Context, sub *subscription.Subscription, namespaceURN string) error {
	return plugins.ErrProviderSyncMethodNotSupported
}

func (s *CortexService) DeleteSubscription(ctx context.Context, sub *subscription.Subscription, namespaceURN string) error {
	return plugins.ErrProviderSyncMethodNotSupported
}

func (s *CortexService) SyncSubscriptions(_ context.Context, subscriptions []subscription.Subscription, namespaceURN string) error {
	amConfig := make([]ReceiverConfig, 0)
	for _, item := range subscriptions {
		amConfig = append(amConfig, getAlertManagerReceiverConfig(&item)...)
	}

	if err := s.cortexClient.CreateAlertmanagerConfig(AlertManagerConfig{
		Receivers: amConfig,
	}, namespaceURN); err != nil {
		return fmt.Errorf("error calling cortex: %w", err)
	}

	return nil
}

func getAlertManagerReceiverConfig(subs *subscription.Subscription) []ReceiverConfig {
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
