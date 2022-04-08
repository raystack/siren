package rules

import (
	"context"
	"fmt"

	cortexClient "github.com/grafana/cortex-tools/pkg/client"
	rwrulefmt "github.com/grafana/cortex-tools/pkg/rules/rwrulefmt"
	"github.com/odpf/siren/domain"
	"github.com/odpf/siren/pkg/templates"
	"github.com/odpf/siren/store"
	"github.com/pkg/errors"
	"github.com/prometheus/prometheus/pkg/rulefmt"
	"gopkg.in/yaml.v3"
)

const (
	namePrefix = "siren_api"
)

type variable struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Value       string `json:"value"`
	Description string `json:"description"`
}

type Variables struct {
	Variables []variable `json:"variables"`
}

var cortexClientInstance = newCortexClient

type cortexCaller interface {
	CreateRuleGroup(ctx context.Context, namespace string, rg rwrulefmt.RuleGroup) error
	DeleteRuleGroup(ctx context.Context, namespace, groupName string) error
	GetRuleGroup(ctx context.Context, namespace, groupName string) (*rwrulefmt.RuleGroup, error)
	ListRules(ctx context.Context, namespace string) (map[string][]rwrulefmt.RuleGroup, error)
}

// Service handles business logic
type Service struct {
	repository       store.RuleRepository
	templateService  domain.TemplatesService
	namespaceService domain.NamespaceService
	providerService  domain.ProviderService
}

// NewService returns repository struct
func NewService(
	repository store.RuleRepository,
	templateRepository store.TemplatesRepository,
	namespaceService domain.NamespaceService,
	providerService domain.ProviderService,
) *Service {
	return &Service{
		repository:       repository,
		templateService:  templates.NewService(templateRepository),
		namespaceService: namespaceService,
		providerService:  providerService,
	}
}

func (s *Service) Migrate() error {
	return s.repository.Migrate()
}

func (s *Service) Upsert(ctx context.Context, rule *domain.Rule) error {
	rule.Name = fmt.Sprintf("%s_%s_%s_%s", namePrefix, rule.Namespace, rule.GroupName, rule.Template)

	template, err := s.templateService.GetByName(rule.Template)
	if err != nil {
		return errors.Wrap(err, "s.templateService.GetByName")
	}
	if template == nil {
		return errors.New("template not found")
	}
	templateVariables := template.Variables
	finalRuleVariables := mergeRuleVariablesWithDefaults(templateVariables, rule.Variables)
	rule.Variables = finalRuleVariables

	namespace, err := s.namespaceService.GetNamespace(rule.ProviderNamespace)
	if err != nil {
		return errors.Wrap(err, "s.namespaceService.GetNamespace")
	}
	if namespace == nil {
		return errors.New("namespace not found")
	}
	provider, err := s.providerService.GetProvider(namespace.Provider)
	if err != nil {
		return errors.Wrap(err, "s.providerService.GetProvider")
	}
	if provider == nil {
		return errors.New("provider not found")
	}

	rule.Name = fmt.Sprintf("%s_%s_%s_%s_%s_%s", namePrefix, provider.Urn,
		namespace.Urn, rule.Namespace, rule.GroupName, rule.Template)

	ctx = s.repository.WithTransaction(ctx)
	if err := s.repository.Upsert(ctx, rule, s.templateService); err != nil {
		if err := s.repository.Rollback(ctx); err != nil {
			return errors.Wrap(err, "s.repository.Rollback")
		}
		return errors.Wrap(err, "s.repository.Upsert")
	}

	if provider.Type == "cortex" {
		client, err := cortexClientInstance(provider.Host)
		if err != nil {
			if err := s.repository.Rollback(ctx); err != nil {
				return errors.Wrap(err, "s.repository.Rollback")
			}
			return errors.Wrap(err, "cortexClientInstance")
		}

		rulesWithinGroup, err := s.repository.ListByGroup(ctx, rule.Namespace, rule.GroupName, rule.ProviderNamespace)
		if err != nil {
			if err := s.repository.Rollback(ctx); err != nil {
				return errors.Wrap(err, "s.repository.Rollback")
			}
			return errors.Wrap(err, "s.repository.ListByGroup")
		}

		if err := s.postRuleGroupWith(rule, rulesWithinGroup, client, namespace.Urn); err != nil {
			if err := s.repository.Rollback(ctx); err != nil {
				return errors.Wrap(err, "s.repository.Rollback")
			}
			return errors.Wrap(err, "s.postRuleGroupWith")
		}
	} else {
		if err := s.repository.Rollback(ctx); err != nil {
			return errors.Wrap(err, "s.repository.Rollback")
		}
		return errors.New("provider not supported")
	}

	if err := s.repository.Commit(ctx); err != nil {
		return errors.Wrap(err, "s.repository.Commit")
	}
	return nil
}

