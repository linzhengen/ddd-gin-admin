/*
Package main ddd-gin-admin

Swagger：https://github.com/swaggo/swag#declarative-comments-format

Usage：

	go get -u github.com/swaggo/swag/main/swag
	swag init --generalInfo ./main/main.go --output ./interfaces/swagger */
package main

import (
	"context"
	"os"

	"github.com/linzhengen/ddd-gin-admin/injector"
	"github.com/linzhengen/ddd-gin-admin/pkg/logger"
	"github.com/urfave/cli/v2"
)

// VERSION You can specify the version number by compiling：go build -ldflags "-X main.VERSION=x.x.x"
var VERSION = "0.2.0"

// @title ddd-gin-admin
// @version 0.2.0
// @description RBAC scaffolding based on GIN + GORM + CASBIN + WIRE.
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
	app.Usage = "RBAC scaffolding based on GIN + GORM + CASBIN + WIRE."
	app.Commands = []*cli.Command{
		newWebCmd(ctx),
	}
	err := app.Run(os.Args)
	if err != nil {
		logger.WithContext(ctx).Errorf(err.Error())
	}
}

func newWebCmd(ctx context.Context) *cli.Command {
	return &cli.Command{
		Name:  "web",
		Usage: "run web server",
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
			return injector.Run(ctx,
				injector.SetConfigFile(c.String("conf")),
				injector.SetModelFile(c.String("model")),
				injector.SetWWWDir(c.String("www")),
				injector.SetMenuFile(c.String("menu")),
				injector.SetVersion(VERSION))
		},
	}
}
