package persistence

import (
	"context"

	"github.com/linzhengen/ddd-gin-admin/infrastructure/gormx"

	"github.com/linzhengen/ddd-gin-admin/domain/repository"

	"github.com/jinzhu/gorm"
	"github.com/linzhengen/ddd-gin-admin/domain/entity"
	"github.com/linzhengen/ddd-gin-admin/domain/schema"
	"github.com/linzhengen/ddd-gin-admin/pkg/errors"
)

func NewUser(db *gorm.DB) repository.UserRepository {
	return &user{
		db: db,
	}
}

type user struct {
	db *gorm.DB
}

func (a *user) getQueryOption(opts ...schema.UserQueryOptions) schema.UserQueryOptions {
	var opt schema.UserQueryOptions
	if len(opts) > 0 {
		opt = opts[0]
	}
	return opt
}

func (a *user) Query(ctx context.Context, params schema.UserQueryParam, opts ...schema.UserQueryOptions) (*schema.UserQueryResult, error) {
	opt := a.getQueryOption(opts...)

	db := entity.GetUserDB(ctx, a.db)
	if v := params.UserName; v != "" {
		db = db.Where("user_name=?", v)
	}
	if v := params.Status; v > 0 {
		db = db.Where("status=?", v)
	}
	if v := params.RoleIDs; len(v) > 0 {
		subQuery := entity.GetUserRoleDB(ctx, a.db).
			Select("user_id").
			Where("role_id IN (?)", v).
			SubQuery()
		db = db.Where("id IN ?", subQuery)
	}
	if v := params.QueryValue; v != "" {
		v = "%" + v + "%"
		db = db.Where("user_name LIKE ? OR real_name LIKE ? OR phone LIKE ? OR email LIKE ?", v, v, v, v)
	}

	opt.OrderFields = append(opt.OrderFields, schema.NewOrderField("id", schema.OrderByDESC))
	db = db.Order(gormx.ParseOrder(opt.OrderFields))

	var list entity.Users
	pr, err := gormx.WrapPageQuery(ctx, db, params.PaginationParam, &list)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	qr := &schema.UserQueryResult{
		PageResult: pr,
		Data:       list.ToSchemaUsers(),
	}
	return qr, nil
}

func (a *user) Get(ctx context.Context, id string, opts ...schema.UserQueryOptions) (*schema.User, error) {
	var item entity.User
	ok, err := gormx.FindOne(ctx, entity.GetUserDB(ctx, a.db).Where("id=?", id), &item)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if !ok {
		return nil, nil
	}

	return item.ToSchemaUser(), nil
}

func (a *user) Create(ctx context.Context, item schema.User) error {
	sitem := entity.SchemaUser(item)
	result := entity.GetUserDB(ctx, a.db).Create(sitem.ToUser())
	return errors.WithStack(result.Error)
}

func (a *user) Update(ctx context.Context, id string, item schema.User) error {
	eitem := entity.SchemaUser(item).ToUser()
	result := entity.GetUserDB(ctx, a.db).Where("id=?", id).Updates(eitem)
	return errors.WithStack(result.Error)
}

func (a *user) Delete(ctx context.Context, id string) error {
	result := entity.GetUserDB(ctx, a.db).Where("id=?", id).Delete(entity.User{})
	return errors.WithStack(result.Error)
}

func (a *user) UpdateStatus(ctx context.Context, id string, status int) error {
	result := entity.GetUserDB(ctx, a.db).Where("id=?", id).Update("status", status)
	return errors.WithStack(result.Error)
}

func (a *user) UpdatePassword(ctx context.Context, id, password string) error {
	result := entity.GetUserDB(ctx, a.db).Where("id=?", id).Update("password", password)
	return errors.WithStack(result.Error)
}
