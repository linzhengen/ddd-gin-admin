package request

import (
	"github.com/linzhengen/ddd-gin-admin/app/domain/pagination"
)

type UserQueryParam struct {
	PaginationParam pagination.Param
	OrderFields     pagination.OrderFields
	UserName        string   `form:"userName"`   // User Name
	QueryValue      string   `form:"queryValue"` // Query search values
	Status          int      `form:"status"`     // Status(1:enable 2:disable)
	RoleIDs         []string `form:"-"`          // Role IDs
}

type UserRoleQueryParam struct {
	PaginationParam pagination.Param
	OrderFields     pagination.OrderFields
	UserID          string
	UserIDs         []string
}
