package injector

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

	menuInfra "github.com/linzhengen/ddd-gin-admin/app/infrastructure/menu"
	menuActionInfra "github.com/linzhengen/ddd-gin-admin/app/infrastructure/menu/menuaction"
	menuActionResourceInfra "github.com/linzhengen/ddd-gin-admin/app/infrastructure/menu/menuactionresource"
	userInfra "github.com/linzhengen/ddd-gin-admin/app/infrastructure/user"
	roleInfra "github.com/linzhengen/ddd-gin-admin/app/infrastructure/user/role"
	roleMenuInfra "github.com/linzhengen/ddd-gin-admin/app/infrastructure/user/rolemenu"
	userRoleInfra "github.com/linzhengen/ddd-gin-admin/app/infrastructure/user/userrole"

	"github.com/linzhengen/ddd-gin-admin/app/infrastructure/gormx"

	"github.com/linzhengen/ddd-gin-admin/configs"

	"github.com/jinzhu/gorm"
)

func InitGormDB() (*gorm.DB, func(), error) {
	cfg := configs.C.Gorm
	db, cleanFunc, err := NewGormDB()
	if err != nil {
		return nil, cleanFunc, err
	}

	if cfg.EnableAutoMigrate {
		err = autoMigrate(db)
		if err != nil {
			return nil, cleanFunc, err
		}
	}

	return db, cleanFunc, nil
}

func NewGormDB() (*gorm.DB, func(), error) {
	cfg := configs.C
	var dsn string
	switch cfg.Gorm.DBType {
	case "mysql":
		dsn = cfg.MySQL.DSN()
	case "sqlite3":
		dsn = cfg.Sqlite3.DSN()
		_ = os.MkdirAll(filepath.Dir(dsn), 0777)
	case "postgres":
		dsn = cfg.Postgres.DSN()
	default:
		return nil, nil, errors.New("unknown db")
	}

	return gormx.NewDB(&gormx.Config{
		Debug:        cfg.Gorm.Debug,
		DBType:       cfg.Gorm.DBType,
		DSN:          dsn,
		MaxIdleConns: cfg.Gorm.MaxIdleConns,
		MaxLifetime:  cfg.Gorm.MaxLifetime,
		MaxOpenConns: cfg.Gorm.MaxOpenConns,
		TablePrefix:  cfg.Gorm.TablePrefix,
	})
}

func autoMigrate(db *gorm.DB) error {
	if dbType := configs.C.Gorm.DBType; strings.ToLower(dbType) == "mysql" {
		db = db.Set("gorm:table_options", "ENGINE=InnoDB")
	}

	return db.AutoMigrate(
		new(menuActionInfra.Model),
		new(menuActionResourceInfra.Model),
		new(menuInfra.Model),
		new(roleInfra.Model),
		new(roleMenuInfra.Model),
		new(userRoleInfra.Model),
		new(userInfra.Model),
	).Error
}
