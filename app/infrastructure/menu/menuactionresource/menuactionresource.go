package menuactionresource

import (
	"github.com/linzhengen/ddd-gin-admin/app/domain/menu/menuactionresource"
	"github.com/linzhengen/ddd-gin-admin/pkg/util/structure"
)

type Model struct {
	ID       string `gorm:"column:id;primary_key;size:36;"`
	ActionID string `gorm:"column:action_id;size:36;index;default:'';not null;"`
	Method   string `gorm:"column:method;size:100;default:'';not null;"`
	Path     string `gorm:"column:path;size:100;default:'';not null;"`
}

func (Model) TableName() string {
	return "menu_action_resources"
}

func (a Model) ToDomain() *menuactionresource.MenuActionResource {
	item := new(menuactionresource.MenuActionResource)
	structure.Copy(a, item)
	return item
}

func toDomainList(ms []*Model) []*menuactionresource.MenuActionResource {
	list := make([]*menuactionresource.MenuActionResource, len(ms))
	for i, item := range ms {
		list[i] = item.ToDomain()
	}
	return list
}

func domainToModel(m *menuactionresource.MenuActionResource) *Model {
	item := new(Model)
	structure.Copy(m, item)
	return item
}
