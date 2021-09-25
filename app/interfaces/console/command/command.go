package command

import (
	"context"

	"github.com/urfave/cli/v2"
)

type Commands []*cli.Command

func NewCliCommands(
	ctx context.Context,
	migrate MigrateCommand,
) Commands {
	return []*cli.Command{
		{
			Name:  "db:migrate",
			Usage: "Migrate database",
			Action: func(c *cli.Context) error {
				return migrate.Migrate(ctx, c)
			},
		},
	}
}
