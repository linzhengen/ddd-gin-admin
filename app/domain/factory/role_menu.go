package factory

import (
	"github.com/linzhengen/ddd-gin-admin/app/domain/entity"
	"github.com/linzhengen/ddd-gin-admin/app/domain/valueobject/schema"
	"github.com/linzhengen/ddd-gin-admin/pkg/util/structure"
)

func NewRoleMenu() RoleMenu {
	return RoleMenu{}
}

type RoleMenu struct{}

func (a RoleMenu) ToEntity(roleMenu *schema.RoleMenu) *entity.RoleMenu {
	item := new(entity.RoleMenu)
	structure.Copy(roleMenu, item)
	return item
}

func (a RoleMenu) ToSchema(roleMenu *entity.RoleMenu) *schema.RoleMenu {
	item := new(schema.RoleMenu)
	structure.Copy(roleMenu, item)
	return item
}

func (a RoleMenu) ToSchemaList(roleMenus []*entity.RoleMenu) schema.RoleMenus {
	list := make([]*schema.RoleMenu, len(roleMenus))
	for i, item := range roleMenus {
		list[i] = a.ToSchema(item)
	}
	return list
}
