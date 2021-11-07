package request

import (
	"time"

	"github.com/linzhengen/ddd-gin-admin/app/domain/user/role"
	"github.com/linzhengen/ddd-gin-admin/pkg/util/structure"

	"github.com/linzhengen/ddd-gin-admin/app/domain/pagination"
)

type Role struct {
	Name      string    `json:"name" binding:"required"`               // Name
	Sequence  int       `json:"sequence"`                              // Sequence
	Memo      string    `json:"memo"`                                  // Memo
	Status    int       `json:"status" binding:"required,max=2,min=1"` // Status(1:enable 2:disable)
	Creator   string    `json:"creator"`                               // Creator
	CreatedAt time.Time `json:"created_at"`                            // CreatedAt
	UpdatedAt time.Time `json:"updated_at"`                            // UpdatedAt
}

func (a *Role) ToDomain() *role.Role {
	item := new(role.Role)
	structure.Copy(a, item)
	return item
}

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
