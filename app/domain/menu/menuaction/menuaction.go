package menuaction

import "github.com/linzhengen/ddd-gin-admin/app/domain/pagination"

type MenuAction struct {
	ID     string
	MenuID string
	Code   string
	Name   string
}

type QueryParam struct {
	PaginationParam pagination.Param
	OrderFields     pagination.OrderFields
	MenuID          string
	IDs             []string
}
