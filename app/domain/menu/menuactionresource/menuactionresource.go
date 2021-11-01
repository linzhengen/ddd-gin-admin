package menuactionresource

import "github.com/linzhengen/ddd-gin-admin/app/domain/pagination"

type MenuActionResource struct {
	ID       string
	ActionID string
	Method   string
	Path     string
}

type QueryParam struct {
	PaginationParam pagination.Param
	OrderFields     pagination.OrderFields
	MenuID          string
	MenuIDs         []string
}
