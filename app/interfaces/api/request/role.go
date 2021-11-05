package request

import (
	"github.com/linzhengen/ddd-gin-admin/app/domain/pagination"
)

type RoleQueryParam struct {
	PaginationParam pagination.Param
	OrderFields     pagination.OrderFields
	IDs             []string `form:"-"`          // IDs
	Name            string   `form:"-"`          // Name
	QueryValue      string   `form:"queryValue"` // Query Search Values
	UserID          string   `form:"-"`          // User ID
	Status          int      `form:"status"`     // Status(1:enable 2:disable)
}

type RoleMenuQueryParam struct {
	PaginationParam pagination.Param
	OrderFields     pagination.OrderFields
	RoleID          string
	RoleIDs         []string
}
