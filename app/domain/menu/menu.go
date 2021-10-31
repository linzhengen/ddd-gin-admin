package menu

import (
	"time"

	"github.com/linzhengen/ddd-gin-admin/app/domain/pagination"
)

type Menu struct {
	ID         string
	Name       string
	Sequence   int
	Icon       *string
	Router     *string
	ParentID   *string
	ParentPath *string
	ShowStatus int
	Status     int
	Memo       *string
	Creator    string
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  *time.Time
}

type QueryParam struct {
	PaginationParam  pagination.Param
	OrderFields      pagination.OrderFields
	IDs              []string
	Name             string
	PrefixParentPath string
	QueryValue       string
	ParentID         *string
	ShowStatus       int
	Status           int
}
