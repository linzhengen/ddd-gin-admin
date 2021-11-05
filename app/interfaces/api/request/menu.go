package request

import (
	"github.com/linzhengen/ddd-gin-admin/app/domain/pagination"
)

type MenuQueryParam struct {
	PaginationParam  pagination.Param
	OrderFields      pagination.OrderFields
	IDs              []string `form:"-"`
	Name             string   `form:"-"`
	PrefixParentPath string   `form:"-"`
	QueryValue       string   `form:"queryValue"`
	ParentID         *string  `form:"parentID"`
	ShowStatus       int      `form:"showStatus"` // 1:show 2:hide
	Status           int      `form:"status"`     // 1:enable 2:disable
}

type MenuActionQueryParam struct {
	PaginationParam pagination.Param
	OrderFields     pagination.OrderFields
	MenuID          string
	IDs             []string
}

type MenuActionResourceQueryParam struct {
	PaginationParam pagination.Param
	OrderFields     pagination.OrderFields
	MenuID          string   // menu id
	MenuIDs         []string // menu ids
}
