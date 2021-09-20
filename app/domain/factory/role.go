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

func (a Role) ToEntity(Role *schema.Role) *entity.Role {
	item := new(entity.Role)
	structure.Copy(Role, item)
	return item
}

func (a Role) ToSchema(Role *entity.Role) *schema.Role {
	item := new(schema.Role)
	structure.Copy(Role, item)
	return item
}

func (a Role) ToSchemaList(Roles []*entity.Role) schema.Roles {
	list := make([]*schema.Role, len(Roles))
	for i, item := range Roles {
		list[i] = a.ToSchema(item)
	}
	return list
}
