package command

import (
	"context"

	"github.com/urfave/cli/v2"
)

type Commands []*cli.Command

func NewCliCommands(
	ctx context.Context,
	migrateCmd MigrateCommand,
	seedCmd SeedCommand,
) Commands {
	return []*cli.Command{
		{
			Name:  "db:migrate",
			Usage: "Migrate database",
			Action: func(c *cli.Context) error {
				return migrateCmd.Migrate(ctx, c)
			},
		},
		{
			Name:  "db:seed",
			Usage: "Running Seeders",
			Action: func(c *cli.Context) error {
				return seedCmd.Seed(ctx, c)
			},
		},
	}
}
