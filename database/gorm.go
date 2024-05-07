package database

import (
	"context"

	"gorm.io/gorm"
)

var _dbKey = struct{}{}

type Repository struct {
	db *gorm.DB
}

func (repo Repository) getDB(ctx context.Context) *gorm.DB {
	db, ok := ctx.Value(_dbKey).(*gorm.DB)
	if ok && db != nil {
		return db
	}

	return repo.db
}

func (repo Repository) Tx(ctx context.Context, fn func(context.Context) error) error {
	return repo.getDB(ctx).Transaction(func(tx *gorm.DB) error {
		return fn(context.WithValue(ctx, _dbKey, tx))
	})
}
