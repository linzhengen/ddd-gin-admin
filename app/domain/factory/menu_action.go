package factory

import (
	"github.com/linzhengen/ddd-gin-admin/app/domain/entity"
	"github.com/linzhengen/ddd-gin-admin/app/domain/valueobject/schema"
	"github.com/linzhengen/ddd-gin-admin/pkg/util/structure"
)

func NewMenuAction() MenuAction {
	return MenuAction{}
}

type MenuAction struct{}

func (a MenuAction) ToEntity(menuAction *schema.MenuAction) *entity.MenuAction {
	item := new(entity.MenuAction)
	structure.Copy(menuAction, item)
	return item
}

func (a MenuAction) ToSchema(menuAction *entity.MenuAction) *schema.MenuAction {
	item := new(schema.MenuAction)
	structure.Copy(menuAction, item)
	return item
}

func (a MenuAction) ToSchemaList(menuActions []*entity.MenuAction) schema.MenuActions {
	list := make([]*schema.MenuAction, len(menuActions))
	for i, item := range menuActions {
		list[i] = a.ToSchema(item)
	}
	return list
}
