package main

import (
	"fmt"

	"github.com/linzhengen/ddd-gin-admin/infrastructure/persistence/mysql"

	"github.com/linzhengen/ddd-gin-admin/configs"

	"github.com/jinzhu/gorm"
)

// InitStore init store.
func InitStore() (*gorm.DB, func(), error) {
	var storeCall func()
	cfg := configs.Env()

	db, err := initGorm()
	if err != nil {
		return nil, nil, err
	}

	storeCall = func() {
		db.Close()
	}

	mysql.SetTablePrefix(cfg.GormTablePrefix)

	if cfg.GormEnableAutoMigrate {
		err = mysql.AutoMigrate(db)
		if err != nil {
			return nil, nil, err
		}
	}

	return db, storeCall, nil
}

// initGorm init gorm.
func initGorm() (*gorm.DB, error) {
	cfg := configs.Env()

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?%s",
		cfg.MysqlUser, cfg.MysqlPassword, cfg.MysqlHost, cfg.MysqlPort, cfg.MysqlDbName, cfg.MysqlParameters)
	return mysql.NewDB(&mysql.DbConfig{
		Debug:        cfg.GormDebug,
		DBType:       cfg.GormDbType,
		DSN:          dsn,
		MaxIdleConns: cfg.GormMaxIdleConns,
		MaxLifetime:  cfg.GormMaxLifetime,
		MaxOpenConns: cfg.GormMaxOpenConns,
	})
}
