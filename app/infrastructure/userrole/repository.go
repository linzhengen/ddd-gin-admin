package userrole

import (
	"context"

	"github.com/linzhengen/ddd-gin-admin/app/domain/errors"

	"github.com/jinzhu/gorm"
	"github.com/linzhengen/ddd-gin-admin/app/domain/pagination"
	"github.com/linzhengen/ddd-gin-admin/app/domain/userrole"
	"github.com/linzhengen/ddd-gin-admin/app/infrastructure/gormx"
)

func NewRepository(db *gorm.DB) userrole.Repository {
	return &repository{
		db: db,
	}
}

type repository struct {
	db *gorm.DB
}

func (a *repository) Query(ctx context.Context, params userrole.QueryParam) ([]*userrole.UserRole, *pagination.Pagination, error) {
	db := GetModelDB(ctx, a.db)
	if v := params.UserID; v != "" {
		db = db.Where("user_id=?", v)
	}
	if v := params.UserIDs; len(v) > 0 {
		db = db.Where("user_id IN (?)", v)
	}

	db = db.Order(gormx.ParseOrder(params.OrderFields.AddIdSortField()))

	var list []*Model
	pr, err := gormx.WrapPageQuery(ctx, db, params.PaginationParam, &list)
	if err != nil {
		return nil, nil, errors.WithStack(err)
	}
	return toDomainList(list), pr, nil
}

func (a *repository) Get(ctx context.Context, id string) (*userrole.UserRole, error) {
	db := GetModelDB(ctx, a.db).Where("id=?", id)
	item := &Model{}
	ok, err := gormx.FindOne(ctx, db, &item)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if !ok {
		return nil, nil
	}

	return item.ToDomain(), nil
}

func (a *repository) Create(ctx context.Context, item *userrole.UserRole) error {
	result := GetModelDB(ctx, a.db).Create(domainToModel(item))
	return errors.WithStack(result.Error)
}

func (a *repository) Update(ctx context.Context, id string, item *userrole.UserRole) error {
	result := GetModelDB(ctx, a.db).Where("id=?", id).Updates(domainToModel(item))
	return errors.WithStack(result.Error)
}

func (a *repository) Delete(ctx context.Context, id string) error {
	result := GetModelDB(ctx, a.db).Where("id=?", id).Delete(Model{})
	return errors.WithStack(result.Error)
}

func (a *repository) DeleteByUserID(ctx context.Context, userID string) error {
	result := GetModelDB(ctx, a.db).Where("user_id=?", userID).Delete(Model{})
	return errors.WithStack(result.Error)
}
