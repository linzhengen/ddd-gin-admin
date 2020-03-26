package entity

import (
	"context"
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
	icontext "github.com/linzhengen/ddd-gin-admin/infrastructure/context"
	"github.com/linzhengen/ddd-gin-admin/infrastructure/json"
)

var tablePrefix string

// SetTablePrefix is table prefix.
func SetTablePrefix(prefix string) {
	tablePrefix = prefix
}

// GetTablePrefix gets table prefix.
func GetTablePrefix() string {
	return tablePrefix
}

// Model base model.
type Model struct {
	ID        uint      `gorm:"column:id;primary_key;auto_increment;"`
	CreatedAt time.Time `gorm:"column:created_at;"`
	UpdatedAt time.Time `gorm:"column:updated_at;"`
	DeletedAt time.Time `gorm:"column:deleted_at;index;"`
}

// TableName table name.
func (Model) TableName(name string) string {
	return fmt.Sprintf("%s%s", GetTablePrefix(), name)
}

func toString(v interface{}) string {
	return json.MarshalToString(v)
}

func getDB(ctx context.Context, defDB *gorm.DB) *gorm.DB {
	trans, ok := icontext.FromTrans(ctx)
	if ok {
		db, ok := trans.(*gorm.DB)
		if ok {
			return db
		}
	}
	return defDB
}

func getDBWithModel(ctx context.Context, defDB *gorm.DB, m interface{}) *gorm.DB {
	return getDB(ctx, defDB).Model(m)
}
