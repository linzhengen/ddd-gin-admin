package entity

import (
	"context"

	"github.com/jinzhu/gorm"
	"github.com/linzhengen/ddd-gin-admin/domain/schema"
	"github.com/linzhengen/ddd-gin-admin/pkg/util/structure"
)

func GetUserRoleDB(ctx context.Context, defDB *gorm.DB) *gorm.DB {
	return GetDBWithModel(ctx, defDB, new(UserRole))
}

type SchemaUserRole schema.UserRole

func (a SchemaUserRole) ToUserRole() *UserRole {
	item := new(UserRole)
	//nolint:errcheck
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
	//nolint:errcheck
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
