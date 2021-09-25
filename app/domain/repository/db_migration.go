package repository

import "context"

type DBMigrationRepository interface {
	Migrate(ctx context.Context) error
}
