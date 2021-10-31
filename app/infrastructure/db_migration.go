package infrastructure

import (
	"context"
	"strings"

	"github.com/linzhengen/ddd-gin-admin/app/infrastructure/menuactionresource"

	"github.com/linzhengen/ddd-gin-admin/app/infrastructure/menuaction"

	"github.com/linzhengen/ddd-gin-admin/app/infrastructure/menu"

	"github.com/linzhengen/ddd-gin-admin/app/infrastructure/role"

	"github.com/linzhengen/ddd-gin-admin/app/infrastructure/userrole"

	"github.com/linzhengen/ddd-gin-admin/app/infrastructure/user"

	"github.com/jinzhu/gorm"
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
		new(user.Model),
		new(userrole.Model),
		new(role.Model),
		new(menu.Model),
		new(menuaction.Model),
		new(menuactionresource.Model),
	).Error
}
