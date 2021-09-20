package casbin

import (
	"context"
	"fmt"

	"github.com/casbin/casbin/v2"

	"github.com/linzhengen/ddd-gin-admin/configs"

	"github.com/linzhengen/ddd-gin-admin/app/domain/entity"

	"github.com/linzhengen/ddd-gin-admin/app/domain/repository"
	"github.com/linzhengen/ddd-gin-admin/app/domain/valueobject/schema"

	casbinModel "github.com/casbin/casbin/v2/model"
	"github.com/casbin/casbin/v2/persist"
	"github.com/linzhengen/ddd-gin-admin/pkg/logger"
)

var autoLoadPolicyChan entity.AutoLoadPolicyChan

func NewCasbinAdapter(
	roleRepo repository.RoleRepository,
	roleMenuRepo repository.RoleMenuRepository,
	menuResourceRepo repository.MenuActionResourceRepository,
	userRepo repository.UserRepository,
	userRoleRepo repository.UserRoleRepository,
) repository.CasbinAdapter {
	return &casbinAdapter{
		roleRepo:         roleRepo,
		roleMenuRepo:     roleMenuRepo,
		menuResourceRepo: menuResourceRepo,
		userRepo:         userRepo,
		userRoleRepo:     userRoleRepo,
	}
}

type casbinAdapter struct {
	roleRepo         repository.RoleRepository
	roleMenuRepo     repository.RoleMenuRepository
	menuResourceRepo repository.MenuActionResourceRepository
	userRepo         repository.UserRepository
	userRoleRepo     repository.UserRoleRepository
}

func (a *casbinAdapter) AddCasbinPolicyItemToChan(ctx context.Context, e *casbin.SyncedEnforcer) {
	if !configs.C.Casbin.Enable {
		return
	}

	if len(autoLoadPolicyChan) > 0 {
		logger.WithContext(ctx).Infof("The load casbin policy is already in the wait queue")
		return
	}

	autoLoadPolicyChan <- &entity.CasbinPolicyItem{
		Ctx:      ctx,
		Enforcer: e,
	}
}

func (a *casbinAdapter) GetAutoLoadPolicyChan() entity.AutoLoadPolicyChan {
	return autoLoadPolicyChan
}

func (a *casbinAdapter) CreateAutoLoadPolicyChan() entity.AutoLoadPolicyChan {
	autoLoadPolicyChan = make(chan *entity.CasbinPolicyItem, 1)
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

func (a *casbinAdapter) LoadPolicy(model casbinModel.Model) error {
	ctx := context.Background()
	err := a.loadRolePolicy(ctx, model)
	if err != nil {
		logger.WithContext(ctx).Errorf("Load casbin role policy error: %s", err.Error())
		return err
	}

	err = a.loadUserPolicy(ctx, model)
	if err != nil {
		logger.WithContext(ctx).Errorf("Load casbin user policy error: %s", err.Error())
		return err
	}

	return nil
}

func (a *casbinAdapter) loadRolePolicy(ctx context.Context, m casbinModel.Model) error {
	roleResult, err := a.roleRepo.Query(ctx, schema.RoleQueryParam{
		Status: 1,
	})
	if err != nil {
		return err
	}
	if len(roleResult.Data) == 0 {
		return nil
	}

	roleMenuResult, err := a.roleMenuRepo.Query(ctx, schema.RoleMenuQueryParam{})
	if err != nil {
		return err
	}
	mRoleMenus := roleMenuResult.Data.ToRoleIDMap()

	menuResourceResult, err := a.menuResourceRepo.Query(ctx, schema.MenuActionResourceQueryParam{})
	if err != nil {
		return err
	}
	mMenuResources := menuResourceResult.Data.ToActionIDMap()

	for _, item := range roleResult.Data {
		mcache := make(map[string]struct{})
		if rms, ok := mRoleMenus[item.ID]; ok {
			for _, actionID := range rms.ToActionIDs() {
				if mrs, ok := mMenuResources[actionID]; ok {
					for _, mr := range mrs {
						if mr.Path == " || mr.Method == " {
							continue
						} else if _, ok := mcache[mr.Path+mr.Method]; ok {
							continue
						}
						mcache[mr.Path+mr.Method] = struct{}{}
						line := fmt.Sprintf("p,%s,%s,%s", item.ID, mr.Path, mr.Method)
						persist.LoadPolicyLine(line, m)
					}
				}
			}
		}
	}

	return nil
}

func (a *casbinAdapter) loadUserPolicy(ctx context.Context, m casbinModel.Model) error {
	userResult, err := a.userRepo.Query(ctx, schema.UserQueryParam{
		Status: 1,
	})
	if err != nil {
		return err
	}
	if len(userResult.Data) > 0 {
		userRoleResult, err := a.userRoleRepo.Query(ctx, schema.UserRoleQueryParam{})
		if err != nil {
			return err
		}

		mUserRoles := userRoleResult.Data.ToUserIDMap()
		for _, uitem := range userResult.Data {
			if urs, ok := mUserRoles[uitem.ID]; ok {
				for _, ur := range urs {
					line := fmt.Sprintf("g,%s,%s", ur.UserID, ur.RoleID)
					persist.LoadPolicyLine(line, m)
				}
			}
		}
	}

	return nil
}

// SavePolicy saves all policy rules to the storage.
func (a *casbinAdapter) SavePolicy(model casbinModel.Model) error {
	return nil
}

// AddPolicy adds a policy rule to the storage.
// This is part of the Auto-Save feature.
func (a *casbinAdapter) AddPolicy(sec string, ptype string, rule []string) error {
	return nil
}

// RemovePolicy removes a policy rule from the storage.
// This is part of the Auto-Save feature.
func (a *casbinAdapter) RemovePolicy(sec string, ptype string, rule []string) error {
	return nil
}

// RemoveFilteredPolicy removes policy rules that match the filter from the storage.
// This is part of the Auto-Save feature.
func (a *casbinAdapter) RemoveFilteredPolicy(sec string, ptype string, fieldIndex int, fieldValues ...string) error {
	return nil
}
