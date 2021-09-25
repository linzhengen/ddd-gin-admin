package command

import (
	"context"

	"github.com/linzhengen/ddd-gin-admin/pkg/logger"

	"github.com/linzhengen/ddd-gin-admin/app/application"

	"github.com/urfave/cli/v2"
)

type SeedCommand interface {
	Seed(ctx context.Context, c *cli.Context) error
}

func NewSeedCommand(seedApp application.DbSeedConsole) SeedCommand {
	return &seedCommand{
		seedApp: seedApp,
	}
}

type seedCommand struct {
	seedApp application.DbSeedConsole
}

func (s seedCommand) Seed(ctx context.Context, c *cli.Context) error {
	logger.WithContext(ctx).Info("DB seeding starting...")
	if err := s.seedApp.Seed(ctx); err != nil {
		logger.WithContext(ctx).Error(err)
		return err
	}
	logger.WithContext(ctx).Info("DB seeding done")
	return nil
}
