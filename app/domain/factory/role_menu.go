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

func (a RoleMenu) ToEntity(RoleMenu *schema.RoleMenu) *entity.RoleMenu {
	item := new(entity.RoleMenu)
	structure.Copy(RoleMenu, item)
	return item
}

func (a RoleMenu) ToSchema(RoleMenu *entity.RoleMenu) *schema.RoleMenu {
	item := new(schema.RoleMenu)
	structure.Copy(RoleMenu, item)
	return item
}

func (a RoleMenu) ToSchemaList(RoleMenus []*entity.RoleMenu) schema.RoleMenus {
	list := make([]*schema.RoleMenu, len(RoleMenus))
	for i, item := range RoleMenus {
		list[i] = a.ToSchema(item)
	}
	return list
}
