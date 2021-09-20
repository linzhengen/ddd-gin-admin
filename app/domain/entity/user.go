package entity

import (
	"time"
)

type User struct {
	ID        string     `gorm:"column:id;primary_key;size:36;"`
	UserName  string     `gorm:"column:user_name;size:64;index;default:'';not null;"`
	RealName  string     `gorm:"column:real_name;size:64;index;default:'';not null;"`
	Password  string     `gorm:"column:password;size:40;default:'';not null;"`
	Email     *string    `gorm:"column:email;size:255;index;"`
	Phone     *string    `gorm:"column:phone;size:20;index;"`
	Status    int        `gorm:"column:status;index;default:0;not null;"`
	Creator   string     `gorm:"column:creator;size:36;"`
	CreatedAt time.Time  `gorm:"column:created_at;index;"`
	UpdatedAt time.Time  `gorm:"column:updated_at;index;"`
	DeletedAt *time.Time `gorm:"column:deleted_at;index;"`
}

type Users []*User
