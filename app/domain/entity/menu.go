package entity

import (
	"time"
)

type Menu struct {
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

type Menus []*Menu
