package role

import (
	"context"

	"github.com/linzhengen/ddd-gin-admin/app/domain/user/role"

	"github.com/linzhengen/ddd-gin-admin/app/infrastructure/user/userrole"

	"github.com/linzhengen/ddd-gin-admin/app/domain/errors"

	"github.com/jinzhu/gorm"
	"github.com/linzhengen/ddd-gin-admin/app/domain/pagination"
	"github.com/linzhengen/ddd-gin-admin/app/infrastructure/gormx"
)

func NewRepository(db *gorm.DB) role.Repository {
	return &repository{
		db: db,
	}
}

type repository struct {
	db *gorm.DB
}

func GetModelDB(ctx context.Context, db *gorm.DB) *gorm.DB {
	return gormx.GetDBWithModel(ctx, db, new(Model))
}

func (a *repository) Query(ctx context.Context, params role.QueryParam) (role.Roles, *pagination.Pagination, error) {
	db := GetModelDB(ctx, a.db)
	if v := params.IDs; len(v) > 0 {
		db = db.Where("id IN (?)", v)
	}
	if v := params.Name; v != "" {
		db = db.Where("name=?", v)
	}
	if v := params.UserID; v != "" {
		// todo: serviceへ移動
		subQuery := userrole.GetModelDB(ctx, a.db).
			Where("deleted_at is null").
			Where("user_id=?", v).
			Select("role_id").SubQuery()
		db = db.Where("id IN ?", subQuery)
	}
	if v := params.QueryValue; v != "" {
		v = "%" + v + "%"
		db = db.Where("name LIKE ? OR memo LIKE ?", v, v)
	}

	db = db.Order(gormx.ParseOrder(params.OrderFields.AddIdSortField()))

	var list []*Model
	pr, err := gormx.WrapPageQuery(ctx, db, params.PaginationParam, &list)
	if err != nil {
		return nil, nil, errors.WithStack(err)
	}

	return toDomainList(list), pr, nil
}

func (a *repository) Get(ctx context.Context, id string) (*role.Role, error) {
	role := &Model{}
	ok, err := gormx.FindOne(ctx, GetModelDB(ctx, a.db).Where("id=?", id), &role)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if !ok {
		return nil, nil
	}

	return role.ToDomain(), nil
}

func (a *repository) Create(ctx context.Context, item *role.Role) error {
	result := GetModelDB(ctx, a.db).Create(domainToModel(item))
	return errors.WithStack(result.Error)
}

func (a *repository) Update(ctx context.Context, id string, item *role.Role) error {
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
