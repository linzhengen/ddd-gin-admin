package menu

import (
	"context"

	"github.com/linzhengen/ddd-gin-admin/app/domain/errors"

	"github.com/linzhengen/ddd-gin-admin/app/domain/pagination"

	"github.com/linzhengen/ddd-gin-admin/app/domain/menu"

	"github.com/jinzhu/gorm"
	"github.com/linzhengen/ddd-gin-admin/app/infrastructure/gormx"
)

func GetModelDB(ctx context.Context, defDB *gorm.DB) *gorm.DB {
	return gormx.GetDBWithModel(ctx, defDB, new(Model))
}

func NewRepository(db *gorm.DB) menu.Repository {
	return &repository{
		db: db,
	}
}

type repository struct {
	db *gorm.DB
}

func (a *repository) Query(ctx context.Context, params menu.QueryParam) (menu.Menus, *pagination.Pagination, error) {
	db := GetModelDB(ctx, a.db)
	if v := params.IDs; len(v) > 0 {
		db = db.Where("id IN (?)", v)
	}
	if v := params.Name; v != "" {
		db = db.Where("name=?", v)
	}
	if v := params.ParentID; v != nil {
		db = db.Where("parent_id=?", *v)
	}
	if v := params.PrefixParentPath; v != "" {
		db = db.Where("parent_path LIKE ?", v+"%")
	}
	if v := params.ShowStatus; v != 0 {
		db = db.Where("show_status=?", v)
	}
	if v := params.Status; v != 0 {
		db = db.Where("status=?", v)
	}
	if v := params.QueryValue; v != "" {
		v = "%" + v + "%"
		db = db.Where("name LIKE ? OR memo LIKE ?", v, v)
	}

	db = db.Order(gormx.ParseOrder(params.OrderFields.AddIdSortField()))

	var list []*Model
	pr, err := gormx.WrapPageQuery(ctx, db, params.PaginationParam, &list)
	if err != nil {
		return nil, nil, errors.WithStack(err)
	}

	return toDomainList(list), pr, nil
}

func (a *repository) Get(ctx context.Context, id string) (*menu.Menu, error) {
	item := &Model{}
	ok, err := gormx.FindOne(ctx, GetModelDB(ctx, a.db).Where("id=?", id), &item)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if !ok {
		return nil, nil
	}

	return item.ToDomain(), nil
}

func (a *repository) Create(ctx context.Context, item *menu.Menu) error {
	result := GetModelDB(ctx, a.db).Create(domainToModel(item))
	return errors.WithStack(result.Error)
}

func (a *repository) Update(ctx context.Context, id string, item *menu.Menu) error {
	result := GetModelDB(ctx, a.db).Where("id=?", id).Updates(domainToModel(item))
	return errors.WithStack(result.Error)
}

func (a *repository) UpdateParentPath(ctx context.Context, id, parentPath string) error {
	result := GetModelDB(ctx, a.db).Where("id=?", id).Update("parent_path", parentPath)
	return errors.WithStack(result.Error)
}

func (a *repository) Delete(ctx context.Context, id string) error {
	result := GetModelDB(ctx, a.db).Where("id=?", id).Delete(Model{})
	return errors.WithStack(result.Error)
}

func (a *repository) UpdateStatus(ctx context.Context, id string, status int) error {
	result := GetModelDB(ctx, a.db).Where("id=?", id).Update("status", status)
	return errors.WithStack(result.Error)
}
