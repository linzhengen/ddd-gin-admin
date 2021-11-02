package menuaction

import "github.com/linzhengen/ddd-gin-admin/app/domain/pagination"

type MenuAction struct {
	ID     string
	MenuID string
	Code   string
	Name   string
}

type MenuActions []*MenuAction

type QueryParam struct {
	PaginationParam pagination.Param
	OrderFields     pagination.OrderFields
	MenuID          string
	IDs             []string
}

func (a MenuActions) ToMenuIDMap() map[string]MenuActions {
	m := make(map[string]MenuActions)
	for _, item := range a {
		m[item.MenuID] = append(m[item.MenuID], item)
	}
	return m
}
