package user

import (
	"time"

	"github.com/linzhengen/ddd-gin-admin/pkg/util/structure"

	"github.com/linzhengen/ddd-gin-admin/app/domain/user/userrole"

	"github.com/linzhengen/ddd-gin-admin/app/domain/pagination"
	"github.com/linzhengen/ddd-gin-admin/app/domain/user/role"
)

type User struct {
	ID        string
	UserName  string
	RealName  string
	Password  string
	Email     *string
	Phone     *string
	Status    int
	Creator   string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
	Roles     role.Roles
}

func (a User) FillRoles(userRoles map[string]userrole.UserRoles, roles map[string]*role.Role) *User {
	u := new(User)
	structure.Copy(a, u)
	for _, roleID := range userRoles[a.ID].ToRoleIDs() {
		if v, ok := roles[roleID]; ok {
			u.Roles = append(u.Roles, v)
		}
	}
	return u
}

type Users []*User

type QueryParams struct {
	PaginationParam pagination.Param
	OrderFields     pagination.OrderFields
	UserName        string
	QueryValue      string
	Status          int
	RoleIDs         []string
}

func (a Users) ToIDs() []string {
	idList := make([]string, len(a))
	for i, item := range a {
		idList[i] = item.ID
	}
	return idList
}

func (a Users) FillRoles(userRoles map[string]userrole.UserRoles, roles map[string]*role.Role) Users {
	list := make(Users, len(a))
	for i, item := range a {
		list[i] = item.FillRoles(userRoles, roles)
	}

	return list
}
