package api

import (
	"time"

	"github.com/casbin/casbin/v3"
	"github.com/casbin/casbin/v3/log"

	"github.com/linzhengen/ddd-gin-admin/app/application"
	"github.com/linzhengen/ddd-gin-admin/configs"
)

func InitCasbin(adapter application.RbacAdapter) (*casbin.SyncedEnforcer, func(), error) {
	adapter.CreateAutoLoadPolicyChan()
	cfg := configs.C.Casbin
	if cfg.Model == "" {
		return new(casbin.SyncedEnforcer), func() {}, nil
	}

	e, err := casbin.NewSyncedEnforcer(cfg.Model)
	if err != nil {
		return nil, nil, err
	}
	if cfg.Debug {
		l := log.NewDefaultLogger()
		_ = l.SetEventTypes([]log.EventType{
			log.EventEnforce,
			log.EventAddPolicy,
			log.EventRemovePolicy,
			log.EventLoadPolicy,
			log.EventSavePolicy,
		})
		e.SetLogger(l)
	}

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
