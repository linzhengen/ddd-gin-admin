/*
Package main ddd-gin-admin

Swagger：https://github.com/swaggo/swag#declarative-comments-format

Usage：

	go get -u github.com/swaggo/swag/main/swag
	swag init --generalInfo ./main/main.go --output ./app/interfaces/api/swagger */
package main

import (
	"context"
	"os"

	"github.com/linzhengen/ddd-gin-admin/injector/console"

	"github.com/linzhengen/ddd-gin-admin/injector/api"

	"github.com/linzhengen/ddd-gin-admin/injector"

	"github.com/linzhengen/ddd-gin-admin/pkg/logger"
	"github.com/urfave/cli/v2"
)

// VERSION You can specify the version number by compiling：go build -ldflags "-X main.VERSION=x.x.x"
var VERSION = "0.4.0"

// @title ddd-gin-admin
// @version 0.2.0
// @description RBAC scaffolding based on DDD + GIN + GORM + CASBIN + WIRE.
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
// @schemes http https
// @basePath /
func main() {
	logger.SetVersion(VERSION)
	ctx := logger.NewTagContext(context.Background(), "__main__")
	app := cli.NewApp()
	app.Name = "ddd-gin-admin"
	app.Version = VERSION
	app.Usage = "RBAC scaffolding based on DDD + GIN + GORM + CASBIN + WIRE."
	app.Commands = []*cli.Command{
		newWebCmd(ctx),
		newConsoleCmd(ctx),
	}
	err := app.Run(os.Args)
	if err != nil {
		logger.WithContext(ctx).Errorf(err.Error())
	}
}

func newWebCmd(ctx context.Context) *cli.Command {
	return &cli.Command{
		Name:  "web",
		Usage: "Run web server",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "conf",
				Aliases:  []string{"c"},
				Usage:    "server config files(.json,.yaml,.toml)",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "model",
				Aliases:  []string{"m"},
				Usage:    "casbin model config(.conf)",
				Required: true,
			},
			&cli.StringFlag{
				Name:  "menu",
				Usage: "default menus config(.yaml)",
			},
			&cli.StringFlag{
				Name:  "www",
				Usage: "static file dir",
			},
		},
		Action: func(c *cli.Context) error {
			return injector.RunServer(ctx,
				api.SetConfigFile(c.String("conf")),
				api.SetModelFile(c.String("model")),
				api.SetWWWDir(c.String("www")),
				api.SetMenuFile(c.String("menu")),
				api.SetVersion(VERSION))
		},
	}
}

func newConsoleCmd(ctx context.Context) *cli.Command {
	commands, cleanFunc, err := injector.InitConsole(ctx,
		console.SetConfigFile("./configs/config.toml"),
		console.SetVersion(VERSION))
	return &cli.Command{
		Name:  "console",
		Usage: "Run console commands",
		Before: func(c *cli.Context) error {
			return err
		},
		After: func(c *cli.Context) error {
			cleanFunc()
			return nil
		},
		Action: func(c *cli.Context) error {
			return nil
		},
		Subcommands: commands,
	}
}
