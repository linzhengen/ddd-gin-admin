package gormx

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gorm.io/gorm/schema"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"

	"github.com/linzhengen/ddd-gin-admin/pkg/logger"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"

	_ "gorm.io/driver/mysql"
	_ "gorm.io/driver/postgres"
	_ "gorm.io/driver/sqlite"
)

type Config struct {
	Debug        bool
	DBType       string
	DSN          string
	MaxLifetime  int
	MaxOpenConns int
	MaxIdleConns int
	TablePrefix  string
}

func NewDB(c *Config) (*gorm.DB, func(), error) {
	var dialector gorm.Dialector

	switch strings.ToLower(c.DBType) {
	case "mysql":
		dialector = mysql.Open(c.DSN)
	case "postgres":
		dialector = postgres.Open(c.DSN)
	case "sqlite3":
		_ = os.MkdirAll(filepath.Dir(c.DSN), os.ModePerm)
		dialector = sqlite.Open(c.DSN)
	default:
		panic(fmt.Sprintf("unsupported database type: %s", c.DBType))
	}

	ormCfg := &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   c.TablePrefix,
			SingularTable: true,
		},
		Logger: gormLogger.Discard,
	}

	if c.Debug {
		ormCfg.Logger = gormLogger.Default
	}

	db, err := gorm.Open(dialector, ormCfg)
	if err != nil {
		return nil, nil, err
	}

	if c.Debug {
		db = db.Debug()
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, nil, err
	}
	cleanFunc := func() {
		err := sqlDB.Close()
		if err != nil {
			logger.Errorf("Gorm db close error: %s", err.Error())
		}
	}

	err = sqlDB.Ping()
	if err != nil {
		return nil, cleanFunc, err
	}

	sqlDB.SetMaxIdleConns(c.MaxIdleConns)
	sqlDB.SetMaxOpenConns(c.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Duration(c.MaxLifetime) * time.Second)
	return db, cleanFunc, nil
}
