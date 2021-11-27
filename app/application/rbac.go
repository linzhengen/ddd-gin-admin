package application

import (
	"context"

	"github.com/casbin/casbin/v2"

	"github.com/linzhengen/ddd-gin-admin/app/domain/rbac"

	casbinModel "github.com/casbin/casbin/v2/model"
	"github.com/linzhengen/ddd-gin-admin/configs"
	"github.com/linzhengen/ddd-gin-admin/pkg/logger"

	"github.com/casbin/casbin/v2/persist"
)

type AutoLoadPolicyChan chan *PolicyItem

type PolicyItem struct {
	Ctx      context.Context
	Enforcer *casbin.SyncedEnforcer
}

type RbacAdapter interface {
	persist.Adapter
	CreateAutoLoadPolicyChan() AutoLoadPolicyChan
	GetAutoLoadPolicyChan() AutoLoadPolicyChan
	AddPolicyItemToChan(ctx context.Context, e *casbin.SyncedEnforcer)
}

type rbacAdapter struct {
	rbacRepo rbac.Repository
}

func NewRbacAdapter(rbacRepo rbac.Repository) RbacAdapter {
	return &rbacAdapter{
		rbacRepo: rbacRepo,
	}
}

var autoLoadPolicyChan AutoLoadPolicyChan

func (a *rbacAdapter) CreateAutoLoadPolicyChan() AutoLoadPolicyChan {
	autoLoadPolicyChan = make(chan *PolicyItem, 1)
	go func() {
		for item := range autoLoadPolicyChan {
			err := item.Enforcer.LoadPolicy()
			if err != nil {
				logger.WithContext(item.Ctx).Errorf("The load casbin policy error: %s", err.Error())
			}
		}
	}()
	return autoLoadPolicyChan
}

func (a *rbacAdapter) GetAutoLoadPolicyChan() AutoLoadPolicyChan {
	return autoLoadPolicyChan
}

func (a *rbacAdapter) AddPolicyItemToChan(ctx context.Context, e *casbin.SyncedEnforcer) {
	if !configs.C.Casbin.Enable {
		return
	}

	if len(autoLoadPolicyChan) > 0 {
		logger.WithContext(ctx).Infof("The load casbin policy is already in the wait queue")
		return
	}

	autoLoadPolicyChan <- &PolicyItem{
		Ctx:      ctx,
		Enforcer: e,
	}
}

func (a *rbacAdapter) LoadPolicy(model casbinModel.Model) error {
	ctx := context.Background()
	policies, err := a.rbacRepo.ListRolesPolices(ctx)
	if err != nil {
		logger.WithContext(ctx).Errorf("Load casbin role policy error: %s", err.Error())
		return err
	}
	if len(policies) > 0 {
		persist.LoadPolicyArray(policies, model)
	}

	policies, err = a.rbacRepo.ListUsersPolices(ctx)
	if err != nil {
		logger.WithContext(ctx).Errorf("Load casbin user policy error: %s", err.Error())
		return err
	}
	if len(policies) > 0 {
		persist.LoadPolicyArray(policies, model)
	}
	return nil
}

// SavePolicy saves all policy rules to the storage.
func (a *rbacAdapter) SavePolicy(model casbinModel.Model) error {
	return nil
}

// AddPolicy adds a policy rule to the storage.
// This is part of the Auto-Save feature.
func (a *rbacAdapter) AddPolicy(sec string, ptype string, rule []string) error {
	return nil
}

// RemovePolicy removes a policy rule from the storage.
// This is part of the Auto-Save feature.
func (a *rbacAdapter) RemovePolicy(sec string, ptype string, rule []string) error {
	return nil
}

// RemoveFilteredPolicy removes policy rules that match the filter from the storage.
// This is part of the Auto-Save feature.
func (a *rbacAdapter) RemoveFilteredPolicy(sec string, ptype string, fieldIndex int, fieldValues ...string) error {
	return nil
}
