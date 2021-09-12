package entity

import (
	"github.com/linzhengen/ddd-gin-admin/app/domain/schema"
	"github.com/linzhengen/ddd-gin-admin/pkg/util/structure"
)

type SchemaRoleMenu schema.RoleMenu

func (a SchemaRoleMenu) ToRoleMenu() *RoleMenu {
	item := new(RoleMenu)
	structure.Copy(a, item)
	return item
}

type RoleMenu struct {
	ID       string `gorm:"column:id;primary_key;size:36;"`
	RoleID   string `gorm:"column:role_id;size:36;index;default:'';not null;"`
	MenuID   string `gorm:"column:menu_id;size:36;index;default:'';not null;"`
	ActionID string `gorm:"column:action_id;size:36;index;default:'';not null;"`
}

func (a RoleMenu) ToSchemaRoleMenu() *schema.RoleMenu {
	item := new(schema.RoleMenu)
	structure.Copy(a, item)
	return item
}

type RoleMenus []*RoleMenu

func (a RoleMenus) ToSchemaRoleMenus() []*schema.RoleMenu {
	list := make([]*schema.RoleMenu, len(a))
	for i, item := range a {
		list[i] = item.ToSchemaRoleMenu()
	}
	return list
}
