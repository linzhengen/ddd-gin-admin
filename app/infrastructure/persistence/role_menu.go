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

func getRoleMenuDB(ctx context.Context, defDB *gorm.DB) *gorm.DB {
	return gormx.GetDBWithModel(ctx, defDB, new(entity.RoleMenu))
}

func NewRoleMenu(db *gorm.DB) repository.RoleMenuRepository {
	return &roleMenu{
		db: db,
	}
}

type roleMenu struct {
	db *gorm.DB
}

func (a *roleMenu) getQueryOption(opts ...schema.RoleMenuQueryOptions) schema.RoleMenuQueryOptions {
	var opt schema.RoleMenuQueryOptions
	if len(opts) > 0 {
		opt = opts[0]
	}
	return opt
}

func (a *roleMenu) Query(ctx context.Context, params schema.RoleMenuQueryParam, opts ...schema.RoleMenuQueryOptions) (*schema.RoleMenuQueryResult, error) {
	opt := a.getQueryOption(opts...)

	db := getRoleMenuDB(ctx, a.db)
	if v := params.RoleID; v != "" {
		db = db.Where("role_id=?", v)
	}
	if v := params.RoleIDs; len(v) > 0 {
		db = db.Where("role_id IN (?)", v)
	}

	opt.OrderFields = append(opt.OrderFields, schema.NewOrderField("id", schema.OrderByDESC))
	db = db.Order(gormx.ParseOrder(opt.OrderFields))

	var list entity.RoleMenus
	pr, err := gormx.WrapPageQuery(ctx, db, params.PaginationParam, &list)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	qr := &schema.RoleMenuQueryResult{
		PageResult: pr,
		Data:       list.ToSchemaRoleMenus(),
	}

	return qr, nil
}

func (a *roleMenu) Get(ctx context.Context, id string, opts ...schema.RoleMenuQueryOptions) (*schema.RoleMenu, error) {
	db := getRoleMenuDB(ctx, a.db).Where("id=?", id)
	var item entity.RoleMenu
	ok, err := gormx.FindOne(ctx, db, &item)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if !ok {
		return nil, nil
	}

	return item.ToSchemaRoleMenu(), nil
}

func (a *roleMenu) Create(ctx context.Context, item schema.RoleMenu) error {
	eitem := entity.SchemaRoleMenu(item).ToRoleMenu()
	result := getRoleMenuDB(ctx, a.db).Create(eitem)
	return errors.WithStack(result.Error)
}

func (a *roleMenu) Update(ctx context.Context, id string, item schema.RoleMenu) error {
	eitem := entity.SchemaRoleMenu(item).ToRoleMenu()
	result := getRoleMenuDB(ctx, a.db).Where("id=?", id).Updates(eitem)
	return errors.WithStack(result.Error)
}

func (a *roleMenu) Delete(ctx context.Context, id string) error {
	result := getRoleMenuDB(ctx, a.db).Where("id=?", id).Delete(entity.RoleMenu{})
	return errors.WithStack(result.Error)
}

func (a *roleMenu) DeleteByRoleID(ctx context.Context, roleID string) error {
	result := getRoleMenuDB(ctx, a.db).Where("role_id=?", roleID).Delete(entity.RoleMenu{})
	return errors.WithStack(result.Error)
}
