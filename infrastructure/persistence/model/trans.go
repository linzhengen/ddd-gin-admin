package model

import (
	"context"

	"github.com/jinzhu/gorm"
	"github.com/linzhengen/ddd-gin-admin/infrastructure/errors"
)

// NewTrans constructor.
func NewTrans(db *gorm.DB) *Trans {
	return &Trans{db}
}

// Trans db transaction.
type Trans struct {
	db *gorm.DB
}

// Begin start db transaction.
func (a *Trans) Begin(ctx context.Context) (interface{}, error) {
	result := a.db.Begin()
	if err := result.Error; err != nil {
		return nil, errors.WithStack(err)
	}
	return result, nil
}

// Commit commit transaction.
func (a *Trans) Commit(ctx context.Context, trans interface{}) error {
	db, ok := trans.(*gorm.DB)
	if !ok {
		return errors.New("unknow trans")
	}

	result := db.Commit()
	if err := result.Error; err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// Rollback roll back.
func (a *Trans) Rollback(ctx context.Context, trans interface{}) error {
	db, ok := trans.(*gorm.DB)
	if !ok {
		return errors.New("unknow trans")
	}

	result := db.Rollback()
	if err := result.Error; err != nil {
		return errors.WithStack(err)
	}
	return nil
}
