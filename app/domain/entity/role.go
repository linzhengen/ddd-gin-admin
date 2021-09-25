package entity

import (
	"time"
)

type Role struct {
	ID        string     `gorm:"column:id;primary_key;size:36;"`
	Name      string     `gorm:"column:name;size:100;index;default:'';not null;"`
	Sequence  int        `gorm:"column:sequence;index;default:0;not null;"`
	Memo      *string    `gorm:"column:memo;size:1024;"`
	Status    int        `gorm:"column:status;index;default:0;not null;"`
	Creator   string     `gorm:"column:creator;size:36;"`
	CreatedAt time.Time  `gorm:"column:created_at;index;"`
	UpdatedAt time.Time  `gorm:"column:updated_at;index;"`
	DeletedAt *time.Time `gorm:"column:deleted_at;index;"`
}

type Roles []*Role