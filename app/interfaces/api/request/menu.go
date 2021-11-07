package request

import (
	"time"

	"github.com/linzhengen/ddd-gin-admin/app/domain/pagination"
)

type Menu struct {
	ID         string    `json:"id"`                                         // ID
	Name       string    `json:"name" binding:"required"`                    // Name
	Sequence   int       `json:"sequence"`                                   // Sequence
	Icon       string    `json:"icon"`                                       // Icon
	Router     string    `json:"router"`                                     // Router
	ParentID   string    `json:"parent_id"`                                  // Parent ID
	ParentPath string    `json:"parent_path"`                                // Parent Path
	ShowStatus int       `json:"show_status" binding:"required,max=2,min=1"` // Show Status(1:show 2:hide)
	Status     int       `json:"status" binding:"required,max=2,min=1"`      // Menu Status(1:enable 2:disable)
	Memo       string    `json:"memo"`                                       // Memo
	Creator    string    `json:"creator"`                                    // Creator
	CreatedAt  time.Time `json:"created_at"`                                 // CreatedAt
	UpdatedAt  time.Time `json:"updated_at"`                                 // UpdatedAt
}

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
