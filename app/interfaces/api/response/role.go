package response

import (
	"time"

	"github.com/linzhengen/ddd-gin-admin/app/domain/user/role"
	"github.com/linzhengen/ddd-gin-admin/pkg/util/structure"

	"github.com/linzhengen/ddd-gin-admin/app/domain/pagination"
)

type Role struct {
	ID        string    `json:"id"`         // ID
	Name      string    `json:"name"`       // Name
	Sequence  int       `json:"sequence"`   // Sequence
	Memo      string    `json:"memo"`       // Memo
	Status    int       `json:"status"`     // Status(1:enable 2:disable)
	Creator   string    `json:"creator"`    // Creator
	CreatedAt time.Time `json:"created_at"` // CreatedAt
	UpdatedAt time.Time `json:"updated_at"` // UpdatedAt
	RoleMenus RoleMenus `json:"role_menus"` // RoleMenus
}

func RoleFromDomain(role *role.Role) *Role {
	item := new(Role)
	structure.Copy(role, item)
	return item
}

type RoleQueryResult struct {
	Data       Roles
	PageResult *pagination.Pagination
}

type Roles []*Role

func RolesFromDomain(roles role.Roles) Roles {
	list := make([]*Role, len(roles))
	for i, item := range roles {
		structure.Copy(item, list[i])
	}
	return list
}

func (a Roles) ToNames() []string {
	names := make([]string, len(a))
	for i, item := range a {
		names[i] = item.Name
	}
	return names
}

func (a Roles) ToMap() map[string]*Role {
	m := make(map[string]*Role)
	for _, item := range a {
		m[item.ID] = item
	}
	return m
}

// ----------------------------------------RoleMenu--------------------------------------

type RoleMenu struct {
	ID       string `json:"id"`                           // ID
	RoleID   string `json:"role_id" binding:"required"`   // Role ID
	MenuID   string `json:"menu_id" binding:"required"`   // Menu ID
	ActionID string `json:"action_id" binding:"required"` // Action ID
}

type RoleMenuQueryResult struct {
	Data       RoleMenus
	PageResult *pagination.Pagination
}

type RoleMenus []*RoleMenu

func (a RoleMenus) ToMap() map[string]*RoleMenu {
	m := make(map[string]*RoleMenu)
	for _, item := range a {
		m[item.MenuID+"-"+item.ActionID] = item
	}
	return m
}

func (a RoleMenus) ToRoleIDMap() map[string]RoleMenus {
	m := make(map[string]RoleMenus)
	for _, item := range a {
		m[item.RoleID] = append(m[item.RoleID], item)
	}
	return m
}

func (a RoleMenus) ToMenuIDs() []string {
	var idList []string
	m := make(map[string]struct{})

	for _, item := range a {
		if _, ok := m[item.MenuID]; ok {
			continue
		}
		idList = append(idList, item.MenuID)
		m[item.MenuID] = struct{}{}
	}

	return idList
}

func (a RoleMenus) ToActionIDs() []string {
	idList := make([]string, len(a))
	m := make(map[string]struct{})
	for i, item := range a {
		if _, ok := m[item.ActionID]; ok {
			continue
		}
		idList[i] = item.ActionID
		m[item.ActionID] = struct{}{}
	}
	return idList
}
