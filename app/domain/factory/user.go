package factory

import (
	"github.com/linzhengen/ddd-gin-admin/app/domain/entity"
	"github.com/linzhengen/ddd-gin-admin/app/domain/valueobject/schema"
	"github.com/linzhengen/ddd-gin-admin/pkg/util/structure"
)

func NewUser() User {
	return User{}
}

type User struct{}

func (a User) ToEntity(User *schema.User) *entity.User {
	item := new(entity.User)
	structure.Copy(User, item)
	return item
}

func (a User) ToSchema(User *entity.User) *schema.User {
	item := new(schema.User)
	structure.Copy(User, item)
	return item
}

func (a User) ToSchemaList(Users []*entity.User) schema.Users {
	list := make([]*schema.User, len(Users))
	for i, item := range Users {
		list[i] = a.ToSchema(item)
	}
	return list
}
