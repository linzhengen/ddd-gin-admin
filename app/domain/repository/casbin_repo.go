package repository

import (
	"context"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/persist"
	"github.com/linzhengen/ddd-gin-admin/app/domain/entity"
)

type CasbinAdapter interface {
	persist.Adapter
	CreateAutoLoadPolicyChan() entity.AutoLoadPolicyChan
	GetAutoLoadPolicyChan() entity.AutoLoadPolicyChan
	AddCasbinPolicyItemToChan(ctx context.Context, e *casbin.SyncedEnforcer)
}
