package mysql

import (
	"time"

	"github.com/jinzhu/gorm"
	"github.com/linzhengen/ddd-gin-admin/domain/entity"
	"github.com/linzhengen/ddd-gin-admin/domain/repository"
)

// DbConfig ...
type DbConfig struct {
	Debug        bool
	DBType       string
	DSN          string
	MaxLifetime  time.Duration
	MaxOpenConns int
	MaxIdleConns int
}

// Repositories ...
type Repositories struct {
	User repository.UserRepository
	db   *gorm.DB
}

// NewRepositories ...
func NewRepositories(db *gorm.DB) (*Repositories, error) {
	return &Repositories{
		User: NewUserRepository(db),
		db:   db,
	}, nil
}

// closes the database connection.
func (s *Repositories) Close() error {
	return s.db.Close()
}

// AutoMigrate migrate all tables.
func AutoMigrate(db *gorm.DB) error {
	db = db.Set("gorm:table_options", "ENGINE=InnoDB")
	return db.AutoMigrate(&entity.User{}).Error
}

// NewDB ...
func NewDB(c *DbConfig) (*gorm.DB, error) {
	db, err := gorm.Open(c.DBType, c.DSN)
	if err != nil {
		return nil, err
	}

	if c.Debug {
		db = db.Debug()
	}

	err = db.DB().Ping()
	if err != nil {
		return nil, err
	}

	db.DB().SetMaxIdleConns(c.MaxIdleConns)
	db.DB().SetMaxOpenConns(c.MaxOpenConns)
	db.DB().SetConnMaxLifetime(c.MaxLifetime)
	return db, nil
}

// SetTablePrefix ...
func SetTablePrefix(prefix string) {
	entity.SetTablePrefix(prefix)
}
