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

func (a UserRole) ToEntity(userRole *schema.UserRole) *entity.UserRole {
	item := new(entity.UserRole)
	structure.Copy(userRole, item)
	return item
}

func (a UserRole) ToSchema(userRole *entity.UserRole) *schema.UserRole {
	item := new(schema.UserRole)
	structure.Copy(userRole, item)
	return item
}

func (a UserRole) ToSchemaList(userRoles []*entity.UserRole) schema.UserRoles {
	list := make([]*schema.UserRole, len(userRoles))
	for i, item := range userRoles {
		list[i] = a.ToSchema(item)
	}
	return list
}
