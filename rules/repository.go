package rules

import (
	"context"
	"encoding/json"
	"fmt"
	cortexClient "github.com/grafana/cortex-tools/pkg/client"
	"github.com/grafana/cortex-tools/pkg/rules/rwrulefmt"
	"github.com/odpf/siren/templates"
	"github.com/prometheus/prometheus/pkg/rulefmt"
	"gopkg.in/yaml.v3"
	"gorm.io/gorm"
)

const (
	namePrefix = "siren_api"
)

type Variable struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Value       string `json:"value"`
	Description string `json:"description"`
}

type Variables struct {
	Variables []Variable `json:"variables"`
}

// Repository talks to the store to read or insert data
type Repository struct {
	db     *gorm.DB
	client *cortexClient.CortexClient
}

// NewRepository returns repository struct
func NewRepository(db *gorm.DB) *Repository {
	cfg := cortexClient.Config{
		Address:         "http://localhost:8080",
		UseLegacyRoutes: true,
	}
	client, err := cortexClient.New(cfg)
	if err != nil {
		return nil
	}
	return &Repository{db: db, client: client}
}

func (r Repository) Migrate() error {
	err := r.db.AutoMigrate(&Rule{})
	if err != nil {
		return err
	}
	return nil
}

func (r Repository) Upsert(rule *Rule) (*Rule, error) {
	rule.Name = fmt.Sprintf("%s_%s_%s_%s_%s", namePrefix,
		rule.Entity, rule.Namespace, rule.GroupName, rule.Template)
	var existingRule Rule
	var rulesWithinGroup []Rule
	err := r.db.Transaction(func(tx *gorm.DB) error {
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

		renderedBodyForThisGroup := ""
		for i := 0; i < len(rulesWithinGroup); i++ {
			inputValue := make(map[string]string)
			var variables []Variable
			jsonBlob := []byte(rulesWithinGroup[i].Variables)
			_ = json.Unmarshal(jsonBlob, &variables)
			for _, v := range variables {
				inputValue[v.Name] = v.Value
			}
			service := templates.NewService(r.db)
			renderedBody, err := service.Render(rulesWithinGroup[i].Template, inputValue)
			if err != nil {
				return nil
			}
			renderedBodyForThisGroup += renderedBody
		}
		if result.Error != nil {
			return result.Error
		}
		var ruleNodes []rulefmt.RuleNode
		err := yaml.Unmarshal([]byte(renderedBodyForThisGroup), &ruleNodes)
		if err != nil {
			fmt.Println(err)
		}
		ctx := cortexClient.NewContextWithTenantId(context.Background(), rule.Entity)
		y := rwrulefmt.RuleGroup{
			RuleGroup: rulefmt.RuleGroup{
				Name:  rule.GroupName,
				Rules: ruleNodes,
			},
		}
		err = r.client.CreateRuleGroup(ctx, rule.Namespace, y)
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

func (r Repository) Get(string) ([]Rule, error) {
	return nil, nil
}
