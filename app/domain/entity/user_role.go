package entity

type UserRole struct {
	ID     string `gorm:"column:id;primary_key;size:36;"`
	UserID string `gorm:"column:user_id;size:36;index;default:'';not null;"`
	RoleID string `gorm:"column:role_id;size:36;index;default:'';not null;"`
}

type UserRoles []*UserRole
