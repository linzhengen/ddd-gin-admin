package persistence

import (
	"context"

	"github.com/linzhengen/ddd-gin-admin/domain/repository"

	"github.com/google/wire"
	"github.com/jinzhu/gorm"
	"github.com/linzhengen/ddd-gin-admin/domain/entity"
	"github.com/linzhengen/ddd-gin-admin/domain/schema"
	"github.com/linzhengen/ddd-gin-admin/pkg/errors"
)

var MenuActionSet = wire.NewSet(wire.Struct(new(MenuAction), "*"))

type MenuAction struct {
	DB *gorm.DB
}

var _ repository.MenuActionRepository = &MenuAction{}

func (a *MenuAction) getQueryOption(opts ...schema.MenuActionQueryOptions) schema.MenuActionQueryOptions {
	var opt schema.MenuActionQueryOptions
	if len(opts) > 0 {
		opt = opts[0]
	}
	return opt
}

func (a *MenuAction) Query(ctx context.Context, params schema.MenuActionQueryParam, opts ...schema.MenuActionQueryOptions) (*schema.MenuActionQueryResult, error) {
	opt := a.getQueryOption(opts...)

	db := entity.GetMenuActionDB(ctx, a.DB)
	if v := params.MenuID; v != "" {
		db = db.Where("menu_id=?", v)
	}
	if v := params.IDs; len(v) > 0 {
		db = db.Where("id IN (?)", v)
	}

	opt.OrderFields = append(opt.OrderFields, schema.NewOrderField("id", schema.OrderByASC))
	db = db.Order(ParseOrder(opt.OrderFields))

	var list entity.MenuActions
	pr, err := WrapPageQuery(ctx, db, params.PaginationParam, &list)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	qr := &schema.MenuActionQueryResult{
		PageResult: pr,
		Data:       list.ToSchemaMenuActions(),
	}

	return qr, nil
}

func (a *MenuAction) Get(ctx context.Context, id string, opts ...schema.MenuActionQueryOptions) (*schema.MenuAction, error) {
	db := entity.GetMenuActionDB(ctx, a.DB).Where("id=?", id)
	var item entity.MenuAction
	ok, err := FindOne(ctx, db, &item)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if !ok {
		return nil, nil
	}

	return item.ToSchemaMenuAction(), nil
}

func (a *MenuAction) Create(ctx context.Context, item schema.MenuAction) error {
	eitem := entity.SchemaMenuAction(item).ToMenuAction()
	result := entity.GetMenuActionDB(ctx, a.DB).Create(eitem)
	return errors.WithStack(result.Error)
}

func (a *MenuAction) Update(ctx context.Context, id string, item schema.MenuAction) error {
	eitem := entity.SchemaMenuAction(item).ToMenuAction()
	result := entity.GetMenuActionDB(ctx, a.DB).Where("id=?", id).Updates(eitem)
	return errors.WithStack(result.Error)
}

func (a *MenuAction) Delete(ctx context.Context, id string) error {
	result := entity.GetMenuActionDB(ctx, a.DB).Where("id=?", id).Delete(entity.MenuAction{})
	return errors.WithStack(result.Error)
}

func (a *MenuAction) DeleteByMenuID(ctx context.Context, menuID string) error {
	result := entity.GetMenuActionDB(ctx, a.DB).Where("menu_id=?", menuID).Delete(entity.MenuAction{})
	return errors.WithStack(result.Error)
}
