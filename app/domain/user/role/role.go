package role

import (
	"time"

	"github.com/linzhengen/ddd-gin-admin/app/domain/user/rolemenu"

	"github.com/linzhengen/ddd-gin-admin/app/domain/pagination"
)

type Role struct {
	ID        string
	Name      string
	Sequence  int
	Memo      *string
	Status    int
	Creator   string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
	RoleMenus rolemenu.RoleMenus
}

type Roles []*Role

type QueryParam struct {
	PaginationParam pagination.Param
	OrderFields     pagination.OrderFields
	IDs             []string
	Name            string
	QueryValue      string
	UserID          string
	Status          int
}
