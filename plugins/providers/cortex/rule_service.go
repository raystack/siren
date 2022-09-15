package cortex

import (
	"context"
	"fmt"

	"github.com/odpf/siren/core/rule"
	"github.com/odpf/siren/core/template"
	"github.com/odpf/siren/pkg/errors"
	"github.com/prometheus/prometheus/pkg/rulefmt"
	"gopkg.in/yaml.v3"
)

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

	newRuleNodes, err := mergeRuleNodes(cortexRuleGroup.Rules, upsertedRuleNodes, rl.Enabled)
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

func mergeRuleNodes(ruleNodes []rulefmt.RuleNode, newRuleNodes []rulefmt.RuleNode, enabled bool) ([]rulefmt.RuleNode, error) {
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
