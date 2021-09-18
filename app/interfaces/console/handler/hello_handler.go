package handler

import (
	"context"

	"github.com/linzhengen/ddd-gin-admin/pkg/logger"

	"github.com/linzhengen/ddd-gin-admin/app/application"

	"github.com/urfave/cli/v2"
)

type HelloHandler interface {
	Hello(ctx context.Context, c *cli.Context)
}

func NewHelloHandler(helloApp application.HelloConsole) HelloHandler {
	return &helloHandler{
		helloApp: helloApp,
	}
}

type helloHandler struct {
	helloApp application.HelloConsole
}

func (a *helloHandler) Hello(ctx context.Context, c *cli.Context) {
	name, err := a.helloApp.GetUserName(ctx, c.String("id"))
	if err != nil {
		logger.Errorf("error %v", err)
		return
	}
	logger.Infof("hello %s", name)
}
