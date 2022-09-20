package postgres

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/odpf/siren/core/rule"
	"github.com/odpf/siren/internal/store/model"
	"github.com/odpf/siren/pkg/errors"
)

const ruleUpsertQuery = `
INSERT INTO
rules
	(name, namespace, group_name, template, enabled, variables, provider_namespace, created_at, updated_at)
VALUES
	($1, $2, $3, $4, $5, $6, $7, now(), now())
ON CONFLICT
	(name)
DO UPDATE SET
	(namespace, group_name, template, enabled, variables, provider_namespace, updated_at) =
	($2, $3, $4, $5, $6, $7, now())
RETURNING *`

var ruleListQueryBuilder = sq.Select(
	"id",
	"name",
	"namespace",
	"group_name",
	"template",
	"enabled",
	"variables",
	"provider_namespace",
	"created_at",
	"updated_at",
).From("rules")

// RuleRepository talks to the store to read or insert data
type RuleRepository struct {
	client *Client
}

// NewRuleRepository returns repository struct
func NewRuleRepository(client *Client) *RuleRepository {
	return &RuleRepository{client}
}

func (r *RuleRepository) Upsert(ctx context.Context, rl *rule.Rule) error {
	if rl == nil {
		return errors.New("rule domain is nil")
	}

	ruleModel := new(model.Rule)
	if err := ruleModel.FromDomain(*rl); err != nil {
		return err
	}

	var newRuleModel model.Rule
	if err := r.client.GetDB(ctx).QueryRowxContext(ctx, ruleUpsertQuery,
		ruleModel.Name,
		ruleModel.Namespace,
		ruleModel.GroupName,
		ruleModel.Template,
		ruleModel.Enabled,
		ruleModel.Variables,
		ruleModel.ProviderNamespace,
	).StructScan(&newRuleModel); err != nil {
		err = checkPostgresError(err)
		if errors.Is(err, errDuplicateKey) {
			return rule.ErrDuplicate
		}
		if errors.Is(err, errForeignKeyViolation) {
			return rule.ErrRelation
		}
		return err
	}

	newRule, err := newRuleModel.ToDomain()
	if err != nil {
		return err
	}

	*rl = *newRule

	return nil
}

func (r *RuleRepository) List(ctx context.Context, flt rule.Filter) ([]rule.Rule, error) {
	var queryBuilder = ruleListQueryBuilder
	if flt.Name != "" {
		queryBuilder = queryBuilder.Where("name = ?", flt.Name)
	}
	if flt.Namespace != "" {
		queryBuilder = queryBuilder.Where("namespace = ?", flt.Namespace)
	}
	if flt.GroupName != "" {
		queryBuilder = queryBuilder.Where("group_name = ?", flt.GroupName)
	}
	if flt.TemplateName != "" {
		queryBuilder = queryBuilder.Where("template = ?", flt.TemplateName)
	}
	if flt.NamespaceID != 0 {
		queryBuilder = queryBuilder.Where("provider_namespace = ?", flt.NamespaceID)
	}

	query, args, err := queryBuilder.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.client.GetDB(ctx).QueryxContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rulesDomain []rule.Rule
	for rows.Next() {
		var ruleModel model.Rule
		if err := rows.StructScan(&ruleModel); err != nil {
			return nil, err
		}

		newRule, err := ruleModel.ToDomain()
		if err != nil {
			return nil, err
		}
		rulesDomain = append(rulesDomain, *newRule)
	}

	return rulesDomain, nil
}

func (r *RuleRepository) WithTransaction(ctx context.Context) context.Context {
	return r.client.WithTransaction(ctx, nil)
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
