package rolemenu

import "github.com/linzhengen/ddd-gin-admin/app/domain/pagination"

type RoleMenu struct {
	ID       string `gorm:"column:id;primary_key;size:36;"`
	RoleID   string `gorm:"column:role_id;size:36;index;default:'';not null;"`
	MenuID   string `gorm:"column:menu_id;size:36;index;default:'';not null;"`
	ActionID string `gorm:"column:action_id;size:36;index;default:'';not null;"`
}

type RoleMenus []*RoleMenu

type QueryParam struct {
	PaginationParam pagination.Param
	OrderFields     pagination.OrderFields
	RoleID          string
	RoleIDs         []string
}

func (a RoleMenus) ToMap() map[string]*RoleMenu {
	m := make(map[string]*RoleMenu)
	for _, item := range a {
		m[item.MenuID+"-"+item.ActionID] = item
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
	var idList []string
	m := make(map[string]struct{})

	for _, item := range a {
		if _, ok := m[item.ActionID]; ok {
			continue
		}
		idList = append(idList, item.ActionID)
		m[item.ActionID] = struct{}{}
	}

	return idList
}

func (a RoleMenus) ToRoleIDMap() map[string]RoleMenus {
	m := make(map[string]RoleMenus)
	for _, item := range a {
		m[item.RoleID] = append(m[item.RoleID], item)
	}
	return m
}
