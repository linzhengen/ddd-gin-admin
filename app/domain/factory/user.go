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

func (a User) ToEntity(user *schema.User) *entity.User {
	item := new(entity.User)
	structure.Copy(user, item)
	return item
}

func (a User) ToSchema(user *entity.User) *schema.User {
	item := new(schema.User)
	structure.Copy(user, item)
	return item
}

func (a User) ToSchemaList(users []*entity.User) schema.Users {
	list := make([]*schema.User, len(users))
	for i, item := range users {
		list[i] = a.ToSchema(item)
	}
	return list
}
