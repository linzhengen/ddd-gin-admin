package entity

import (
	"context"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/linzhengen/ddd-gin-admin/domain/schema"
	"github.com/linzhengen/ddd-gin-admin/pkg/util/structure"
)

func GetDemoDB(ctx context.Context, defDB *gorm.DB) *gorm.DB {
	return GetDBWithModel(ctx, defDB, new(Demo))
}

type SchemaDemo schema.Demo

func (a SchemaDemo) ToDemo() *Demo {
	item := new(Demo)
	structure.Copy(a, item)
	return item
}

type Demo struct {
	ID        string     `gorm:"column:id;primary_key;size:36;"`
	Code      string     `gorm:"column:code;size:50;index;default:'';not null;"`
	Name      string     `gorm:"column:name;size:100;index;default:'';not null;"`
	Memo      *string    `gorm:"column:memo;size:200;"`
	Status    int        `gorm:"column:status;index;default:0;not null;"`
	Creator   string     `gorm:"column:creator;size:36;"`
	CreatedAt time.Time  `gorm:"column:created_at;index;"`
	UpdatedAt time.Time  `gorm:"column:updated_at;index;"`
	DeletedAt *time.Time `gorm:"column:deleted_at;index;"`
}

func (a Demo) ToSchemaDemo() *schema.Demo {
	item := new(schema.Demo)
	structure.Copy(a, item)
	return item
}

type Demos []*Demo

func (a Demos) ToSchemaDemos() []*schema.Demo {
	list := make([]*schema.Demo, len(a))
	for i, item := range a {
		list[i] = item.ToSchemaDemo()
	}
	return list
}
