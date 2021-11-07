package trans

import (
	"context"

	"github.com/linzhengen/ddd-gin-admin/app/domain/trans"

	"github.com/jinzhu/gorm"
	"github.com/linzhengen/ddd-gin-admin/app/domain/contextx"
)

func NewRepository(db *gorm.DB) trans.Repository {
	return &repository{
		db: db,
	}
}

type repository struct {
	db *gorm.DB
}

func (a *repository) Exec(ctx context.Context, fn func(context.Context) error) error {
	if _, ok := contextx.FromTrans(ctx); ok {
		return fn(ctx)
	}

	return a.db.Transaction(func(db *gorm.DB) error {
		return fn(contextx.NewTrans(ctx, db))
	})
}
