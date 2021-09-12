package persistence

import (
	"context"

	"github.com/linzhengen/ddd-gin-admin/app/domain/errors"

	"github.com/linzhengen/ddd-gin-admin/app/domain/entity"
	"github.com/linzhengen/ddd-gin-admin/app/domain/repository"
	"github.com/linzhengen/ddd-gin-admin/app/domain/schema"
	"github.com/linzhengen/ddd-gin-admin/app/infrastructure/gormx"

	"github.com/jinzhu/gorm"
)

func getRoleDB(ctx context.Context, defDB *gorm.DB) *gorm.DB {
	return gormx.GetDBWithModel(ctx, defDB, new(entity.Role))
}

func NewRole(db *gorm.DB) repository.RoleRepository {
	return &role{
		db: db,
	}
}

type role struct {
	db *gorm.DB
}

func (a *role) getQueryOption(opts ...schema.RoleQueryOptions) schema.RoleQueryOptions {
	var opt schema.RoleQueryOptions
	if len(opts) > 0 {
		opt = opts[0]
	}
	return opt
}

func (a *role) Query(ctx context.Context, params schema.RoleQueryParam, opts ...schema.RoleQueryOptions) (*schema.RoleQueryResult, error) {
	opt := a.getQueryOption(opts...)

	db := getRoleDB(ctx, a.db)
	if v := params.IDs; len(v) > 0 {
		db = db.Where("id IN (?)", v)
	}
	if v := params.Name; v != "" {
		db = db.Where("name=?", v)
	}
	if v := params.UserID; v != "" {
		subQuery := getUserRoleDB(ctx, a.db).
			Where("deleted_at is null").
			Where("user_id=?", v).
			Select("role_id").SubQuery()
		db = db.Where("id IN ?", subQuery)
	}
	if v := params.QueryValue; v != "" {
		v = "%" + v + "%"
		db = db.Where("name LIKE ? OR memo LIKE ?", v, v)
	}

	opt.OrderFields = append(opt.OrderFields, schema.NewOrderField("id", schema.OrderByDESC))
	db = db.Order(gormx.ParseOrder(opt.OrderFields))

	var list entity.Roles
	pr, err := gormx.WrapPageQuery(ctx, db, params.PaginationParam, &list)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	qr := &schema.RoleQueryResult{
		PageResult: pr,
		Data:       list.ToSchemaRoles(),
	}

	return qr, nil
}

func (a *role) Get(ctx context.Context, id string, opts ...schema.RoleQueryOptions) (*schema.Role, error) {
	var role entity.Role
	ok, err := gormx.FindOne(ctx, getRoleDB(ctx, a.db).Where("id=?", id), &role)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if !ok {
		return nil, nil
	}

	return role.ToSchemaRole(), nil
}

func (a *role) Create(ctx context.Context, item schema.Role) error {
	eitem := entity.SchemaRole(item).ToRole()
	result := getRoleDB(ctx, a.db).Create(eitem)
	return errors.WithStack(result.Error)
}

func (a *role) Update(ctx context.Context, id string, item schema.Role) error {
	eitem := entity.SchemaRole(item).ToRole()
	result := getRoleDB(ctx, a.db).Where("id=?", id).Updates(eitem)
	return errors.WithStack(result.Error)
}

func (a *role) Delete(ctx context.Context, id string) error {
	result := getRoleDB(ctx, a.db).Where("id=?", id).Delete(entity.Role{})
	return errors.WithStack(result.Error)
}

func (a *role) UpdateStatus(ctx context.Context, id string, status int) error {
	result := getRoleDB(ctx, a.db).Where("id=?", id).Update("status", status)
	return errors.WithStack(result.Error)
}
