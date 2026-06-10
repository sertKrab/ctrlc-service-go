package store

import (
	"context"

	"gorm.io/gorm"
)

// ExecuteOptions controls transaction behaviour for a store operation.
type ExecuteOptions struct {
	UseTransaction bool
}

// DefaultOpts runs without a transaction.
var DefaultOpts = ExecuteOptions{UseTransaction: false}

// TxOpts runs inside a transaction.
var TxOpts = ExecuteOptions{UseTransaction: true}

// Execute runs fn against the database, optionally inside a transaction.
// Service/usecase methods must call Execute and never call db.Begin, tx.Commit,
// or tx.Rollback directly — the executor manages the full connection lifecycle.
func Execute[T any](ctx context.Context, db *gorm.DB, fn func(tx *gorm.DB) (T, error), opts ExecuteOptions) (T, error) {
	if !opts.UseTransaction {
		return fn(db.WithContext(ctx))
	}
	var result T
	err := db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var e error
		result, e = fn(tx)
		return e
	})
	return result, err
}
