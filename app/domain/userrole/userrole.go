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
