package gormx

import (
	"time"

	"github.com/jinzhu/gorm"
	"github.com/linzhengen/ddd-gin-admin/pkg/logger"

	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
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
	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		return c.TablePrefix + defaultTableName
	}

	db, err := gorm.Open(c.DBType, c.DSN)
	if err != nil {
		return nil, nil, err
	}

	if c.Debug {
		db = db.Debug()
	}

	cleanFunc := func() {
		err := db.Close()
		if err != nil {
			logger.Errorf("Gorm db close error: %s", err.Error())
		}
	}

	err = db.DB().Ping()
	if err != nil {
		return nil, cleanFunc, err
	}

	db.SingularTable(true)
	db.DB().SetMaxIdleConns(c.MaxIdleConns)
	db.DB().SetMaxOpenConns(c.MaxOpenConns)
	db.DB().SetConnMaxLifetime(time.Duration(c.MaxLifetime) * time.Second)
	return db, cleanFunc, nil
}
