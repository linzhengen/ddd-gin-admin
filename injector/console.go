package injector

import (
	"context"

	"github.com/linzhengen/ddd-gin-admin/app/interfaces/console/command"

	"github.com/linzhengen/ddd-gin-admin/configs"
	"github.com/linzhengen/ddd-gin-admin/injector/console"
)

func InitConsole(ctx context.Context, opts ...console.Option) (command.Commands, func(), error) {
	var o console.Options
	for _, opt := range opts {
		opt(&o)
	}

	configs.MustLoad(o.ConfigFile)
	configs.PrintWithJSON()

	loggerCleanFunc, err := console.InitLogger()
	if err != nil {
		return nil, nil, err
	}

	commands, consoleCleanFunc, err := BuildConsoleInjector(ctx)
	if err != nil {
		return nil, nil, err
	}

	return commands, func() {
		loggerCleanFunc()
		consoleCleanFunc()
	}, nil
}
