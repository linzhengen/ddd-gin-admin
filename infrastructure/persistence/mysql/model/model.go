package model

import (
	"context"

	"github.com/jinzhu/gorm"
	"github.com/linzhengen/ddd-gin-admin/domain/schema"
	icontext "github.com/linzhengen/ddd-gin-admin/infrastructure/context"
)

// ExecTrans execute db trans.
func ExecTrans(ctx context.Context, db *gorm.DB, fn func(context.Context) error) error {
	if _, ok := icontext.FromTrans(ctx); ok {
		return fn(ctx)
	}

	transModel := NewTrans(db)
	trans, err := transModel.Begin(ctx)
	if err != nil {
		return err
	}

	ctx = icontext.NewTrans(ctx, trans)
	err = fn(ctx)
	if err != nil {
		_ = transModel.Rollback(ctx, trans)
		return err
	}
	return transModel.Commit(ctx, trans)
}

// WrapPageQuery wrap page query.
func WrapPageQuery(ctx context.Context, db *gorm.DB, pp *schema.PaginationParam, out interface{}) (*schema.PaginationResult, error) {
	if pp != nil {
		total, err := FindPage(ctx, db, pp.PageIndex, pp.PageSize, out)
		if err != nil {
			return nil, err
		}
		return &schema.PaginationResult{
			Total: total,
		}, nil
	}

	result := db.Find(out)
	return nil, result.Error
}

// FindPage find data with page index.
func FindPage(ctx context.Context, db *gorm.DB, pageIndex, pageSize int, out interface{}) (int, error) {
	var count int
	result := db.Count(&count)
	if err := result.Error; err != nil {
		return 0, err
	} else if count == 0 {
		return 0, nil
	}

	// 如果分页大小小于0或者分页索引小于0，则不查询数据
	if pageSize < 0 || pageIndex < 0 {
		return count, nil
	}

	if pageIndex > 0 && pageSize > 0 {
		db = db.Offset((pageIndex - 1) * pageSize)
	}
	if pageSize > 0 {
		db = db.Limit(pageSize)
	}
	result = db.Find(out)
	if err := result.Error; err != nil {
		return 0, err
	}

	return count, nil
}

// FindOne find one recode.
func FindOne(ctx context.Context, db *gorm.DB, out interface{}) (bool, error) {
	result := db.First(out)
	if err := result.Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// Check check data exists.
func Check(ctx context.Context, db *gorm.DB) (bool, error) {
	var count int
	result := db.Count(&count)
	if err := result.Error; err != nil {
		return false, err
	}
	return count > 0, nil
}
