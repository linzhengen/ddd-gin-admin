package entity

import (
	"github.com/linzhengen/ddd-gin-admin/app/domain/schema"
	"github.com/linzhengen/ddd-gin-admin/pkg/util/structure"
)

type SchemaMenuAction schema.MenuAction

func (a SchemaMenuAction) ToMenuAction() *MenuAction {
	item := new(MenuAction)
	structure.Copy(a, item)
	return item
}

type MenuAction struct {
	ID     string `gorm:"column:id;primary_key;size:36;"`
	MenuID string `gorm:"column:menu_id;size:36;index;default:'';not null;"`
	Code   string `gorm:"column:code;size:100;default:'';not null;"`
	Name   string `gorm:"column:name;size:100;default:'';not null;"`
}

func (a MenuAction) ToSchemaMenuAction() *schema.MenuAction {
	item := new(schema.MenuAction)
	structure.Copy(a, item)
	return item
}

type MenuActions []*MenuAction

func (a MenuActions) ToSchemaMenuActions() []*schema.MenuAction {
	list := make([]*schema.MenuAction, len(a))
	for i, item := range a {
		list[i] = item.ToSchemaMenuAction()
	}
	return list
}