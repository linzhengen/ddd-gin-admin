package menuaction

import (
	"github.com/linzhengen/ddd-gin-admin/app/domain/menu/menuaction"
	"github.com/linzhengen/ddd-gin-admin/pkg/util/structure"
)

type Model struct {
	ID     string `gorm:"column:id;primary_key;size:36;"`
	MenuID string `gorm:"column:menu_id;size:36;index;default:'';not null;"`
	Code   string `gorm:"column:code;size:100;default:'';not null;"`
	Name   string `gorm:"column:name;size:100;default:'';not null;"`
}

func (Model) TableName() string {
	return "menu_actions"
}

func (a Model) ToDomain() *menuaction.MenuAction {
	item := new(menuaction.MenuAction)
	structure.Copy(a, item)
	return item
}

func toDomainList(ms []*Model) []*menuaction.MenuAction {
	list := make([]*menuaction.MenuAction, len(ms))
	for i, item := range ms {
		list[i] = item.ToDomain()
	}
	return list
}

func domainToModel(m *menuaction.MenuAction) *Model {
	item := new(Model)
	structure.Copy(m, item)
	return item
}
