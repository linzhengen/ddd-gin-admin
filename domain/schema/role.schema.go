package schema

import "time"

type Role struct {
	ID        string    `json:"id"`                                    // ID
	Name      string    `json:"name" binding:"required"`               // Name
	Sequence  int       `json:"sequence"`                              // Sequence
	Memo      string    `json:"memo"`                                  // Memo
	Status    int       `json:"status" binding:"required,max=2,min=1"` // Status(1:enable 2:disable)
	Creator   string    `json:"creator"`                               // Creator
	CreatedAt time.Time `json:"created_at"`                            // CreatedAt
	UpdatedAt time.Time `json:"updated_at"`                            // UpdatedAt
	RoleMenus RoleMenus `json:"role_menus" binding:"required,gt=0"`    // RoleMenus
}

type RoleQueryParam struct {
	PaginationParam
	IDs        []string `form:"-"`          // IDs
	Name       string   `form:"-"`          // Name
	QueryValue string   `form:"queryValue"` // Query Search Values
	UserID     string   `form:"-"`          // User ID
	Status     int      `form:"status"`     // Status(1:enable 2:disable)
}

type RoleQueryOptions struct {
	OrderFields []*OrderField // Order Fields
}

type RoleQueryResult struct {
	Data       Roles
	PageResult *PaginationResult
}

type Roles []*Role

func (a Roles) ToNames() []string {
	names := make([]string, len(a))
	for i, item := range a {
		names[i] = item.Name
	}
	return names
}

func (a Roles) ToMap() map[string]*Role {
	m := make(map[string]*Role)
	for _, item := range a {
		m[item.ID] = item
	}
	return m
}

// ----------------------------------------RoleMenu--------------------------------------

type RoleMenu struct {
	ID       string `json:"id"`                           // ID
	RoleID   string `json:"role_id" binding:"required"`   // Role ID
	MenuID   string `json:"menu_id" binding:"required"`   // Menu ID
	ActionID string `json:"action_id" binding:"required"` // Action ID
}

type RoleMenuQueryParam struct {
	PaginationParam
	RoleID  string
	RoleIDs []string
}

type RoleMenuQueryOptions struct {
	OrderFields []*OrderField
}

type RoleMenuQueryResult struct {
	Data       RoleMenus
	PageResult *PaginationResult
}

type RoleMenus []*RoleMenu

func (a RoleMenus) ToMap() map[string]*RoleMenu {
	m := make(map[string]*RoleMenu)
	for _, item := range a {
		m[item.MenuID+"-"+item.ActionID] = item
	}
	return m
}

func (a RoleMenus) ToRoleIDMap() map[string]RoleMenus {
	m := make(map[string]RoleMenus)
	for _, item := range a {
		m[item.RoleID] = append(m[item.RoleID], item)
	}
	return m
}

func (a RoleMenus) ToMenuIDs() []string {
	var idList []string
	m := make(map[string]struct{})

	for _, item := range a {
		if _, ok := m[item.MenuID]; ok {
			continue
		}
		idList = append(idList, item.MenuID)
		m[item.MenuID] = struct{}{}
	}

	return idList
}

func (a RoleMenus) ToActionIDs() []string {
	idList := make([]string, len(a))
	m := make(map[string]struct{})
	for i, item := range a {
		if _, ok := m[item.ActionID]; ok {
			continue
		}
		idList[i] = item.ActionID
		m[item.ActionID] = struct{}{}
	}
	return idList
}
