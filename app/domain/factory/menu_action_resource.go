package factory

import (
	"github.com/linzhengen/ddd-gin-admin/app/domain/entity"
	"github.com/linzhengen/ddd-gin-admin/app/domain/valueobject/schema"
	"github.com/linzhengen/ddd-gin-admin/pkg/util/structure"
)

func NewMenuActionResource() MenuActionResource {
	return MenuActionResource{}
}

type MenuActionResource struct{}

func (a MenuActionResource) ToEntity(menuActionResource *schema.MenuActionResource) *entity.MenuActionResource {
	item := new(entity.MenuActionResource)
	structure.Copy(menuActionResource, item)
	return item
}

func (a MenuActionResource) ToSchema(menuActionResource *entity.MenuActionResource) *schema.MenuActionResource {
	item := new(schema.MenuActionResource)
	structure.Copy(menuActionResource, item)
	return item
}

func (a MenuActionResource) ToSchemaList(menuActionResources []*entity.MenuActionResource) schema.MenuActionResources {
	list := make([]*schema.MenuActionResource, len(menuActionResources))
	for i, item := range menuActionResources {
		list[i] = a.ToSchema(item)
	}
	return list
}
