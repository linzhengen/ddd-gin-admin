package rbac

import (
	"context"
	"fmt"

	"github.com/casbin/casbin/v2/persist"

	"github.com/linzhengen/ddd-gin-admin/app/domain/rbac"

	"github.com/linzhengen/ddd-gin-admin/app/domain/menu/menuactionresource"
	"github.com/linzhengen/ddd-gin-admin/app/domain/user"
	"github.com/linzhengen/ddd-gin-admin/app/domain/user/role"
	"github.com/linzhengen/ddd-gin-admin/app/domain/user/rolemenu"
	"github.com/linzhengen/ddd-gin-admin/app/domain/user/userrole"
	"github.com/linzhengen/ddd-gin-admin/configs"
	"github.com/linzhengen/ddd-gin-admin/pkg/logger"

	casbinModel "github.com/casbin/casbin/v2/model"
)

var autoLoadPolicyChan rbac.AutoLoadPolicyChan

func NewRepository(
	roleRepo role.Repository,
	roleMenuRepo rolemenu.Repository,
	menuResourceRepo menuactionresource.Repository,
	userRepo user.Repository,
	userRoleRepo userrole.Repository,
) rbac.Repository {
	return &repository{
		roleRepo:         roleRepo,
		roleMenuRepo:     roleMenuRepo,
		menuResourceRepo: menuResourceRepo,
		userRepo:         userRepo,
		userRoleRepo:     userRoleRepo,
	}
}

type repository struct {
	roleRepo         role.Repository
	roleMenuRepo     rolemenu.Repository
	menuResourceRepo menuactionresource.Repository
	userRepo         user.Repository
	userRoleRepo     userrole.Repository
}

func (a *repository) CreateAutoLoadPolicyChan() rbac.AutoLoadPolicyChan {
	autoLoadPolicyChan = make(chan *rbac.Policy, 1)
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

func (a *repository) GetAutoLoadPolicyChan() rbac.AutoLoadPolicyChan {
	return autoLoadPolicyChan
}

func (a *repository) AddPolicyItemToChan(ctx context.Context, e rbac.SyncedEnforcer) {
	if !configs.C.Casbin.Enable {
		return
	}

	if len(autoLoadPolicyChan) > 0 {
		logger.WithContext(ctx).Infof("The load casbin policy is already in the wait queue")
		return
	}

	autoLoadPolicyChan <- &rbac.Policy{
		Ctx:      ctx,
		Enforcer: e,
	}
}

func (a *repository) LoadPolicy(model casbinModel.Model) error {
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

func (a *repository) loadRolePolicy(ctx context.Context, m casbinModel.Model) error {
	roleResult, _, err := a.roleRepo.Query(ctx, role.QueryParam{
		Status: 1,
	})
	if err != nil {
		return err
	}
	if len(roleResult) == 0 {
		return nil
	}

	roleMenuResult, _, err := a.roleMenuRepo.Query(ctx, rolemenu.QueryParam{})
	if err != nil {
		return err
	}
	mRoleMenus := roleMenuResult.ToRoleIDMap()

	menuResourceResult, _, err := a.menuResourceRepo.Query(ctx, menuactionresource.QueryParam{})
	if err != nil {
		return err
	}
	mMenuResources := menuResourceResult.ToActionIDMap()

	for _, item := range roleResult {
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

func (a *repository) loadUserPolicy(ctx context.Context, m casbinModel.Model) error {
	userResult, _, err := a.userRepo.Query(ctx, user.QueryParams{
		Status: 1,
	})
	if err != nil {
		return err
	}
	if len(userResult) > 0 {
		userRoleResult, _, err := a.userRoleRepo.Query(ctx, userrole.QueryParam{})
		if err != nil {
			return err
		}

		mUserRoles := userRoleResult.ToUserIDMap()
		for _, uitem := range userResult {
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
