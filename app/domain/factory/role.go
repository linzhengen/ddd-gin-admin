package factory

import (
	"github.com/linzhengen/ddd-gin-admin/app/domain/entity"
	"github.com/linzhengen/ddd-gin-admin/app/domain/valueobject/schema"

	"github.com/linzhengen/ddd-gin-admin/pkg/util/structure"
)

func NewRole() Role {
	return Role{}
}

type Role struct{}

func (a Role) ToEntity(role *schema.Role) *entity.Role {
	item := new(entity.Role)
	structure.Copy(role, item)
	return item
}

func (a Role) ToSchema(role *entity.Role) *schema.Role {
	item := new(schema.Role)
	structure.Copy(role, item)
	return item
}

func (a Role) ToSchemaList(roles []*entity.Role) schema.Roles {
	list := make([]*schema.Role, len(roles))
	for i, item := range roles {
		list[i] = a.ToSchema(item)
	}
	return list
}
