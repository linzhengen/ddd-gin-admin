package rbac

import (
	"context"
)

type AutoLoadPolicyChan chan *Policy

type SyncedEnforcer interface {
	LoadPolicy() error
}

type Policy struct {
	Ctx      context.Context
	Enforcer SyncedEnforcer
}
