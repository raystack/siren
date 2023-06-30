package postgresq

import (
	"context"
	"errors"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/raystack/siren/pkg/pgc"
	"github.com/raystack/siren/plugins/queues"
)

const defaultPublishedTimeThreshold = time.Duration(7) * time.Hour

func (q *Queue) Cleanup(ctx context.Context, filter queues.FilterCleanup) error {

	// validate filter
	var (
		publishedTimeThreshold int
		pendingTimeThreshold   int
	)
	if filter.MessagePendingTimeThreshold == "" {
		publishedTimeThreshold = int(defaultPublishedTimeThreshold.Seconds())
	} else {
		dur, err := time.ParseDuration(filter.MessagePublishedTimeThreshold)
		if err != nil {
			return err
		}
		publishedTimeThreshold = int(dur.Seconds())
	}

	if filter.MessagePendingTimeThreshold != "" {
		dur, err := time.ParseDuration(filter.MessagePendingTimeThreshold)
		if err != nil {
			return err
		}
		pendingTimeThreshold = int(dur.Seconds())
	}

	var filterExpr sq.Sqlizer
	messagePublishedExpr := sq.And{
		sq.Expr("status = 'published'"),
		sq.Expr(fmt.Sprintf("now() - interval '%d seconds' > updated_at", publishedTimeThreshold)),
	}

	if pendingTimeThreshold != 0 {
		messagePendingExpr := sq.And{
			sq.Expr("status = 'pending'"),
			sq.Expr(fmt.Sprintf("now() - interval '%d seconds' > updated_at", pendingTimeThreshold)),
		}
		filterExpr = sq.Or{
			messagePublishedExpr,
			messagePendingExpr,
		}
	} else {
		filterExpr = messagePublishedExpr
	}

	query, args, err := sq.Delete(MessageQueueTableFullName).Where(filterExpr).ToSql()
	if err != nil {
		return err
	}

	res, err := q.pgClient.ExecContext(ctx, pgc.OpDelete, MessageQueueTableFullName, query, args...)
	if err != nil {
		return err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("no rows affected when cleanup messages")
	}
	return nil
}
