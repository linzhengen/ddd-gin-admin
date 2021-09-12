package persistence

import (
	"context"

	"github.com/linzhengen/ddd-gin-admin/app/infrastructure/persistence/gormx"

	"github.com/linzhengen/ddd-gin-admin/app/domain/errors"

	"github.com/jinzhu/gorm"
	"github.com/linzhengen/ddd-gin-admin/app/domain/entity"
	"github.com/linzhengen/ddd-gin-admin/app/domain/repository"
	"github.com/linzhengen/ddd-gin-admin/app/domain/schema"
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

func (a *menuAction) getQueryOption(opts ...schema.MenuActionQueryOptions) schema.MenuActionQueryOptions {
	var opt schema.MenuActionQueryOptions
	if len(opts) > 0 {
		opt = opts[0]
	}
	return opt
}

func (a *menuAction) Query(ctx context.Context, params schema.MenuActionQueryParam, opts ...schema.MenuActionQueryOptions) (*schema.MenuActionQueryResult, error) {
	opt := a.getQueryOption(opts...)

	db := getMenuActionDB(ctx, a.db)
	if v := params.MenuID; v != "" {
		db = db.Where("menu_id=?", v)
	}
	if v := params.IDs; len(v) > 0 {
		db = db.Where("id IN (?)", v)
	}

	opt.OrderFields = append(opt.OrderFields, schema.NewOrderField("id", schema.OrderByASC))
	db = db.Order(gormx.ParseOrder(opt.OrderFields))

	var list entity.MenuActions
	pr, err := gormx.WrapPageQuery(ctx, db, params.PaginationParam, &list)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	qr := &schema.MenuActionQueryResult{
		PageResult: pr,
		Data:       list.ToSchemaMenuActions(),
	}

	return qr, nil
}

func (a *menuAction) Get(ctx context.Context, id string, opts ...schema.MenuActionQueryOptions) (*schema.MenuAction, error) {
	db := getMenuActionDB(ctx, a.db).Where("id=?", id)
	var item entity.MenuAction
	ok, err := gormx.FindOne(ctx, db, &item)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if !ok {
		return nil, nil
	}

	return item.ToSchemaMenuAction(), nil
}

func (a *menuAction) Create(ctx context.Context, item schema.MenuAction) error {
	eitem := entity.SchemaMenuAction(item).ToMenuAction()
	result := getMenuActionDB(ctx, a.db).Create(eitem)
	return errors.WithStack(result.Error)
}

func (a *menuAction) Update(ctx context.Context, id string, item schema.MenuAction) error {
	eitem := entity.SchemaMenuAction(item).ToMenuAction()
	result := getMenuActionDB(ctx, a.db).Where("id=?", id).Updates(eitem)
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
