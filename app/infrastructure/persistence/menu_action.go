package persistence

import (
	"context"

	"github.com/linzhengen/ddd-gin-admin/app/infrastructure/persistence/gormx"

	"github.com/linzhengen/ddd-gin-admin/app/domain/valueobject/errors"

	"github.com/jinzhu/gorm"
	"github.com/linzhengen/ddd-gin-admin/app/domain/entity"
	"github.com/linzhengen/ddd-gin-admin/app/domain/repository"
	"github.com/linzhengen/ddd-gin-admin/app/domain/valueobject/schema"
)

func getMenuActionDB(ctx context.Context, defDB *gorm.DB) *gorm.DB {
	return gormx.GetDBWithModel(ctx, defDB, new(entity.MenuAction))
}

func NewMenuAction(db *gorm.DB) repository.MenuActionRepository {
	return &menuAction{
		db: db,
	}
}

type menuAction struct {
	db *gorm.DB
}

func (a *menuAction) Query(ctx context.Context, params schema.MenuActionQueryParam) (entity.MenuActions, *schema.PaginationResult, error) {
	db := getMenuActionDB(ctx, a.db)
	if v := params.MenuID; v != "" {
		db = db.Where("menu_id=?", v)
	}
	if v := params.IDs; len(v) > 0 {
		db = db.Where("id IN (?)", v)
	}

	db = db.Order(gormx.ParseOrder(params.OrderFields.AddIdSortField()))

	var list entity.MenuActions
	pr, err := gormx.WrapPageQuery(ctx, db, params.PaginationParam, &list)
	if err != nil {
		return nil, nil, errors.WithStack(err)
	}

	return list, pr, nil
}

func (a *menuAction) Get(ctx context.Context, id string) (*entity.MenuAction, error) {
	db := getMenuActionDB(ctx, a.db).Where("id=?", id)
	item := &entity.MenuAction{}
	ok, err := gormx.FindOne(ctx, db, &item)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if !ok {
		return nil, nil
	}

	return item, nil
}

func (a *menuAction) Create(ctx context.Context, item entity.MenuAction) error {
	result := getMenuActionDB(ctx, a.db).Create(item)
	return errors.WithStack(result.Error)
}

func (a *menuAction) Update(ctx context.Context, id string, item entity.MenuAction) error {
	result := getMenuActionDB(ctx, a.db).Where("id=?", id).Updates(item)
	return errors.WithStack(result.Error)
}

func (a *menuAction) Delete(ctx context.Context, id string) error {
	result := getMenuActionDB(ctx, a.db).Where("id=?", id).Delete(entity.MenuAction{})
	return errors.WithStack(result.Error)
}

func (a *menuAction) DeleteByMenuID(ctx context.Context, menuID string) error {
	result := getMenuActionDB(ctx, a.db).Where("menu_id=?", menuID).Delete(entity.MenuAction{})
	return errors.WithStack(result.Error)
}
