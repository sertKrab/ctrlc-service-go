package store

import (
	"context"

	"gorm.io/gorm"
)

type UseCase interface {
	UseTransaction() bool
}

func Execute(ctx context.Context, db *gorm.DB, fn func(tx *gorm.DB) error, useTransaction bool) error {
	if !useTransaction {
		return fn(db)
	}
	return db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return fn(tx)
	})
}
