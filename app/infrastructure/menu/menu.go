package menu

import (
	"time"

	"github.com/linzhengen/ddd-gin-admin/app/domain/menu"
)

type Model struct {
	ID         string     `gorm:"column:id;primary_key;size:36;"`
	Name       string     `gorm:"column:name;size:50;index;default:'';not null;"`
	Sequence   int        `gorm:"column:sequence;index;default:0;not null;"`
	Icon       *string    `gorm:"column:icon;size:255;"`
	Router     *string    `gorm:"column:router;size:255;"`
	ParentID   *string    `gorm:"column:parent_id;size:36;index;"`
	ParentPath *string    `gorm:"column:parent_path;size:518;index;"`
	ShowStatus int        `gorm:"column:show_status;index;default:0;not null;"`
	Status     int        `gorm:"column:status;index;default:0;not null;"`
	Memo       *string    `gorm:"column:memo;size:1024;"`
	Creator    string     `gorm:"column:creator;size:36;"`
	CreatedAt  time.Time  `gorm:"column:created_at;index;"`
	UpdatedAt  time.Time  `gorm:"column:updated_at;index;"`
	DeletedAt  *time.Time `gorm:"column:deleted_at;index;"`
}

func (Model) TableName() string {
	return "menus"
}

func (a Model) ToDomain() *menu.Menu {
	item := new(menu.Menu)
	item.ID = a.ID
	item.Name = a.Name
	item.Sequence = a.Sequence
	if a.Icon != nil {
		item.Icon = *a.Icon
	}
	if a.Router != nil {
		item.Router = *a.Router
	}
	if a.ParentID != nil {
		item.ParentID = *a.ParentID
	}
	if a.ParentPath != nil {
		item.ParentPath = *a.ParentPath
	}
	item.ShowStatus = a.ShowStatus
	item.Status = a.Status
	if a.Memo != nil {
		item.Memo = *a.Memo
	}
	item.Creator = a.Creator
	item.CreatedAt = a.CreatedAt
	item.UpdatedAt = a.UpdatedAt
	item.DeletedAt = a.DeletedAt
	return item
}

func toDomainList(menus []*Model) []*menu.Menu {
	list := make([]*menu.Menu, len(menus))
	for i, item := range menus {
		list[i] = item.ToDomain()
	}
	return list
}

func domainToModel(m *menu.Menu) *Model {
	item := new(Model)
	item.ID = m.ID
	item.Name = m.Name
	item.Sequence = m.Sequence
	if m.Icon != "" {
		item.Icon = &m.Icon
	}
	if m.Router != "" {
		item.Router = &m.Router
	}
	if m.ParentID != "" {
		item.ParentID = &m.ParentID
	}
	if m.ParentPath != "" {
		item.ParentPath = &m.ParentPath
	}
	item.ShowStatus = m.ShowStatus
	item.Status = m.Status
	if m.Memo != "" {
		item.Memo = &m.Memo
	}
	item.Creator = m.Creator
	item.CreatedAt = m.CreatedAt
	item.UpdatedAt = m.UpdatedAt
	item.DeletedAt = m.DeletedAt
	return item
}
