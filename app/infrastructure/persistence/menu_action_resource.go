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

func getMenuActionResourceDB(ctx context.Context, defDB *gorm.DB) *gorm.DB {
	return gormx.GetDBWithModel(ctx, defDB, new(entity.MenuActionResource))
}

func NewMenuActionResource(db *gorm.DB) repository.MenuActionResourceRepository {
	return &menuActionResource{
		db: db,
	}
}

type menuActionResource struct {
	db *gorm.DB
}

func (a *menuActionResource) Query(ctx context.Context, req request.MenuActionResourceQueryRequest) (entity.MenuActionResources, *response.Pagination, error) {
	db := getMenuActionResourceDB(ctx, a.db)
	if v := req.MenuID; v != "" {
		subQuery := getMenuActionDB(ctx, a.db).
			Where("menu_id=?", v).
			Select("id").SubQuery()
		db = db.Where("action_id IN ?", subQuery)
	}
	if v := req.MenuIDs; len(v) > 0 {
		subQuery := getMenuActionDB(ctx, a.db).Where("menu_id IN (?)", v).Select("id").SubQuery()
		db = db.Where("action_id IN ?", subQuery)
	}

	db = db.Order(gormx.ParseOrder(req.OrderFields.AddIdSortKey()))

	var list entity.MenuActionResources
	p, err := gormx.WrapPageQuery(ctx, db, req.PaginationParam, &list)
	if err != nil {
		return nil, nil, errors.WithStack(err)
	}

	return list, p, nil
}

func (a *menuActionResource) Get(ctx context.Context, id string) (*entity.MenuActionResource, error) {
	db := getMenuActionResourceDB(ctx, a.db).Where("id=?", id)
	var item *entity.MenuActionResource
	ok, err := gormx.FindOne(ctx, db, &item)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if !ok {
		return nil, nil
	}

	return item, nil
}

func (a *menuActionResource) Create(ctx context.Context, item entity.MenuActionResource) error {
	result := getMenuActionResourceDB(ctx, a.db).Create(item)
	return errors.WithStack(result.Error)
}

func (a *menuActionResource) Update(ctx context.Context, id string, item entity.MenuActionResource) error {
	result := getMenuActionResourceDB(ctx, a.db).Where("id=?", id).Updates(item)
	return errors.WithStack(result.Error)
}

func (a *menuActionResource) Delete(ctx context.Context, id string) error {
	result := getMenuActionResourceDB(ctx, a.db).Where("id=?", id).Delete(entity.MenuActionResource{})
	return errors.WithStack(result.Error)
}

func (a *menuActionResource) DeleteByActionID(ctx context.Context, actionID string) error {
	result := getMenuActionResourceDB(ctx, a.db).Where("action_id =?", actionID).Delete(entity.MenuActionResource{})
	return errors.WithStack(result.Error)
}

func (a *menuActionResource) DeleteByMenuID(ctx context.Context, menuID string) error {
	subQuery := getMenuActionDB(ctx, a.db).Where("menu_id=?", menuID).Select("id").SubQuery()
	result := getMenuActionResourceDB(ctx, a.db).Where("action_id IN ?", subQuery).Delete(entity.MenuActionResource{})
	return errors.WithStack(result.Error)
}
