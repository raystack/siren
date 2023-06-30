package rule

import (
	"context"
	"fmt"

	"github.com/raystack/siren/core/namespace"
	"github.com/raystack/siren/core/template"
	"github.com/raystack/siren/pkg/errors"
)

const (
	namePrefix = "siren_api"
)

//go:generate mockery --name=NamespaceService -r --case underscore --with-expecter --structname NamespaceService --filename namespace_service.go --output=./mocks
type NamespaceService interface {
	List(context.Context) ([]namespace.Namespace, error)
	Create(context.Context, *namespace.Namespace) error
	Get(context.Context, uint64) (*namespace.Namespace, error)
	Update(context.Context, *namespace.Namespace) error
	Delete(context.Context, uint64) error
}

//go:generate mockery --name=TemplateService -r --case underscore --with-expecter --structname TemplateService --filename template_service.go --output=./mocks
type TemplateService interface {
	Upsert(context.Context, *template.Template) error
	List(context.Context, template.Filter) ([]template.Template, error)
	GetByName(context.Context, string) (*template.Template, error)
	Delete(context.Context, string) error
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
	repository            Repository
	templateService       TemplateService
	namespaceService      NamespaceService
	ruleUploadersRegistry map[string]RuleUploader
}

// NewService returns repository struct
func NewService(
	repository Repository,
	templateService TemplateService,
	namespaceService NamespaceService,
	ruleUploadersRegistry map[string]RuleUploader,
) *Service {
	return &Service{
		repository:            repository,
		templateService:       templateService,
		namespaceService:      namespaceService,
		ruleUploadersRegistry: ruleUploadersRegistry,
	}
}

func (s *Service) Upsert(ctx context.Context, rl *Rule) error {
	ns, err := s.namespaceService.Get(ctx, rl.ProviderNamespace)
	if err != nil {
		return err
	}

	templateToUpdate, err := s.templateService.GetByName(ctx, rl.Template)
	if err != nil {
		return err
	}

	templateVariables := templateToUpdate.Variables
	finalRuleVariables := mergeRuleVariablesWithDefaults(templateVariables, rl.Variables)
	rl.Variables = finalRuleVariables

	rl.Name = fmt.Sprintf("%s_%s_%s_%s_%s_%s", namePrefix, ns.Provider.URN,
		ns.URN, rl.Namespace, rl.GroupName, rl.Template)

	ctx = s.repository.WithTransaction(ctx)
	if err = s.repository.Upsert(ctx, rl); err != nil {
		if err := s.repository.Rollback(ctx, err); err != nil {
			return err
		}
		return err
	}

	pluginService, err := s.getProviderPluginService(ns.Provider.Type)
	if err != nil {
		if err := s.repository.Rollback(ctx, err); err != nil {
			return err
		}
		return err
	}

	if err := pluginService.UpsertRule(ctx, ns.URN, ns.Provider, rl, templateToUpdate); err != nil {
		if err := s.repository.Rollback(ctx, err); err != nil {
			return err
		}
		return err
	}

	if err := s.repository.Commit(ctx); err != nil {
		return err
	}

	return nil
}

func (s *Service) getProviderPluginService(providerType string) (RuleUploader, error) {
	pluginService, exist := s.ruleUploadersRegistry[providerType]
	if !exist {
		return nil, errors.ErrInvalid.WithMsgf("unsupported provider type: %q", providerType)
	}
	return pluginService, nil
}

func (s *Service) List(ctx context.Context, flt Filter) ([]Rule, error) {
	return s.repository.List(ctx, flt)
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
