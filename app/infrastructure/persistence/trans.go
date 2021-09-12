package persistence

import (
	"context"

	"github.com/linzhengen/ddd-gin-admin/app/domain/contextx"

	"github.com/jinzhu/gorm"
	"github.com/linzhengen/ddd-gin-admin/app/domain/repository"
)

func NewTrans(db *gorm.DB) repository.TransRepository {
	return &trans{
		db: db,
	}
}

type trans struct {
	db *gorm.DB
}

func (a *trans) Exec(ctx context.Context, fn func(context.Context) error) error {
	if _, ok := contextx.FromTrans(ctx); ok {
		return fn(ctx)
	}

	return a.db.Transaction(func(db *gorm.DB) error {
		return fn(contextx.NewTrans(ctx, db))
	})
}
