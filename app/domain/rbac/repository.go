package rbac

import (
	"context"
)

type Repository interface {
	CreateAutoLoadPolicyChan() AutoLoadPolicyChan
	GetAutoLoadPolicyChan() AutoLoadPolicyChan
	AddPolicyItemToChan(ctx context.Context, e SyncedEnforcer)
}
