package command

import (
	"context"

	"github.com/urfave/cli/v2"
)

type Commands []*cli.Command

func NewCliCommands(
	ctx context.Context,
	hello HelloCommand,
) Commands {
	return []*cli.Command{
		{
			Name:  "hello",
			Usage: "echo hello",
			Action: func(c *cli.Context) error {
				return hello.Hello(ctx, c)
			},
		},
	}
}
