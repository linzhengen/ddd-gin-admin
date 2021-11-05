package user

import (
	"context"

	"github.com/linzhengen/ddd-gin-admin/app/infrastructure/user/userrole"

	"github.com/linzhengen/ddd-gin-admin/app/domain/errors"

	"github.com/jinzhu/gorm"
	"github.com/linzhengen/ddd-gin-admin/app/domain/pagination"
	"github.com/linzhengen/ddd-gin-admin/app/domain/user"
	"github.com/linzhengen/ddd-gin-admin/app/infrastructure/gormx"
)

func NewRepository(db *gorm.DB) user.Repository {
	return &repository{
		db: db,
	}
}

type repository struct {
	db *gorm.DB
}

func GetModelDB(ctx context.Context, defDB *gorm.DB) *gorm.DB {
	return gormx.GetDBWithModel(ctx, defDB, new(Model))
}

func (a *repository) Query(ctx context.Context, params user.QueryParams) (user.Users, *pagination.Pagination, error) {
	db := GetModelDB(ctx, a.db)
	if v := params.UserName; v != "" {
		db = db.Where("user_name=?", v)
	}
	if v := params.Status; v > 0 {
		db = db.Where("status=?", v)
	}
	if v := params.RoleIDs; len(v) > 0 {
		// todo: serviceへ移動
		subQuery := userrole.GetModelDB(ctx, a.db).
			Select("user_id").
			Where("role_id IN (?)", v).
			SubQuery()
		db = db.Where("id IN ?", subQuery)
	}
	if v := params.QueryValue; v != "" {
		v = "%" + v + "%"
		db = db.Where("user_name LIKE ? OR real_name LIKE ? OR phone LIKE ? OR email LIKE ?", v, v, v, v)
	}

	db = db.Order(gormx.ParseOrder(params.OrderFields.AddIdSortField()))

	var list []*Model
	pr, err := gormx.WrapPageQuery(ctx, db, params.PaginationParam, &list)
	if err != nil {
		return nil, nil, errors.WithStack(err)
	}
	return toDomainList(list), pr, nil
}

func (a *repository) Get(ctx context.Context, id string) (*user.User, error) {
	item := &Model{}
	ok, err := gormx.FindOne(ctx, GetModelDB(ctx, a.db).Where("id=?", id), item)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if !ok {
		return nil, nil
	}

	return item.ToDomain(), nil
}

func (a *repository) Create(ctx context.Context, item *user.User) error {
	result := GetModelDB(ctx, a.db).Create(domainToModel(item))
	return errors.WithStack(result.Error)
}

func (a *repository) Update(ctx context.Context, id string, item *user.User) error {
	result := GetModelDB(ctx, a.db).Where("id=?", id).Updates(domainToModel(item))
	return errors.WithStack(result.Error)
}

func (a *repository) Delete(ctx context.Context, id string) error {
	result := GetModelDB(ctx, a.db).Where("id=?", id).Delete(Model{})
	return errors.WithStack(result.Error)
}

func (a *repository) UpdateStatus(ctx context.Context, id string, status int) error {
	result := GetModelDB(ctx, a.db).Where("id=?", id).Update("status", status)
	return errors.WithStack(result.Error)
}

func (a *repository) UpdatePassword(ctx context.Context, id, password string) error {
	result := GetModelDB(ctx, a.db).Where("id=?", id).Update("password", password)
	return errors.WithStack(result.Error)
}
