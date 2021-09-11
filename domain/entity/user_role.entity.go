package entity

import (
	"github.com/linzhengen/ddd-gin-admin/domain/schema"
	"github.com/linzhengen/ddd-gin-admin/pkg/util/structure"
)

type SchemaUserRole schema.UserRole

func (a SchemaUserRole) ToUserRole() *UserRole {
	item := new(UserRole)
	structure.Copy(a, item)
	return item
}

type UserRole struct {
	ID     string `gorm:"column:id;primary_key;size:36;"`
	UserID string `gorm:"column:user_id;size:36;index;default:'';not null;"`
	RoleID string `gorm:"column:role_id;size:36;index;default:'';not null;"`
}

func (a UserRole) ToSchemaUserRole() *schema.UserRole {
	item := new(schema.UserRole)
	structure.Copy(a, item)
	return item
}

type UserRoles []*UserRole

func (a UserRoles) ToSchemaUserRoles() []*schema.UserRole {
	list := make([]*schema.UserRole, len(a))
	for i, item := range a {
		list[i] = item.ToSchemaUserRole()
	}
	return list
}
