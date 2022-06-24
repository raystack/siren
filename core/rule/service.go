package rule

import (
	"context"
	"fmt"

	rwrulefmt "github.com/grafana/cortex-tools/pkg/rules/rwrulefmt"
	"github.com/odpf/siren/core/namespace"
	"github.com/odpf/siren/core/provider"
	"github.com/odpf/siren/core/template"
	"github.com/odpf/siren/pkg/cortex"
	"github.com/prometheus/prometheus/pkg/rulefmt"
	"gopkg.in/yaml.v3"
)

const (
	namePrefix = "siren_api"
)

//go:generate mockery --name=NamespaceService -r --case underscore --with-expecter --structname NamespaceService --filename namespace_service.go --output=./mocks
type NamespaceService interface {
	List(context.Context) ([]namespace.Namespace, error)
	Create(context.Context, *namespace.Namespace) (uint64, error)
	Get(context.Context, uint64) (*namespace.Namespace, error)
	Update(context.Context, *namespace.Namespace) (uint64, error)
	Delete(context.Context, uint64) error
}

//go:generate mockery --name=ProviderService -r --case underscore --with-expecter --structname ProviderService --filename provider_service.go --output=./mocks
type ProviderService interface {
	List(context.Context, provider.Filter) ([]provider.Provider, error)
	Create(context.Context, *provider.Provider) (uint64, error)
	Get(context.Context, uint64) (*provider.Provider, error)
	Update(context.Context, *provider.Provider) (uint64, error)
	Delete(context.Context, uint64) error
}

//go:generate mockery --name=TemplateService -r --case underscore --with-expecter --structname TemplateService --filename template_service.go --output=./mocks
type TemplateService interface {
	Upsert(context.Context, *template.Template) (uint64, error)
	List(context.Context, template.Filter) ([]template.Template, error)
	GetByName(context.Context, string) (*template.Template, error)
	Delete(context.Context, string) error
	Render(context.Context, string, map[string]string) (string, error)
}

type variable struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Value       string `json:"value"`
	Description string `json:"description"`
}

type Variables struct {
	Variables []variable `json:"variables"`
}

// Service handles business logic
type Service struct {
	repository       Repository
	templateService  TemplateService
	namespaceService NamespaceService
	providerService  ProviderService
	cortexClient     CortexClient
}

// NewService returns repository struct
func NewService(
	repository Repository,
	templateService TemplateService,
	namespaceService NamespaceService,
	providerService ProviderService,
	cortexClient CortexClient,
) *Service {
	return &Service{
		repository:       repository,
		templateService:  templateService,
		namespaceService: namespaceService,
		providerService:  providerService,
		cortexClient:     cortexClient,
	}
}

func (s *Service) Upsert(ctx context.Context, rl *Rule) (uint64, error) {
	rl.Name = fmt.Sprintf("%s_%s_%s_%s", namePrefix, rl.Namespace, rl.GroupName, rl.Template)

	tmpl, err := s.templateService.GetByName(ctx, rl.Template)
	if err != nil {
		return 0, err
	}

	templateVariables := tmpl.Variables
	finalRuleVariables := mergeRuleVariablesWithDefaults(templateVariables, rl.Variables)
	rl.Variables = finalRuleVariables

	ns, err := s.namespaceService.Get(ctx, rl.ProviderNamespace)
	if err != nil {
		return 0, err
	}

	prov, err := s.providerService.Get(ctx, ns.Provider)
	if err != nil {
		return 0, err
	}

	rl.Name = fmt.Sprintf("%s_%s_%s_%s_%s_%s", namePrefix, prov.URN,
		ns.URN, rl.Namespace, rl.GroupName, rl.Template)

	return s.repository.UpsertWithTx(ctx, rl, func(rulesWithinGroup []Rule) error {
		if prov.Type == "cortex" {
			if err := PostRuleGroupWithCortex(ctx, s.cortexClient, s.templateService, rl, rulesWithinGroup, ns.URN); err != nil {
				return err
			}
		}
		return nil
	})
}

func (s *Service) List(ctx context.Context, flt Filter) ([]Rule, error) {
	return s.repository.List(ctx, flt)
}

func PostRuleGroupWithCortex(ctx context.Context, client CortexClient, templateService TemplateService, rl *Rule, rulesWithinGroup []Rule, tenantName string) error {
	renderedBodyForThisGroup := ""
	for _, ruleWithinGroup := range rulesWithinGroup {
		if !ruleWithinGroup.Enabled {
			continue
		}
		inputValue := make(map[string]string)

		for _, v := range ruleWithinGroup.Variables {
			inputValue[v.Name] = v.Value
		}

		renderedBody, err := templateService.Render(ctx, ruleWithinGroup.Template, inputValue)
		if err != nil {
			return err
		}
		renderedBodyForThisGroup += renderedBody
	}

	if renderedBodyForThisGroup == "" {
		err := client.DeleteRuleGroup(cortex.NewContext(ctx, tenantName), rl.Namespace, rl.GroupName)
		if err != nil {
			if err.Error() == "requested resource not found" {
				return nil
			}
			return err
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
			Name:  rl.GroupName,
			Rules: ruleNodes,
		},
	}
	if err := client.CreateRuleGroup(ctx, rl.Namespace, y); err != nil {
		return err
	}
	return nil
}

func mergeRuleVariablesWithDefaults(templateVariables []template.Variable, ruleVariables []RuleVariable) []RuleVariable {
	var finalRuleVariables []RuleVariable
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
			finalRuleVariables = append(finalRuleVariables, RuleVariable{
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
