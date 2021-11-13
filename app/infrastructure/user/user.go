package user

import (
	"time"

	"github.com/linzhengen/ddd-gin-admin/pkg/util/structure"

	"github.com/linzhengen/ddd-gin-admin/app/domain/user"
)

type Model struct {
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

func (Model) TableName() string {
	return "users"
}

func (a Model) ToDomain() *user.User {
	item := new(user.User)
	structure.Copy(a, item)
	return item
}

func toDomainList(ms []*Model) []*user.User {
	list := make([]*user.User, len(ms))
	for i, item := range ms {
		list[i] = item.ToDomain()
	}
	return list
}

func domainToModel(u *user.User) *Model {
	item := new(Model)
	structure.Copy(u, item)
	return item
}
