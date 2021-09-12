package application

import (
	"context"

	"github.com/linzhengen/ddd-gin-admin/configs"

	"github.com/casbin/casbin/v2"
	"github.com/linzhengen/ddd-gin-admin/pkg/logger"
)

var chCasbinPolicy chan *chCasbinPolicyItem

type chCasbinPolicyItem struct {
	ctx context.Context
	e   *casbin.SyncedEnforcer
}

//nolint:gochecknoinits
func init() {
	chCasbinPolicy = make(chan *chCasbinPolicyItem, 1)
	go func() {
		for item := range chCasbinPolicy {
			err := item.e.LoadPolicy()
			if err != nil {
				logger.WithContext(item.ctx).Errorf("The load casbin policy error: %s", err.Error())
			}
		}
	}()
}

func LoadCasbinPolicy(ctx context.Context, e *casbin.SyncedEnforcer) {
	if !configs.C.Casbin.Enable {
		return
	}

	if len(chCasbinPolicy) > 0 {
		logger.WithContext(ctx).Infof("The load casbin policy is already in the wait queue")
		return
	}

	chCasbinPolicy <- &chCasbinPolicyItem{
		ctx: ctx,
		e:   e,
	}
}
