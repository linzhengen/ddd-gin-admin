package request

import (
	"time"

	"github.com/linzhengen/ddd-gin-admin/pkg/util/json"
)

type Menu struct {
	ID         string      `json:"id"`                                         // ID
	Name       string      `json:"name" binding:"required"`                    // Name
	Sequence   int         `json:"sequence"`                                   // Sequence
	Icon       string      `json:"icon"`                                       // Icon
	Router     string      `json:"router"`                                     // Router
	ParentID   string      `json:"parent_id"`                                  // Parent ID
	ParentPath string      `json:"parent_path"`                                // Parent Path
	ShowStatus int         `json:"show_status" binding:"required,max=2,min=1"` // Show Status(1:show 2:hide)
	Status     int         `json:"status" binding:"required,max=2,min=1"`      // Menu Status(1:enable 2:disable)
	Memo       string      `json:"memo"`                                       // Memo
	Creator    string      `json:"creator"`                                    // Creator
	CreatedAt  time.Time   `json:"created_at"`                                 // CreatedAt
	UpdatedAt  time.Time   `json:"updated_at"`                                 // UpdatedAt
	Actions    MenuActions `json:"actions"`                                    // Actions
}

func (a *Menu) String() string {
	return json.MarshalToString(a)
}

type MenuActions []*MenuAction

func (a MenuActions) ToMap() map[string]*MenuAction {
	m := make(map[string]*MenuAction)
	for _, item := range a {
		m[item.Code] = item
	}
	return m
}

func (a MenuActions) FillResources(mResources map[string]MenuActionResources) {
	for i, item := range a {
		a[i].Resources = mResources[item.ID]
	}
}

func (a MenuActions) ToMenuIDMap() map[string]MenuActions {
	m := make(map[string]MenuActions)
	for _, item := range a {
		m[item.MenuID] = append(m[item.MenuID], item)
	}
	return m
}

type MenuQuery struct {
	Pagination
	OrderFields
	IDs              []string `form:"-"`
	Name             string   `form:"-"`
	PrefixParentPath string   `form:"-"`
	QueryValue       string   `form:"queryValue"`
	ParentID         *string  `form:"parentID"`
	ShowStatus       int      `form:"showStatus"` // 1:show 2:hide
	Status           int      `form:"status"`     // 1:enable 2:disable
}

type MenuActionQuery struct {
	Pagination
	OrderFields
	MenuID string
	IDs    []string
}

type MenuActionResourceQuery struct {
	Pagination
	OrderFields
	MenuID  string
	MenuIDs []string
}

type MenuAction struct {
	ID        string              `yaml:"-" json:"id"`
	MenuID    string              `yaml:"-" binding:"required" json:"menu_id"`
	Code      string              `yaml:"code" binding:"required" json:"code"`
	Name      string              `yaml:"name" binding:"required" json:"name"`
	Resources MenuActionResources `yaml:"resources,omitempty" json:"resources"`
}

type MenuActionResources []*MenuActionResource

type MenuActionResource struct {
	ID       string
	ActionID string
	Method   string
	Path     string
}
