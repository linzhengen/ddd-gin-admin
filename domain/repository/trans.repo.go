package repository

import (
	"context"

	"github.com/google/wire"
	"github.com/jinzhu/gorm"
	"github.com/linzhengen/ddd-gin-admin/infrastructure/contextx"
)

var TransSet = wire.NewSet(wire.Struct(new(Trans), "*"))

type Trans struct {
	DB *gorm.DB
}

func (a *Trans) Exec(ctx context.Context, fn func(context.Context) error) error {
	if _, ok := contextx.FromTrans(ctx); ok {
		return fn(ctx)
	}

	return a.DB.Transaction(func(db *gorm.DB) error {
		return fn(contextx.NewTrans(ctx, db))
	})
}
