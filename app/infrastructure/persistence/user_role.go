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

func getUserRoleDB(ctx context.Context, defDB *gorm.DB) *gorm.DB {
	return gormx.GetDBWithModel(ctx, defDB, new(entity.UserRole))
}

func NewUserRole(db *gorm.DB) repository.UserRoleRepository {
	return &userRole{
		db: db,
	}
}

type userRole struct {
	db *gorm.DB
}

func (a *userRole) getQueryOption(opts ...schema.UserRoleQueryOptions) schema.UserRoleQueryOptions {
	var opt schema.UserRoleQueryOptions
	if len(opts) > 0 {
		opt = opts[0]
	}
	return opt
}

func (a *userRole) Query(ctx context.Context, params schema.UserRoleQueryParam, opts ...schema.UserRoleQueryOptions) (*schema.UserRoleQueryResult, error) {
	opt := a.getQueryOption(opts...)

	db := getUserRoleDB(ctx, a.db)
	if v := params.UserID; v != "" {
		db = db.Where("user_id=?", v)
	}
	if v := params.UserIDs; len(v) > 0 {
		db = db.Where("user_id IN (?)", v)
	}

	opt.OrderFields = append(opt.OrderFields, schema.NewOrderField("id", schema.OrderByDESC))
	db = db.Order(gormx.ParseOrder(opt.OrderFields))

	var list entity.UserRoles
	pr, err := gormx.WrapPageQuery(ctx, db, params.PaginationParam, &list)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	qr := &schema.UserRoleQueryResult{
		PageResult: pr,
		Data:       list.ToSchemaUserRoles(),
	}

	return qr, nil
}

func (a *userRole) Get(ctx context.Context, id string, opts ...schema.UserRoleQueryOptions) (*schema.UserRole, error) {
	db := getUserRoleDB(ctx, a.db).Where("id=?", id)
	var item entity.UserRole
	ok, err := gormx.FindOne(ctx, db, &item)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if !ok {
		return nil, nil
	}

	return item.ToSchemaUserRole(), nil
}

func (a *userRole) Create(ctx context.Context, item schema.UserRole) error {
	eitem := entity.SchemaUserRole(item).ToUserRole()
	result := getUserRoleDB(ctx, a.db).Create(eitem)
	return errors.WithStack(result.Error)
}

func (a *userRole) Update(ctx context.Context, id string, item schema.UserRole) error {
	eitem := entity.SchemaUserRole(item).ToUserRole()
	result := getUserRoleDB(ctx, a.db).Where("id=?", id).Updates(eitem)
	return errors.WithStack(result.Error)
}

func (a *userRole) Delete(ctx context.Context, id string) error {
	result := getUserRoleDB(ctx, a.db).Where("id=?", id).Delete(entity.UserRole{})
	return errors.WithStack(result.Error)
}

func (a *userRole) DeleteByUserID(ctx context.Context, userID string) error {
	result := getUserRoleDB(ctx, a.db).Where("user_id=?", userID).Delete(entity.UserRole{})
	return errors.WithStack(result.Error)
}
