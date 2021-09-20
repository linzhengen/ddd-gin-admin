package factory

import (
	"github.com/linzhengen/ddd-gin-admin/app/domain/entity"
	"github.com/linzhengen/ddd-gin-admin/app/domain/valueobject/schema"

	"github.com/linzhengen/ddd-gin-admin/pkg/util/structure"
)

func NewMenu() Menu {
	return Menu{}
}

type Menu struct{}

func (a Menu) ToEntity(menu *schema.Menu) *entity.Menu {
	item := new(entity.Menu)
	structure.Copy(menu, item)
	return item
}

func (a Menu) ToSchema(menu *entity.Menu) *schema.Menu {
	item := new(schema.Menu)
	structure.Copy(menu, item)
	return item
}

func (a Menu) ToSchemaList(menus []*entity.Menu) schema.Menus {
	list := make([]*schema.Menu, len(menus))
	for i, item := range menus {
		list[i] = a.ToSchema(item)
	}
	return list
}
