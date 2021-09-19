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

func (a *userRole) Query(ctx context.Context, req request.UserRoleQueryRequest) (entity.UserRoles, *response.Pagination, error) {
	db := getUserRoleDB(ctx, a.db)
	if v := req.UserID; v != "" {
		db = db.Where("user_id=?", v)
	}
	if v := req.UserIDs; len(v) > 0 {
		db = db.Where("user_id IN (?)", v)
	}

	db = db.Order(gormx.ParseOrder(req.OrderFields.AddIdSortKey()))

	var list entity.UserRoles
	p, err := gormx.WrapPageQuery(ctx, db, req.PaginationParam, &list)
	if err != nil {
		return nil, nil, errors.WithStack(err)
	}
	return list, p, nil
}

func (a *userRole) Get(ctx context.Context, id string) (*entity.UserRole, error) {
	db := getUserRoleDB(ctx, a.db).Where("id=?", id)
	var item *entity.UserRole
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
