package userrole

import (
	"context"

	"github.com/linzhengen/ddd-gin-admin/pkg/util/structure"

	"github.com/linzhengen/ddd-gin-admin/app/domain/userrole"

	"github.com/jinzhu/gorm"
	"github.com/linzhengen/ddd-gin-admin/app/infrastructure/gormx"
)

type Model struct {
	ID     string `gorm:"column:id;primary_key;size:36;"`
	UserID string `gorm:"column:user_id;size:36;index;default:'';not null;"`
	RoleID string `gorm:"column:role_id;size:36;index;default:'';not null;"`
}

func (a Model) ToDomain() *userrole.UserRole {
	item := new(userrole.UserRole)
	structure.Copy(a, item)
	return item
}

func toDomainList(userroles []*Model) []*userrole.UserRole {
	list := make([]*userrole.UserRole, len(userroles))
	for i, item := range userroles {
		list[i] = item.ToDomain()
	}
	return list
}

func domainToModel(u *userrole.UserRole) *Model {
	item := new(Model)
	structure.Copy(u, item)
	return item
}

func GetModelDB(ctx context.Context, defDB *gorm.DB) *gorm.DB {
	return gormx.GetDBWithModel(ctx, defDB, new(Model))
}
