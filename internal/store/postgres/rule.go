package postgres

import (
	"context"

	"github.com/odpf/siren/core/rule"
	"github.com/odpf/siren/internal/store/model"
	"github.com/odpf/siren/pkg/errors"
	"gorm.io/gorm"
)

// RuleRepository talks to the store to read or insert data
type RuleRepository struct {
	client *Client
}

// NewRuleRepository returns repository struct
func NewRuleRepository(client *Client) *RuleRepository {
	return &RuleRepository{client}
}

func (r *RuleRepository) UpsertWithTx(ctx context.Context, rl *rule.Rule, postProcessFn func([]rule.Rule) error) error {
	m := new(model.Rule)
	if err := m.FromDomain(rl); err != nil {
		return err
	}

	if txErr := r.client.db.Transaction(func(tx *gorm.DB) error {
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

		rulesWithinGroup, err := r.list(ctx, tx, rule.Filter{
			Namespace:   rl.Namespace,
			GroupName:   rl.GroupName,
			NamespaceID: rl.ProviderNamespace,
		})
		if err != nil {
			return err
		}

		return postProcessFn(rulesWithinGroup)
	}); txErr != nil {
		return txErr
	}

	return nil
}

func (r *RuleRepository) List(ctx context.Context, flt rule.Filter) ([]rule.Rule, error) {
	return r.list(ctx, r.client.db, flt)
}

func (r *RuleRepository) list(ctx context.Context, tx *gorm.DB, flt rule.Filter) ([]rule.Rule, error) {
	var rules []model.Rule
	txdb := tx.WithContext(ctx)
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
