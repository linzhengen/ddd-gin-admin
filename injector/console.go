package injector

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/linzhengen/ddd-gin-admin/configs"
	"github.com/linzhengen/ddd-gin-admin/injector/console"

	"github.com/linzhengen/ddd-gin-admin/pkg/logger"
)

type ConsoleInjector struct {
}

func initConsole(ctx context.Context, opts ...console.Option) (func(), error) {
	var o console.Options
	for _, opt := range opts {
		opt(&o)
	}

	configs.MustLoad(o.ConfigFile)
	configs.PrintWithJSON()

	loggerCleanFunc, err := console.InitLogger()
	if err != nil {
		return nil, err
	}

	//injector, err := BuildConsoleInjector()
	//if err != nil {
	//	return nil, err
	//}

	return func() {
		loggerCleanFunc()
	}, nil
}

func RunConsole(ctx context.Context, opts ...console.Option) error {
	state := 1
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	cleanFunc, err := initConsole(ctx, opts...)
	if err != nil {
		return err
	}

EXIT:
	for {
		sig := <-sc
		logger.WithContext(ctx).Infof("catched signal[%s]", sig.String())
		switch sig {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			state = 0
			break EXIT
		case syscall.SIGHUP:
		default:
			break EXIT
		}
	}

	cleanFunc()
	logger.WithContext(ctx).Infof("stopping console")
	time.Sleep(time.Second)
	os.Exit(state)
	return nil
}
