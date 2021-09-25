package application

import (
	"context"

	"github.com/linzhengen/ddd-gin-admin/app/domain/repository"
)

type DbMigrationConsole interface {
	Migrate(ctx context.Context) error
}

func NewDbMigrationConsole(dbMigrationRepo repository.DBMigrationRepository) DbMigrationConsole {
	return &dbMigrationConsole{
		dbMigrationRepo: dbMigrationRepo,
	}
}

type dbMigrationConsole struct {
	dbMigrationRepo repository.DBMigrationRepository
}

func (a dbMigrationConsole) Migrate(ctx context.Context) error {
	return a.dbMigrationRepo.Migrate(ctx)
}
