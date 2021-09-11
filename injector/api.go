package injector

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/linzhengen/ddd-gin-admin/app/application"
	"github.com/linzhengen/ddd-gin-admin/app/infrastructure/config"

	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"github.com/linzhengen/ddd-gin-admin/pkg/auth"

	"github.com/linzhengen/ddd-gin-admin/injector/api"

	"github.com/linzhengen/ddd-gin-admin/pkg/logger"

	_ "github.com/linzhengen/ddd-gin-admin/app/interfaces/api/swagger"
)

func NewApiInjector(
	engine *gin.Engine,
	auth auth.Author,
	casbinEnforcer *casbin.SyncedEnforcer,
	menuBll application.Menu,
) *ApiInjector {
	return &ApiInjector{
		engine:         engine,
		auth:           auth,
		casbinEnforcer: casbinEnforcer,
		menuBll:        menuBll,
	}
}

type ApiInjector struct {
	engine         *gin.Engine
	auth           auth.Author
	casbinEnforcer *casbin.SyncedEnforcer
	menuBll        application.Menu
}

func initHttpServer(ctx context.Context, opts ...api.Option) (func(), error) {
	var o api.Options
	for _, opt := range opts {
		opt(&o)
	}

	config.MustLoad(o.ConfigFile)
	if v := o.ModelFile; v != "" {
		config.C.Casbin.Model = v
	}
	if v := o.WWWDir; v != "" {
		config.C.WWW = v
	}
	if v := o.MenuFile; v != "" {
		config.C.Menu.Data = v
	}
	config.PrintWithJSON()

	logger.WithContext(ctx).Printf("starting server，run mode：%s，ver：%s，pid：%d", config.C.RunMode, o.Version, os.Getpid())

	loggerCleanFunc, err := api.InitLogger()
	if err != nil {
		return nil, err
	}

	monitorCleanFunc := api.InitMonitor(ctx)

	api.InitCaptcha()

	injector, injectorCleanFunc, err := BuildApiInjector()
	if err != nil {
		return nil, err
	}

	if config.C.Menu.Enable && config.C.Menu.Data != "" {
		err = injector.menuBll.InitData(ctx, config.C.Menu.Data)
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
