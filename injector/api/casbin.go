package api

import (
	"time"

	"github.com/linzhengen/ddd-gin-admin/app/application"

	"github.com/linzhengen/ddd-gin-admin/configs"

	"github.com/casbin/casbin/v2"
)

func InitCasbin(adapter application.RbacAdapter) (*casbin.SyncedEnforcer, func(), error) {
	adapter.CreateAutoLoadPolicyChan()
	cfg := configs.C.Casbin
	if cfg.Model == "" {
		return new(casbin.SyncedEnforcer), nil, nil
	}

	e, err := casbin.NewSyncedEnforcer(cfg.Model)
	if err != nil {
		return nil, nil, err
	}
	e.EnableLog(cfg.Debug)

	err = e.InitWithModelAndAdapter(e.GetModel(), adapter)
	if err != nil {
		return nil, nil, err
	}
	e.EnableEnforce(cfg.Enable)

	cleanFunc := func() {}
	if cfg.AutoLoad {
		e.StartAutoLoadPolicy(time.Duration(cfg.AutoLoadInternal) * time.Second)
		cleanFunc = func() {
			e.StopAutoLoadPolicy()
		}
	}

	return e, cleanFunc, nil
}
