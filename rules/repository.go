package rules

import (
	"context"
	"encoding/json"
	"fmt"
	cortexClient "github.com/grafana/cortex-tools/pkg/client"
	"github.com/grafana/cortex-tools/pkg/rules/rwrulefmt"
	"github.com/odpf/siren/domain"
	"github.com/odpf/siren/templates"
	"github.com/prometheus/prometheus/pkg/rulefmt"
	"gopkg.in/yaml.v3"
	"gorm.io/gorm"
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

// Repository talks to the store to read or insert data
type Repository struct {
	db *gorm.DB
}

// NewRepository returns repository struct
func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r Repository) Migrate() error {
	err := r.db.AutoMigrate(&Rule{})
	if err != nil {
		return err
	}
	return nil
}

func postRuleGroupWith(rule *Rule, rulesWithinGroup []Rule, client cortexCaller, db *gorm.DB) error {
	renderedBodyForThisGroup := ""
	for i := 0; i < len(rulesWithinGroup); i++ {
		if rulesWithinGroup[i].Status == "disabled" {
			continue
		}
		inputValue := make(map[string]string)
		var variables []variable
		jsonBlob := []byte(rulesWithinGroup[i].Variables)
		_ = json.Unmarshal(jsonBlob, &variables)
		for _, v := range variables {
			inputValue[v.Name] = v.Value
		}
		service := templates.NewService(db)
		renderedBody, err := service.Render(rulesWithinGroup[i].Template, inputValue)
		if err != nil {
			return nil
		}
		renderedBodyForThisGroup += renderedBody
	}
	ctx := cortexClient.NewContextWithTenantId(context.Background(), rule.Entity)
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
		fmt.Println(err)
	}
	y := rwrulefmt.RuleGroup{
		RuleGroup: rulefmt.RuleGroup{
			Name:  rule.GroupName,
			Rules: ruleNodes,
		},
	}
	err = client.CreateRuleGroup(ctx, rule.Namespace, y)
	if err != nil {
		return err
	}
	return nil
}

func (r Repository) Upsert(rule *Rule, client cortexCaller) (*Rule, error) {
	rule.Name = fmt.Sprintf("%s_%s_%s_%s_%s", namePrefix,
		rule.Entity, rule.Namespace, rule.GroupName, rule.Template)
	var existingRule Rule
	var rulesWithinGroup []Rule

	service := templates.NewService(r.db)
	template, err := service.GetByName(rule.Template)
	templateVariables := template.Variables

	var ruleVariables []domain.RuleVariable
	jsonBlob := []byte(rule.Variables)
	err = json.Unmarshal(jsonBlob, &ruleVariables)
	if err != nil {
		return nil, err
	}
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
	jsonBytes, err := json.Marshal(finalRuleVariables)
	rule.Variables = string(jsonBytes)
	if err != nil {
		return nil, err
	}
	err = r.db.Transaction(func(tx *gorm.DB) error {
		result := tx.Where(fmt.Sprintf("name = '%s'", rule.Name)).Find(&existingRule)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			result = tx.Create(rule)
		} else {
			result = tx.Where("id = ?", existingRule.ID).Updates(rule)
		}
		if result.Error != nil {
			return result.Error
		}
		result = tx.Where(fmt.Sprintf("name = '%s'", rule.Name)).Find(&existingRule)
		result = tx.Where(fmt.Sprintf("namespace = '%s' AND entity = '%s' AND group_name = '%s'",
			rule.Namespace, rule.Entity, rule.GroupName)).Find(&rulesWithinGroup)
		if result.Error != nil {
			return result.Error
		}
		err = postRuleGroupWith(rule, rulesWithinGroup, client, r.db)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &existingRule, err
}

func (r Repository) Get(namespace string, entity string, groupName string, status string, template string) ([]Rule, error) {
	var rules []Rule
	selectQuery := `SELECT * from rules`
	selectQueryWithWhereClause := `SELECT * from rules WHERE `
	var filterConditions []string
	if namespace != "" {
		filterConditions = append(filterConditions, fmt.Sprintf("namespace = '%s' ", namespace))
	}
	if entity != "" {
		filterConditions = append(filterConditions, fmt.Sprintf("entity = '%s' ", entity))
	}
	if groupName != "" {
		filterConditions = append(filterConditions, fmt.Sprintf("group_name = '%s' ", groupName))
	}
	if status != "" {
		filterConditions = append(filterConditions, fmt.Sprintf("status = '%s' ", status))
	}
	if template != "" {
		filterConditions = append(filterConditions, fmt.Sprintf("template = '%s' ", template))
	}
	var finalSelectQuery string
	if len(filterConditions) == 0 {
		finalSelectQuery = selectQuery
	} else {
		finalSelectQuery = selectQueryWithWhereClause
		for i := 0; i < len(filterConditions); i++ {
			if i == 0 {
				finalSelectQuery += filterConditions[i]
			} else {
				finalSelectQuery += " AND " + filterConditions[i]
			}
		}
	}
	result := r.db.Raw(finalSelectQuery).Scan(&rules)
	if result.Error != nil {
		return nil, result.Error
	}
	return rules, nil
}
