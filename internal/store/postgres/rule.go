package postgres

import (
	"context"
	"fmt"

	"github.com/odpf/siren/core/rule"
	"github.com/odpf/siren/internal/store/model"
	"github.com/odpf/siren/pkg/errors"
)

// RuleRepository talks to the store to read or insert data
type RuleRepository struct {
	client *Client
}

// NewRuleRepository returns repository struct
func NewRuleRepository(client *Client) *RuleRepository {
	return &RuleRepository{client}
}

func (r *RuleRepository) Upsert(ctx context.Context, rl *rule.Rule) error {
	m := new(model.Rule)
	if err := m.FromDomain(rl); err != nil {
		return err
	}

	tx := r.client.GetDB(ctx)

	result := tx.WithContext(ctx).Where("name = ?", m.Name).Updates(&m)
	if result.Error != nil {
		err := checkPostgresError(result.Error)
		if errors.Is(err, errDuplicateKey) {
			return rule.ErrDuplicate
		}
		if errors.Is(err, errForeignKeyViolation) {
			return rule.ErrRelation
		}
		return err
	}

	if result.RowsAffected == 0 {
		if err := tx.WithContext(ctx).Create(&m).Error; err != nil {
			err = checkPostgresError(err)
			if errors.Is(err, errDuplicateKey) {
				return rule.ErrDuplicate
			}
			if errors.Is(err, errForeignKeyViolation) {
				return rule.ErrRelation
			}
			return err
		}
	}

	newRule, err := m.ToDomain()
	if err != nil {
		return err
	}

	*rl = *newRule
	return nil
}

func (r *RuleRepository) List(ctx context.Context, flt rule.Filter) ([]rule.Rule, error) {
	var rules []model.Rule
	txdb := r.client.GetDB(ctx).WithContext(ctx)
	if flt.Name != "" {
		txdb = txdb.Where("name = ?", flt.Name)
	}
	if flt.Namespace != "" {
		txdb = txdb.Where("namespace = ?", flt.Namespace)
	}
	if flt.GroupName != "" {
		txdb = txdb.Where("group_name = ?", flt.GroupName)
	}
	if flt.TemplateName != "" {
		txdb = txdb.Where("template = ?", flt.TemplateName)
	}
	if flt.NamespaceID != 0 {
		txdb = txdb.Where("provider_namespace = ?", flt.NamespaceID)
	}

	if err := txdb.Find(&rules).Error; err != nil {
		return nil, err
	}

	var domainRules []rule.Rule
	for _, r := range rules {
		rule, err := r.ToDomain()
		if err != nil {
			return nil, err
		}
		domainRules = append(domainRules, *rule)
	}

	return domainRules, nil
}

func (r *RuleRepository) WithTransaction(ctx context.Context) context.Context {
	return r.client.WithTransaction(ctx)
}

func (r *RuleRepository) Rollback(ctx context.Context, err error) error {
	if txErr := r.client.Rollback(ctx); txErr != nil {
		return fmt.Errorf("rollback error %s with error: %w", txErr.Error(), err)
	}
	return nil
}

func (r *RuleRepository) Commit(ctx context.Context) error {
	return r.client.Commit(ctx)
}
