package command

import (
	"context"

	"github.com/linzhengen/ddd-gin-admin/pkg/logger"

	"github.com/linzhengen/ddd-gin-admin/app/application"

	"github.com/urfave/cli/v2"
)

type HelloCommand interface {
	Hello(ctx context.Context, c *cli.Context) error
}

func NewHelloCommand(helloApp application.HelloConsole) HelloCommand {
	return &helloCommand{
		helloApp: helloApp,
	}
}

type helloCommand struct {
	helloApp application.HelloConsole
}

func (a *helloCommand) Hello(ctx context.Context, c *cli.Context) error {
	name, err := a.helloApp.GetUserName(ctx, c.String("id"))
	if err != nil {
		return err
	}
	logger.Infof("hello %s", name)
	return nil
}
