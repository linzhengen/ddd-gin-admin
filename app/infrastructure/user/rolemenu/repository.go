package rolemenu

import (
	"context"

	"github.com/linzhengen/ddd-gin-admin/app/domain/pagination"

	"github.com/linzhengen/ddd-gin-admin/app/domain/errors"
	"github.com/linzhengen/ddd-gin-admin/app/domain/user/rolemenu"

	"github.com/jinzhu/gorm"
	"github.com/linzhengen/ddd-gin-admin/app/infrastructure/gormx"
)

func NewRepository(db *gorm.DB) rolemenu.Repository {
	return &repository{
		db: db,
	}
}

type repository struct {
	db *gorm.DB
}

func GetModelDB(ctx context.Context, defDB *gorm.DB) *gorm.DB {
	return gormx.GetDBWithModel(ctx, defDB, new(Model))
}

func (a *repository) Query(ctx context.Context, params rolemenu.QueryParam) (rolemenu.RoleMenus, *pagination.Pagination, error) {
	db := GetModelDB(ctx, a.db)
	if v := params.RoleID; v != "" {
		db = db.Where("role_id=?", v)
	}
	if v := params.RoleIDs; len(v) > 0 {
		db = db.Where("role_id IN (?)", v)
	}

	db = db.Order(gormx.ParseOrder(params.OrderFields.AddIdSortField()))

	var list []*Model
	pr, err := gormx.WrapPageQuery(ctx, db, params.PaginationParam, &list)
	if err != nil {
		return nil, nil, errors.WithStack(err)
	}
	return toDomainList(list), pr, nil
}

func (a *repository) Get(ctx context.Context, id string) (*rolemenu.RoleMenu, error) {
	db := GetModelDB(ctx, a.db).Where("id=?", id)
	item := &Model{}
	ok, err := gormx.FindOne(ctx, db, &item)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if !ok {
		return nil, nil
	}

	return item.ToDomain(), nil
}

func (a *repository) Create(ctx context.Context, item *rolemenu.RoleMenu) error {
	result := GetModelDB(ctx, a.db).Create(domainToModel(item))
	return errors.WithStack(result.Error)
}

func (a *repository) Update(ctx context.Context, id string, item *rolemenu.RoleMenu) error {
	result := GetModelDB(ctx, a.db).Where("id=?", id).Updates(domainToModel(item))
	return errors.WithStack(result.Error)
}

func (a *repository) Delete(ctx context.Context, id string) error {
	result := GetModelDB(ctx, a.db).Where("id=?", id).Delete(Model{})
	return errors.WithStack(result.Error)
}

func (a *repository) DeleteByRoleID(ctx context.Context, roleID string) error {
	result := GetModelDB(ctx, a.db).Where("role_id=?", roleID).Delete(Model{})
	return errors.WithStack(result.Error)
}
