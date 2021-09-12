package entity

import (
	"context"

	"github.com/casbin/casbin/v2"
)

type AutoLoadPolicyChan chan *CasbinPolicyItem

type CasbinPolicyItem struct {
	Ctx      context.Context
	Enforcer *casbin.SyncedEnforcer
}
