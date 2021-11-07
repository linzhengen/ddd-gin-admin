package gormx

import (
	"context"
	"fmt"
	"strings"

	"github.com/linzhengen/ddd-gin-admin/app/domain/contextx"
	"github.com/linzhengen/ddd-gin-admin/app/domain/pagination"

	"github.com/linzhengen/ddd-gin-admin/configs"

	"github.com/jinzhu/gorm"
)

type TransFunc func(context.Context) error

func ExecTrans(ctx context.Context, db *gorm.DB, fn TransFunc) error {
	if _, ok := contextx.FromTrans(ctx); ok {
		return fn(ctx)
	}

	return db.Transaction(func(db *gorm.DB) error {
		return fn(contextx.NewTrans(ctx, db))
	})
}

func ExecTransWithLock(ctx context.Context, db *gorm.DB, fn TransFunc) error {
	if !contextx.FromTransLock(ctx) {
		ctx = contextx.NewTransLock(ctx)
	}
	return ExecTrans(ctx, db, fn)
}

func GetDB(ctx context.Context, defDB *gorm.DB) *gorm.DB {
	trans, ok := contextx.FromTrans(ctx)
	if ok && !contextx.FromNoTrans(ctx) {
		db, ok := trans.(*gorm.DB)
		if ok {
			if contextx.FromTransLock(ctx) {
				if dbType := configs.C.Gorm.DBType; dbType == "mysql" ||
					dbType == "postgres" {
					db = db.Set("gorm:query_option", "FOR UPDATE")
				}
			}
			return db
		}
	}
	return defDB
}

func GetDBWithModel(ctx context.Context, defDB *gorm.DB, m interface{}) *gorm.DB {
	return GetDB(ctx, defDB).Model(m)
}

func WrapPageQuery(ctx context.Context, db *gorm.DB, pp pagination.Param, out interface{}) (*pagination.Pagination, error) {
	if pp.OnlyCount {
		var count int
		err := db.Count(&count).Error
		if err != nil {
			return nil, err
		}
		return &pagination.Pagination{Total: count}, nil
	}
	if !pp.Pagination {
		err := db.Find(out).Error
		return nil, err
	}

	total, err := FindPage(ctx, db, pp, out)
	if err != nil {
		return nil, err
	}

	return &pagination.Pagination{
		Total:    total,
		Current:  pp.GetCurrent(),
		PageSize: pp.GetPageSize(),
	}, nil
}

func FindPage(ctx context.Context, db *gorm.DB, pp pagination.Param, out interface{}) (int, error) {
	var count int
	err := db.Count(&count).Error
	if err != nil {
		return 0, err
	}
	if count == 0 {
		return count, nil
	}

	current, pageSize := pp.GetCurrent(), pp.GetPageSize()
	if current > 0 && pageSize > 0 {
		db = db.Offset((current - 1) * pageSize).Limit(pageSize)
	}
	if pageSize > 0 {
		db = db.Limit(pageSize)
	}

	err = db.Find(out).Error
	return count, err
}

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

func Check(ctx context.Context, db *gorm.DB) (bool, error) {
	var count int
	result := db.Count(&count)
	if err := result.Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

type OrderFieldFunc func(string) string

func ParseOrder(items []*pagination.OrderField, handle ...OrderFieldFunc) string {
	orders := make([]string, len(items))

	for i, item := range items {
		key := item.Key
		if len(handle) > 0 {
			key = handle[0](key)
		}

		direction := "ASC"
		if item.Direction == pagination.OrderByDESC {
			direction = "DESC"
		}
		orders[i] = fmt.Sprintf("%s %s", key, direction)
	}

	return strings.Join(orders, ",")
}
