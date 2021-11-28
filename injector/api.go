package injector

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/linzhengen/ddd-gin-admin/app/application"

	"github.com/casbin/casbin/v2"

	"github.com/linzhengen/ddd-gin-admin/app/domain/auth"

	"github.com/linzhengen/ddd-gin-admin/configs"

	"github.com/gin-gonic/gin"
	"github.com/linzhengen/ddd-gin-admin/injector/api"

	"github.com/linzhengen/ddd-gin-admin/pkg/logger"

	_ "github.com/linzhengen/ddd-gin-admin/app/interfaces/api/swagger"
)

func NewApiInjector(
	engine *gin.Engine,
	auth auth.Repository,
	casbinEnforcer *casbin.SyncedEnforcer,
	seed application.Seed,
) *ApiInjector {
	return &ApiInjector{
		engine:         engine,
		auth:           auth,
		casbinEnforcer: casbinEnforcer,
		seed:           seed,
	}
}

type ApiInjector struct {
	engine         *gin.Engine
	auth           auth.Repository
	casbinEnforcer *casbin.SyncedEnforcer
	seed           application.Seed
}

func initHttpServer(ctx context.Context, opts ...api.Option) (func(), error) {
	var o api.Options
	for _, opt := range opts {
		opt(&o)
	}

	configs.MustLoad(o.ConfigFile)
	if v := o.ModelFile; v != "" {
		configs.C.Casbin.Model = v
	}
	if v := o.WWWDir; v != "" {
		configs.C.WWW = v
	}
	if v := o.MenuFile; v != "" {
		configs.C.Menu.Data = v
	}
	configs.PrintWithJSON()

	logger.WithContext(ctx).Printf("starting server，run mode：%s，ver：%s，pid：%d", configs.C.RunMode, o.Version, os.Getpid())

	loggerCleanFunc, err := InitLogger()
	if err != nil {
		return nil, err
	}

	monitorCleanFunc := api.InitMonitor(ctx)

	api.InitCaptcha()

	injector, injectorCleanFunc, err := BuildApiInjector()
	if err != nil {
		return nil, err
	}
	if configs.C.Menu.Enable && configs.C.Menu.Data != "" {
		err = injector.seed.Execute(ctx, configs.C.Menu.Data)
		if err != nil {
			return nil, err
		}
	}
	httpServerCleanFunc := api.InitHTTPServer(ctx, injector.engine)

	return func() {
		httpServerCleanFunc()
		injectorCleanFunc()
		monitorCleanFunc()
		loggerCleanFunc()
	}, nil
}

func RunServer(ctx context.Context, opts ...api.Option) error {
	state := 1
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	cleanFunc, err := initHttpServer(ctx, opts...)
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
	logger.WithContext(ctx).Infof("stopping server")
	time.Sleep(time.Second)
	os.Exit(state)
	return nil
}
