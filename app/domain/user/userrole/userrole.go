package userrole

import "github.com/linzhengen/ddd-gin-admin/app/domain/pagination"

type UserRole struct {
	ID     string
	UserID string
	RoleID string
}

type QueryParam struct {
	PaginationParam pagination.Param
	OrderFields     pagination.OrderFields
	UserID          string
	UserIDs         []string
}

type UserRoles []*UserRole

func (a UserRoles) ToRoleIDs() []string {
	ids := make([]string, len(a))
	for i, item := range a {
		ids[i] = item.RoleID
	}
	return ids
}

func (a UserRoles) ToUserIDMap() map[string]UserRoles {
	m := make(map[string]UserRoles)
	for _, item := range a {
		m[item.UserID] = append(m[item.UserID], item)
	}
	return m
}

func (a UserRoles) ToMap() map[string]*UserRole {
	m := make(map[string]*UserRole)
	for _, item := range a {
		m[item.RoleID] = item
	}
	return m
}
