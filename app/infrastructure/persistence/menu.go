package persistence

import (
	"context"

	"github.com/linzhengen/ddd-gin-admin/app/domain/valueobject/errors"

	"github.com/linzhengen/ddd-gin-admin/app/domain/valueobject/response"

	"github.com/linzhengen/ddd-gin-admin/app/domain/valueobject/request"

	"github.com/linzhengen/ddd-gin-admin/app/infrastructure/persistence/gormx"

	"github.com/jinzhu/gorm"
	"github.com/linzhengen/ddd-gin-admin/app/domain/entity"
	"github.com/linzhengen/ddd-gin-admin/app/domain/repository"
)

func getMenuDB(ctx context.Context, defDB *gorm.DB) *gorm.DB {
	return gormx.GetDBWithModel(ctx, defDB, new(entity.Menu))
}

func NewMenu(db *gorm.DB) repository.MenuRepository {
	return &menu{
		db: db,
	}
}

type menu struct {
	db *gorm.DB
}

func (a *menu) Query(ctx context.Context, req request.MenuQueryRequest) (entity.Menus, *response.Pagination, error) {
	db := getMenuDB(ctx, a.db)
	if v := req.IDs; len(v) > 0 {
		db = db.Where("id IN (?)", v)
	}
	if v := req.Name; v != "" {
		db = db.Where("name=?", v)
	}
	if v := req.ParentID; v != nil {
		db = db.Where("parent_id=?", *v)
	}
	if v := req.PrefixParentPath; v != "" {
		db = db.Where("parent_path LIKE ?", v+"%")
	}
	if v := req.ShowStatus; v != 0 {
		db = db.Where("show_status=?", v)
	}
	if v := req.Status; v != 0 {
		db = db.Where("status=?", v)
	}
	if v := req.QueryValue; v != "" {
		v = "%" + v + "%"
		db = db.Where("name LIKE ? OR memo LIKE ?", v, v)
	}

	db = db.Order(gormx.ParseOrder(req.OrderFields.AddIdSortKey()))

	var list entity.Menus
	p, err := gormx.WrapPageQuery(ctx, db, req.PaginationParam, &list)
	if err != nil {
		return nil, nil, errors.WithStack(err)
	}

	return list, p, nil
}

func (a *menu) Get(ctx context.Context, id string) (*entity.Menu, error) {
	var item *entity.Menu
	ok, err := gormx.FindOne(ctx, getMenuDB(ctx, a.db).Where("id=?", id), &item)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if !ok {
		return nil, nil
	}

	return item, nil
}

func (a *menu) Create(ctx context.Context, item entity.Menu) error {
	result := getMenuDB(ctx, a.db).Create(item)
	return errors.WithStack(result.Error)
}

func (a *menu) Update(ctx context.Context, id string, item entity.Menu) error {
	result := getMenuDB(ctx, a.db).Where("id=?", id).Updates(item)
	return errors.WithStack(result.Error)
}

func (a *menu) UpdateParentPath(ctx context.Context, id, parentPath string) error {
	result := getMenuDB(ctx, a.db).Where("id=?", id).Update("parent_path", parentPath)
	return errors.WithStack(result.Error)
}

func (a *menu) Delete(ctx context.Context, id string) error {
	result := getMenuDB(ctx, a.db).Where("id=?", id).Delete(entity.Menu{})
	return errors.WithStack(result.Error)
}

func (a *menu) UpdateStatus(ctx context.Context, id string, status int) error {
	result := getMenuDB(ctx, a.db).Where("id=?", id).Update("status", status)
	return errors.WithStack(result.Error)
}
