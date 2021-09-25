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

func (a *userRole) Query(ctx context.Context, params schema.UserRoleQueryParam) (entity.UserRoles, *schema.PaginationResult, error) {
	db := getUserRoleDB(ctx, a.db)
	if v := params.UserID; v != "" {
		db = db.Where("user_id=?", v)
	}
	if v := params.UserIDs; len(v) > 0 {
		db = db.Where("user_id IN (?)", v)
	}

	db = db.Order(gormx.ParseOrder(params.OrderFields.AddIdSortField()))

	var list entity.UserRoles
	pr, err := gormx.WrapPageQuery(ctx, db, params.PaginationParam, &list)
	if err != nil {
		return nil, nil, errors.WithStack(err)
	}
	return list, pr, nil
}

func (a *userRole) Get(ctx context.Context, id string) (*entity.UserRole, error) {
	db := getUserRoleDB(ctx, a.db).Where("id=?", id)
	item := &entity.UserRole{}
	ok, err := gormx.FindOne(ctx, db, &item)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if !ok {
		return nil, nil
	}

	return item, nil
}

func (a *userRole) Create(ctx context.Context, item entity.UserRole) error {
	result := getUserRoleDB(ctx, a.db).Create(item)
	return errors.WithStack(result.Error)
}

func (a *userRole) Update(ctx context.Context, id string, item entity.UserRole) error {
	result := getUserRoleDB(ctx, a.db).Where("id=?", id).Updates(item)
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
