package rbac

import (
	"context"
	"fmt"

	"github.com/linzhengen/ddd-gin-admin/app/domain/rbac"

	"github.com/linzhengen/ddd-gin-admin/app/domain/menu/menuactionresource"
	"github.com/linzhengen/ddd-gin-admin/app/domain/user"
	"github.com/linzhengen/ddd-gin-admin/app/domain/user/role"
	"github.com/linzhengen/ddd-gin-admin/app/domain/user/rolemenu"
	"github.com/linzhengen/ddd-gin-admin/app/domain/user/userrole"
)

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

func (a *repository) ListRolesPolices(ctx context.Context) ([]string, error) {
	roleResult, _, err := a.roleRepo.Query(ctx, role.QueryParam{
		Status: 1,
	})
	if err != nil {
		return nil, err
	}
	if len(roleResult) == 0 {
		return nil, nil
	}

	roleMenuResult, _, err := a.roleMenuRepo.Query(ctx, rolemenu.QueryParam{})
	if err != nil {
		return nil, err
	}
	mRoleMenus := roleMenuResult.ToRoleIDMap()

	menuResourceResult, _, err := a.menuResourceRepo.Query(ctx, menuactionresource.QueryParam{})
	if err != nil {
		return nil, err
	}
	mMenuResources := menuResourceResult.ToActionIDMap()

	var policies []string
	for _, item := range roleResult {
		mcache := make(map[string]struct{})
		if rms, ok := mRoleMenus[item.ID]; ok {
			for _, actionID := range rms.ToActionIDs() {
				if mrs, ok := mMenuResources[actionID]; ok {
					for _, mr := range mrs {
						if mr.Path == "" || mr.Method == "" {
							continue
						} else if _, ok := mcache[mr.Path+mr.Method]; ok {
							continue
						}
						mcache[mr.Path+mr.Method] = struct{}{}
						policy := fmt.Sprintf("p,%s,%s,%s", item.ID, mr.Path, mr.Method)
						policies = append(policies, policy)
					}
				}
			}
		}
	}

	return policies, nil
}

func (a *repository) ListUsersPolices(ctx context.Context) ([]string, error) {
	userResult, _, err := a.userRepo.Query(ctx, user.QueryParams{
		Status: 1,
	})
	if err != nil {
		return nil, err
	}
	var policies []string
	if len(userResult) > 0 {
		userRoleResult, _, err := a.userRoleRepo.Query(ctx, userrole.QueryParam{})
		if err != nil {
			return nil, err
		}

		mUserRoles := userRoleResult.ToUserIDMap()
		for _, uitem := range userResult {
			if urs, ok := mUserRoles[uitem.ID]; ok {
				for _, ur := range urs {
					policy := fmt.Sprintf("g,%s,%s", ur.UserID, ur.RoleID)
					policies = append(policies, policy)
				}
			}
		}
	}

	return policies, nil
}
