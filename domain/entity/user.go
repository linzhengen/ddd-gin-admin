package entity

import (
	"context"

	"github.com/jinzhu/gorm"
	"github.com/linzhengen/ddd-gin-admin/domain/schema"
)

// GetUserDB gets user db.
func GetUserDB(ctx context.Context, defDB *gorm.DB) *gorm.DB {
	return getDBWithModel(ctx, defDB, User{})
}

// User is struct.
type User struct {
	Model
	RecordID string `gorm:"column:record_id;size:36;index;"`
	UserName string `gorm:"column:user_name;size:64;index;"`
	RealName string `gorm:"column:real_name;size:64;index;"`
	Password string `gorm:"column:password;size:40;"`
	Email    string `gorm:"column:email;size:255;index;"`
	Phone    string `gorm:"column:phone;size:20;index;"`
	Status   int    `gorm:"column:status;index;"`
	Creator  string `gorm:"column:creator;size:36;"`
}

// TableName users.
func (u User) TableName() string {
	return u.Model.TableName("users")
}

// ToSchemaUser to schema user.
func (u User) ToSchemaUser() *schema.User {
	item := &schema.User{
		RecordID:  u.RecordID,
		UserName:  u.UserName,
		RealName:  u.RealName,
		Password:  u.Password,
		Status:    u.Status,
		Creator:   u.Creator,
		Email:     u.Email,
		Phone:     u.Phone,
		CreatedAt: u.CreatedAt,
	}
	return item
}

// Users list of User
type Users []*User

// ToSchemaUsers to schema users
func (us Users) ToSchemaUsers() []*schema.User {
	list := make([]*schema.User, len(us))
	for i, item := range us {
		list[i] = item.ToSchemaUser()
	}
	return list
}

// SchemaUser schema user.
type SchemaUser schema.User

// ToUser to user from schema data.
func (a SchemaUser) ToUser() *User {
	item := &User{
		RecordID: a.RecordID,
		UserName: a.UserName,
		RealName: a.RealName,
		Password: a.Password,
		Status:   a.Status,
		Creator:  a.Creator,
		Email:    a.Email,
		Phone:    a.Phone,
	}
	return item
}
