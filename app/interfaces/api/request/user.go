package request

import (
	"time"

	"github.com/linzhengen/ddd-gin-admin/app/domain/user"
	"github.com/linzhengen/ddd-gin-admin/pkg/util/structure"

	"github.com/linzhengen/ddd-gin-admin/app/domain/pagination"
)

type User struct {
	ID        string    `json:"id"`         // ID
	UserName  string    `json:"user_name"`  // User Name
	RealName  string    `json:"real_name"`  // Real Name
	Password  string    `json:"password"`   // Password
	Phone     string    `json:"phone"`      // Phone
	Email     string    `json:"email"`      // Email
	Status    int       `json:"status"`     // Status(1:enable 2:disable)
	Creator   string    `json:"creator"`    // Creator
	CreatedAt time.Time `json:"created_at"` // CreatedAt
	RoleIDs   []string  `json:"role_ids"`   // RoleIDs
}

func (a *User) ToDomain() *user.User {
	item := new(user.User)
	structure.Copy(a, item)
	return item
}

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
