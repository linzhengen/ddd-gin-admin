package factory

import (
	"github.com/linzhengen/ddd-gin-admin/app/domain/entity"
	"github.com/linzhengen/ddd-gin-admin/app/domain/valueobject/schema"
	"github.com/linzhengen/ddd-gin-admin/pkg/util/structure"
)

func NewUserRole() UserRole {
	return UserRole{}
}

type UserRole struct{}

func (a UserRole) ToEntity(UserRole *schema.UserRole) *entity.UserRole {
	item := new(entity.UserRole)
	structure.Copy(UserRole, item)
	return item
}

func (a UserRole) ToSchema(UserRole *entity.UserRole) *schema.UserRole {
	item := new(schema.UserRole)
	structure.Copy(UserRole, item)
	return item
}

func (a UserRole) ToSchemaList(UserRoles []*entity.UserRole) schema.UserRoles {
	list := make([]*schema.UserRole, len(UserRoles))
	for i, item := range UserRoles {
		list[i] = a.ToSchema(item)
	}
	return list
}
