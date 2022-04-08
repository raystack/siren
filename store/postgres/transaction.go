package postgres

import (
	"context"
	"errors"

	"gorm.io/gorm"
)

var (
	transactionContextKey = struct{}{}
)

type transaction struct {
	db *gorm.DB
}

func (t *transaction) WithTransaction(ctx context.Context) context.Context {
	tx := t.db.Begin()
	return context.WithValue(ctx, transactionContextKey, tx)
}

func (t *transaction) Rollback(ctx context.Context) error {
	if tx := extractTransaction(ctx); tx != nil {
		tx = tx.Rollback()
		if tx.Error != nil {
			return t.db.Error
		}
		return nil
	}
	return errors.New("no transaction")
}

func (t *transaction) Commit(ctx context.Context) error {
	if tx := extractTransaction(ctx); tx != nil {
		tx = tx.Commit()
		if tx.Error != nil {
			return t.db.Error
		}
		return nil
	}
	return errors.New("no transaction")
}

func (t *transaction) getDb(ctx context.Context) *gorm.DB {
	db := t.db
	if tx := extractTransaction(ctx); tx != nil {
		db = tx
	}
	return db
}

func extractTransaction(ctx context.Context) *gorm.DB {
	if tx, ok := ctx.Value(transactionContextKey).(*gorm.DB); !ok {
		return nil
	} else {
		return tx
	}
}
