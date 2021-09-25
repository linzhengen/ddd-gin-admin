package command

import (
	"context"

	"github.com/linzhengen/ddd-gin-admin/pkg/logger"

	"github.com/linzhengen/ddd-gin-admin/app/application"

	"github.com/urfave/cli/v2"
)

type MigrateCommand interface {
	Migrate(ctx context.Context, c *cli.Context) error
}

func NewMigrateCommand(migrateApp application.DbMigrationConsole) MigrateCommand {
	return &migrateCommand{
		migrateApp: migrateApp,
	}
}

type migrateCommand struct {
	migrateApp application.DbMigrationConsole
}

func (m migrateCommand) Migrate(ctx context.Context, c *cli.Context) error {
	logger.WithContext(ctx).Info("DB migration starting...")
	if err := m.migrateApp.Migrate(ctx); err != nil {
		logger.WithContext(ctx).Error(err)
		return err
	}
	logger.WithContext(ctx).Info("DB migration done")
	return nil
}
