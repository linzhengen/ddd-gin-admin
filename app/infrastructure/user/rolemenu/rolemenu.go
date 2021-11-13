package rolemenu

import (
	"github.com/linzhengen/ddd-gin-admin/app/domain/user/rolemenu"
	"github.com/linzhengen/ddd-gin-admin/pkg/util/structure"
)

type Model struct {
	ID       string `gorm:"column:id;primary_key;size:36;"`
	RoleID   string `gorm:"column:role_id;size:36;index;default:'';not null;"`
	MenuID   string `gorm:"column:menu_id;size:36;index;default:'';not null;"`
	ActionID string `gorm:"column:action_id;size:36;index;default:'';not null;"`
}

func (Model) TableName() string {
	return "role_menus"
}

func (a Model) ToDomain() *rolemenu.RoleMenu {
	item := new(rolemenu.RoleMenu)
	structure.Copy(a, item)
	return item
}

func toDomainList(ms []*Model) []*rolemenu.RoleMenu {
	list := make([]*rolemenu.RoleMenu, len(ms))
	for i, item := range ms {
		list[i] = item.ToDomain()
	}
	return list
}

func domainToModel(u *rolemenu.RoleMenu) *Model {
	item := new(Model)
	structure.Copy(u, item)
	return item
}
