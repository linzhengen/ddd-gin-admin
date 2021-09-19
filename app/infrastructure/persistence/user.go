package persistence

import (
	"context"

	"github.com/linzhengen/ddd-gin-admin/app/domain/valueobject/errors"

	"github.com/linzhengen/ddd-gin-admin/app/domain/valueobject/request"
	"github.com/linzhengen/ddd-gin-admin/app/domain/valueobject/response"

	"github.com/linzhengen/ddd-gin-admin/app/infrastructure/persistence/gormx"

	"github.com/jinzhu/gorm"
	"github.com/linzhengen/ddd-gin-admin/app/domain/entity"
	"github.com/linzhengen/ddd-gin-admin/app/domain/repository"
)

func getUserDB(ctx context.Context, defDB *gorm.DB) *gorm.DB {
	return gormx.GetDBWithModel(ctx, defDB, new(entity.User))
}

func NewUser(db *gorm.DB) repository.UserRepository {
	return &user{
		db: db,
	}
}

type user struct {
	db *gorm.DB
}

func (a *user) Query(ctx context.Context, req request.UserQueryRequest) (entity.Users, *response.Pagination, error) {
	db := getUserDB(ctx, a.db)
	if v := req.UserName; v != "" {
		db = db.Where("user_name=?", v)
	}
	if v := req.Status; v > 0 {
		db = db.Where("status=?", v)
	}
	if v := req.RoleIDs; len(v) > 0 {
		subQuery := getUserRoleDB(ctx, a.db).
			Select("user_id").
			Where("role_id IN (?)", v).
			SubQuery()
		db = db.Where("id IN ?", subQuery)
	}
	if v := req.QueryValue; v != "" {
		v = "%" + v + "%"
		db = db.Where("user_name LIKE ? OR real_name LIKE ? OR phone LIKE ? OR email LIKE ?", v, v, v, v)
	}

	db = db.Order(gormx.ParseOrder(req.OrderFields.AddIdSortKey()))

	var list entity.Users
	p, err := gormx.WrapPageQuery(ctx, db, req.PaginationParam, &list)
	if err != nil {
		return nil, nil, errors.WithStack(err)
	}
	return list, p, nil
}

func (a *user) Get(ctx context.Context, id string) (*entity.User, error) {
	var item *entity.User
	ok, err := gormx.FindOne(ctx, getUserDB(ctx, a.db).Where("id=?", id), &item)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if !ok {
		return nil, nil
	}

	return item, nil
}

func (a *user) Create(ctx context.Context, item entity.User) error {
	result := getUserDB(ctx, a.db).Create(item)
	return errors.WithStack(result.Error)
}

func (a *user) Update(ctx context.Context, id string, item entity.User) error {
	result := getUserDB(ctx, a.db).Where("id=?", id).Updates(item)
	return errors.WithStack(result.Error)
}

func (a *user) Delete(ctx context.Context, id string) error {
	result := getUserDB(ctx, a.db).Where("id=?", id).Delete(entity.User{})
	return errors.WithStack(result.Error)
}

func (a *user) UpdateStatus(ctx context.Context, id string, status int) error {
	result := getUserDB(ctx, a.db).Where("id=?", id).Update("status", status)
	return errors.WithStack(result.Error)
}

func (a *user) UpdatePassword(ctx context.Context, id, password string) error {
	result := getUserDB(ctx, a.db).Where("id=?", id).Update("password", password)
	return errors.WithStack(result.Error)
}
