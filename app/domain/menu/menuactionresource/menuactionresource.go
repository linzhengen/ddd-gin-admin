package menuactionresource

import "github.com/linzhengen/ddd-gin-admin/app/domain/pagination"

type MenuActionResource struct {
	ID       string
	ActionID string
	Method   string
	Path     string
}
type MenuActionResources []*MenuActionResource

type QueryParam struct {
	PaginationParam pagination.Param
	OrderFields     pagination.OrderFields
	MenuID          string
	MenuIDs         []string
}

func (a MenuActionResources) ToMenuActionIDMap() map[string]MenuActionResources {
	m := make(map[string]MenuActionResources)
	for _, item := range a {
		m[item.ActionID] = append(m[item.ActionID], item)
	}
	return m
}
