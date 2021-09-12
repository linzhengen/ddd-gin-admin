package entity

import (
	"github.com/linzhengen/ddd-gin-admin/app/domain/schema"
	"github.com/linzhengen/ddd-gin-admin/pkg/util/structure"
)

type SchemaMenuActionResource schema.MenuActionResource

func (a SchemaMenuActionResource) ToMenuActionResource() *MenuActionResource {
	item := new(MenuActionResource)
	structure.Copy(a, item)
	return item
}

type MenuActionResource struct {
	ID       string `gorm:"column:id;primary_key;size:36;"`
	ActionID string `gorm:"column:action_id;size:36;index;default:'';not null;"`
	Method   string `gorm:"column:method;size:100;default:'';not null;"`
	Path     string `gorm:"column:path;size:100;default:'';not null;"`
}

func (a MenuActionResource) ToSchemaMenuActionResource() *schema.MenuActionResource {
	item := new(schema.MenuActionResource)
	structure.Copy(a, item)
	return item
}

type MenuActionResources []*MenuActionResource

func (a MenuActionResources) ToSchemaMenuActionResources() []*schema.MenuActionResource {
	list := make([]*schema.MenuActionResource, len(a))
	for i, item := range a {
		list[i] = item.ToSchemaMenuActionResource()
	}
	return list
}
