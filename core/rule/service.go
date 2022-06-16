package rule

import (
	"context"
	"fmt"

	cortexClient "github.com/grafana/cortex-tools/pkg/client"
	rwrulefmt "github.com/grafana/cortex-tools/pkg/rules/rwrulefmt"
	"github.com/odpf/siren/core/namespace"
	"github.com/odpf/siren/core/provider"
	"github.com/odpf/siren/core/template"
	"github.com/pkg/errors"
	"github.com/prometheus/prometheus/pkg/rulefmt"
	"gopkg.in/yaml.v3"
)

const (
	namePrefix = "siren_api"
)

//go:generate mockery --name=NamespaceService -r --case underscore --with-expecter --structname NamespaceService --filename namespace_service.go --output=./mocks
type NamespaceService interface { //TODO to be refactored, for temporary only
	ListNamespaces() ([]*namespace.Namespace, error)
	CreateNamespace(*namespace.Namespace) error
	GetNamespace(uint64) (*namespace.Namespace, error)
	UpdateNamespace(*namespace.Namespace) error
	DeleteNamespace(uint64) error
	Migrate() error
}

//go:generate mockery --name=ProviderService -r --case underscore --with-expecter --structname ProviderService --filename provider_service.go --output=./mocks
type ProviderService interface { //TODO to be refactored, for temporary only
	ListProviders(map[string]interface{}) ([]*provider.Provider, error)
	CreateProvider(*provider.Provider) (*provider.Provider, error)
	GetProvider(uint64) (*provider.Provider, error)
	UpdateProvider(*provider.Provider) (*provider.Provider, error)
	DeleteProvider(uint64) error
	Migrate() error
}

//go:generate mockery --name=TemplatesService -r --case underscore --with-expecter --structname TemplatesService --filename template_service.go --output=./mocks
type TemplatesService interface {
	Upsert(*template.Template) error
	Index(string) ([]template.Template, error)
	GetByName(string) (*template.Template, error)
	Delete(string) error
	Render(string, map[string]string) (string, error)
	Migrate() error
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
	templateService  TemplatesService
	namespaceService NamespaceService
	providerService  ProviderService
	cortexClient     CortexClient
}

// NewService returns repository struct
func NewService(
	repository Repository,
	templateService TemplatesService,
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

func (s *Service) Migrate() error {
	return s.repository.Migrate()
}

func (s *Service) Upsert(ctx context.Context, rule *Rule) error {
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

	rule.Name = fmt.Sprintf("%s_%s_%s_%s_%s_%s", namePrefix, provider.URN,
		namespace.URN, rule.Namespace, rule.GroupName, rule.Template)

	ctx = s.repository.WithTransaction(ctx)
	if err := s.repository.Upsert(ctx, rule); err != nil {
		if err := s.repository.Rollback(ctx); err != nil {
			return errors.Wrap(err, "s.repository.Rollback")
		}
		return errors.Wrap(err, "s.repository.Upsert")
	}

	if provider.Type == "cortex" {
		rulesWithinGroup, err := s.repository.Get(ctx, "", rule.Namespace, rule.GroupName, "", rule.ProviderNamespace)
		if err != nil {
			if err := s.repository.Rollback(ctx); err != nil {
				return errors.Wrap(err, "s.repository.Rollback")
			}
			return errors.Wrap(err, "s.repository.Get")
		}

		if err := s.postRuleGroupWith(ctx, rule, rulesWithinGroup, s.cortexClient, namespace.URN); err != nil {
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

func (s *Service) Get(ctx context.Context, name, namespace, groupName, template string, providerNamespace uint64) ([]Rule, error) {
	return s.repository.Get(ctx, name, namespace, groupName, template, providerNamespace)
}

func (s *Service) postRuleGroupWith(ctx context.Context, rule *Rule, rulesWithinGroup []Rule, client CortexClient, tenantName string) error {
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
			return errors.Wrap(err, "s.templateService.Render")
		}
		renderedBodyForThisGroup += renderedBody
	}
	ctx = cortexClient.NewContextWithTenantId(ctx, tenantName)
	if renderedBodyForThisGroup == "" {
		err := client.DeleteRuleGroup(ctx, rule.Namespace, rule.GroupName)
		if err != nil {
			if err.Error() == "requested resource not found" {
				return nil
			} else {
				return errors.Wrap(err, "client.DeleteRuleGroup")
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
	if err := client.CreateRuleGroup(ctx, rule.Namespace, y); err != nil {
		return errors.Wrap(err, "client.CreateRuleGroup")
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
