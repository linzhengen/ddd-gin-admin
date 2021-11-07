package menuaction

import (
	"context"

	"github.com/linzhengen/ddd-gin-admin/app/domain/menu/menuaction"

	"github.com/linzhengen/ddd-gin-admin/app/domain/errors"

	"github.com/linzhengen/ddd-gin-admin/app/domain/pagination"

	"github.com/jinzhu/gorm"
	"github.com/linzhengen/ddd-gin-admin/app/infrastructure/gormx"
)

func GetModelDB(ctx context.Context, defDB *gorm.DB) *gorm.DB {
	return gormx.GetDBWithModel(ctx, defDB, new(Model))
}

func NewRepository(db *gorm.DB) menuaction.Repository {
	return &repository{
		db: db,
	}
}

type repository struct {
	db *gorm.DB
}

func (a *repository) Query(ctx context.Context, params menuaction.QueryParam) (menuaction.MenuActions, *pagination.Pagination, error) {
	db := GetModelDB(ctx, a.db)
	if v := params.MenuID; v != "" {
		db = db.Where("menu_id=?", v)
	}
	if v := params.IDs; len(v) > 0 {
		db = db.Where("id IN (?)", v)
	}

	db = db.Order(gormx.ParseOrder(params.OrderFields.AddIdSortField()))

	var list []*Model
	pr, err := gormx.WrapPageQuery(ctx, db, params.PaginationParam, &list)
	if err != nil {
		return nil, nil, errors.WithStack(err)
	}

	return toDomainList(list), pr, nil
}

func (a *repository) Get(ctx context.Context, id string) (*menuaction.MenuAction, error) {
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

func (a *repository) Create(ctx context.Context, item *menuaction.MenuAction) error {
	result := GetModelDB(ctx, a.db).Create(domainToModel(item))
	return errors.WithStack(result.Error)
}

func (a *repository) Update(ctx context.Context, id string, item *menuaction.MenuAction) error {
	result := GetModelDB(ctx, a.db).Where("id=?", id).Updates(domainToModel(item))
	return errors.WithStack(result.Error)
}

func (a *repository) Delete(ctx context.Context, id string) error {
	result := GetModelDB(ctx, a.db).Where("id=?", id).Delete(Model{})
	return errors.WithStack(result.Error)
}

func (a *repository) DeleteByMenuID(ctx context.Context, menuID string) error {
	result := GetModelDB(ctx, a.db).Where("menu_id=?", menuID).Delete(Model{})
	return errors.WithStack(result.Error)
}
