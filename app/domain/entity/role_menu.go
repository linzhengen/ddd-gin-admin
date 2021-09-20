package entity

type RoleMenu struct {
	ID       string `gorm:"column:id;primary_key;size:36;"`
	RoleID   string `gorm:"column:role_id;size:36;index;default:'';not null;"`
	MenuID   string `gorm:"column:menu_id;size:36;index;default:'';not null;"`
	ActionID string `gorm:"column:action_id;size:36;index;default:'';not null;"`
}

type RoleMenus []*RoleMenu
