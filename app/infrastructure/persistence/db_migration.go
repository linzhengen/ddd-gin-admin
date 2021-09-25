package persistence

import (
	"context"
	"strings"

	"github.com/jinzhu/gorm"
	"github.com/linzhengen/ddd-gin-admin/app/domain/entity"
	"github.com/linzhengen/ddd-gin-admin/configs"

	"github.com/linzhengen/ddd-gin-admin/app/domain/repository"
)

func NewDbMigration(db *gorm.DB) repository.DBMigrationRepository {
	return &dbMigration{
		db: db,
	}
}

type dbMigration struct {
	db *gorm.DB
}

func (d *dbMigration) Migrate(ctx context.Context) error {
	if dbType := configs.C.Gorm.DBType; strings.ToLower(dbType) == "mysql" {
		d.db = d.db.Set("gorm:table_options", "ENGINE=InnoDB")
	}

	return d.db.AutoMigrate(
		new(entity.MenuAction),
		new(entity.MenuActionResource),
		new(entity.Menu),
		new(entity.RoleMenu),
		new(entity.Role),
		new(entity.UserRole),
		new(entity.User),
	).Error
}
