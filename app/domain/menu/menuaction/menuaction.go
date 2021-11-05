package menuaction

import (
	"github.com/linzhengen/ddd-gin-admin/app/domain/menu/menuactionresource"
	"github.com/linzhengen/ddd-gin-admin/app/domain/pagination"
)

type MenuAction struct {
	ID        string
	MenuID    string
	Code      string
	Name      string
	Resources menuactionresource.MenuActionResources
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

func (a MenuActions) FillResources(mResources map[string]menuactionresource.MenuActionResources) {
	for i, item := range a {
		a[i].Resources = mResources[item.ID]
	}
}

func (a MenuActions) ToMap() map[string]*MenuAction {
	m := make(map[string]*MenuAction)
	for _, item := range a {
		m[item.Code] = item
	}
	return m
}
