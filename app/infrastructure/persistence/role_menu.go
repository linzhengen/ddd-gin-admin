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

func (a *roleMenu) Query(ctx context.Context, req request.RoleMenuQueryRequest) (entity.RoleMenus, *response.Pagination, error) {
	db := getRoleMenuDB(ctx, a.db)
	if v := req.RoleID; v != "" {
		db = db.Where("role_id=?", v)
	}
	if v := req.RoleIDs; len(v) > 0 {
		db = db.Where("role_id IN (?)", v)
	}

	db = db.Order(gormx.ParseOrder(req.OrderFields.AddIdSortKey()))

	var list entity.RoleMenus
	p, err := gormx.WrapPageQuery(ctx, db, req.PaginationParam, &list)
	if err != nil {
		return nil, nil, errors.WithStack(err)
	}

	return list, p, nil
}

func (a *roleMenu) Get(ctx context.Context, id string) (*entity.RoleMenu, error) {
	db := getRoleMenuDB(ctx, a.db).Where("id=?", id)
	var item *entity.RoleMenu
	ok, err := gormx.FindOne(ctx, db, &item)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if !ok {
		return nil, nil
	}

	return item, nil
}

func (a *roleMenu) Create(ctx context.Context, item entity.RoleMenu) error {
	result := getRoleMenuDB(ctx, a.db).Create(item)
	return errors.WithStack(result.Error)
}

func (a *roleMenu) Update(ctx context.Context, id string, item entity.RoleMenu) error {
	result := getRoleMenuDB(ctx, a.db).Where("id=?", id).Updates(item)
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
