package postgres

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/raystack/siren/core/rule"
	"github.com/raystack/siren/internal/store/model"
	"github.com/raystack/siren/pkg/errors"
	"github.com/raystack/siren/pkg/pgc"
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
	client    *pgc.Client
	tableName string
}

// NewRuleRepository returns repository struct
func NewRuleRepository(client *pgc.Client) *RuleRepository {
	return &RuleRepository{client, "rules"}
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
	if err := r.client.QueryRowxContext(ctx, "UPSERT", r.tableName, ruleUpsertQuery,
		ruleModel.Name,
		ruleModel.Namespace,
		ruleModel.GroupName,
		ruleModel.Template,
		ruleModel.Enabled,
		ruleModel.Variables,
		ruleModel.ProviderNamespace,
	).StructScan(&newRuleModel); err != nil {
		err = pgc.CheckError(err)
		if errors.Is(err, pgc.ErrDuplicateKey) {
			return rule.ErrDuplicate
		}
		if errors.Is(err, pgc.ErrForeignKeyViolation) {
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

	rows, err := r.client.QueryxContext(ctx, pgc.OpSelectAll, r.tableName, query, args...)
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