func (s *Service) Get(ctx context.Context, name, namespace, groupName, template string, providerNamespace uint64) ([]domain.Rule, error) {
	return s.repository.Get(ctx, name, namespace, groupName, template, providerNamespace)
}

func (s *Service) postRuleGroupWith(rule *domain.Rule, rulesWithinGroup []*domain.Rule, client cortexCaller, tenantName string) error {
	renderedBodyForThisGroup := ""
	for i := 0; i < len(rulesWithinGroup); i++ {
		if !rulesWithinGroup[i].Enabled {
			continue
		}
		inputValue := make(map[string]string)

		for _, v := range rulesWithinGroup[i].Variables {
			inputValue[v.Name] = v.Value
		}

		renderedBody, err := s.templateService.Render(rulesWithinGroup[i].Template, inputValue)
		if err != nil {
			return err
		}
		renderedBodyForThisGroup += renderedBody
	}
	ctx := cortexClient.NewContextWithTenantId(context.Background(), tenantName)
	if renderedBodyForThisGroup == "" {
		err := client.DeleteRuleGroup(ctx, rule.Namespace, rule.GroupName)
		if err != nil {
			if err.Error() == "requested resource not found" {
				return nil
			} else {
				return err
			}
		}
		return nil
	}
	var ruleNodes []rulefmt.RuleNode
	err := yaml.Unmarshal([]byte(renderedBodyForThisGroup), &ruleNodes)
	if err != nil {
		return err
	}
	y := rwrulefmt.RuleGroup{
		RuleGroup: rulefmt.RuleGroup{
			Name:  rule.GroupName,
			Rules: ruleNodes,
		},
	}
	return client.CreateRuleGroup(ctx, rule.Namespace, y)
}

func mergeRuleVariablesWithDefaults(templateVariables []domain.Variable, ruleVariables []domain.RuleVariable) []domain.RuleVariable {
	var finalRuleVariables []domain.RuleVariable
	for j := 0; j < len(templateVariables); j++ {
		variableExist := false
		matchingIndex := 0
		for k := 0; k < len(ruleVariables); k++ {
			if ruleVariables[k].Name == templateVariables[j].Name {
				variableExist = true
				matchingIndex = k
			}
		}
		if !variableExist {
			finalRuleVariables = append(finalRuleVariables, domain.RuleVariable{
				Name:        templateVariables[j].Name,
				Type:        templateVariables[j].Type,
				Value:       templateVariables[j].Default,
				Description: templateVariables[j].Description,
			})
		} else {
			finalRuleVariables = append(finalRuleVariables, ruleVariables[matchingIndex])
		}
	}
	return finalRuleVariables
}

func newCortexClient(host string) (cortexCaller, error) {
	cortexConfig := cortexClient.Config{
		Address:         host,
		UseLegacyRoutes: false,
	}

	client, err := cortexClient.New(cortexConfig)
	if err != nil {
		return nil, err
	}

	return client, nil
}
