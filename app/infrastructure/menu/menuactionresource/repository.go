package menuactionresource

import (
	"context"

	"github.com/linzhengen/ddd-gin-admin/app/domain/menu/menuactionresource"

	"github.com/linzhengen/ddd-gin-admin/app/infrastructure/menu/menuaction"

	"github.com/linzhengen/ddd-gin-admin/app/domain/errors"

	"github.com/linzhengen/ddd-gin-admin/app/domain/pagination"

	"github.com/jinzhu/gorm"
	"github.com/linzhengen/ddd-gin-admin/app/infrastructure/gormx"
)

func GetModelDB(ctx context.Context, defDB *gorm.DB) *gorm.DB {
	return gormx.GetDBWithModel(ctx, defDB, new(Model))
}

func NewRepository(db *gorm.DB) menuactionresource.Repository {
	return &repository{
		db: db,
	}
}

type repository struct {
	db *gorm.DB
}

func (a *repository) Query(ctx context.Context, params menuactionresource.QueryParam) (menuactionresource.MenuActionResources, *pagination.Pagination, error) {
	db := GetModelDB(ctx, a.db)
	if v := params.MenuID; v != "" {
		subQuery := menuaction.GetModelDB(ctx, a.db).
			Where("menu_id=?", v).
			Select("id").SubQuery()
		db = db.Where("action_id IN ?", subQuery)
	}
	if v := params.MenuIDs; len(v) > 0 {
		subQuery := menuaction.GetModelDB(ctx, a.db).Where("menu_id IN (?)", v).Select("id").SubQuery()
		db = db.Where("action_id IN ?", subQuery)
	}

	db = db.Order(gormx.ParseOrder(params.OrderFields.AddIdSortField()))

	var list []*Model
	pr, err := gormx.WrapPageQuery(ctx, db, params.PaginationParam, &list)
	if err != nil {
		return nil, nil, errors.WithStack(err)
	}

	return toDomainList(list), pr, nil
}

func (a *repository) Get(ctx context.Context, id string) (*menuactionresource.MenuActionResource, error) {
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

func (a *repository) Create(ctx context.Context, item *menuactionresource.MenuActionResource) error {
	result := GetModelDB(ctx, a.db).Create(domainToModel(item))
	return errors.WithStack(result.Error)
}

func (a *repository) Update(ctx context.Context, id string, item *menuactionresource.MenuActionResource) error {
	result := GetModelDB(ctx, a.db).Where("id=?", id).Updates(domainToModel(item))
	return errors.WithStack(result.Error)
}

func (a *repository) Delete(ctx context.Context, id string) error {
	result := GetModelDB(ctx, a.db).Where("id=?", id).Delete(Model{})
	return errors.WithStack(result.Error)
}

func (a *repository) DeleteByActionID(ctx context.Context, actionID string) error {
	result := GetModelDB(ctx, a.db).Where("action_id =?", actionID).Delete(Model{})
	return errors.WithStack(result.Error)
}

func (a *repository) DeleteByMenuID(ctx context.Context, menuID string) error {
	subQuery := menuaction.GetModelDB(ctx, a.db).Where("menu_id=?", menuID).Select("id").SubQuery()
	result := GetModelDB(ctx, a.db).Where("action_id IN ?", subQuery).Delete(Model{})
	return errors.WithStack(result.Error)
}
