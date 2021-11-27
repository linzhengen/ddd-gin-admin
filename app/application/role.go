package application

import (
	"context"

	"github.com/casbin/casbin/v2"

	"github.com/linzhengen/ddd-gin-admin/app/domain/errors"
	"github.com/linzhengen/ddd-gin-admin/app/domain/pagination"
	"github.com/linzhengen/ddd-gin-admin/app/domain/trans"
	"github.com/linzhengen/ddd-gin-admin/app/domain/user"
	"github.com/linzhengen/ddd-gin-admin/app/domain/user/role"
	"github.com/linzhengen/ddd-gin-admin/app/domain/user/rolemenu"
	"github.com/linzhengen/ddd-gin-admin/pkg/util/uuid"
)

type Role interface {
	Query(ctx context.Context, params role.QueryParam) (role.Roles, *pagination.Pagination, error)
	Get(ctx context.Context, id string) (*role.Role, error)
	QueryRoleMenus(ctx context.Context, roleID string) (rolemenu.RoleMenus, error)
	Create(ctx context.Context, item *role.Role) (string, error)
	Update(ctx context.Context, id string, item *role.Role) error
	Delete(ctx context.Context, id string) error
	UpdateStatus(ctx context.Context, id string, status int) error
}

func NewRole(
	rbacAdapter RbacAdapter,
	enforcer *casbin.SyncedEnforcer,
	transRepo trans.Repository,
	roleRepo role.Repository,
	roleMenuRepo rolemenu.Repository,
	userRepo user.Repository,
) Role {
	return &roleApp{
		rbacAdapter:  rbacAdapter,
		enforcer:     enforcer,
		transRepo:    transRepo,
		roleRepo:     roleRepo,
		roleMenuRepo: roleMenuRepo,
		userRepo:     userRepo,
	}
}

type roleApp struct {
	rbacAdapter  RbacAdapter
	enforcer     *casbin.SyncedEnforcer
	transRepo    trans.Repository
	roleRepo     role.Repository
	roleMenuRepo rolemenu.Repository
	userRepo     user.Repository
}

func (a *roleApp) Query(ctx context.Context, params role.QueryParam) (role.Roles, *pagination.Pagination, error) {
	result, pr, err := a.roleRepo.Query(ctx, params)
	if err != nil {
		return nil, nil, err
	}
	return result, pr, nil
}

func (a *roleApp) Get(ctx context.Context, id string) (*role.Role, error) {
	item, err := a.roleRepo.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	if item == nil {
		return nil, errors.ErrNotFound
	}

	roleMenus, err := a.QueryRoleMenus(ctx, id)
	if err != nil {
		return nil, err
	}
	item.RoleMenus = roleMenus

	return item, nil
}

func (a *roleApp) QueryRoleMenus(ctx context.Context, roleID string) (rolemenu.RoleMenus, error) {
	result, _, err := a.roleMenuRepo.Query(ctx, rolemenu.QueryParam{
		RoleID: roleID,
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (a *roleApp) Create(ctx context.Context, item *role.Role) (string, error) {
	err := a.checkName(ctx, item)
	if err != nil {
		return "", err
	}

	item.ID = uuid.MustString()
	err = a.transRepo.Exec(ctx, func(ctx context.Context) error {
		for _, rmItem := range item.RoleMenus {
			rmItem.ID = uuid.MustString()
			rmItem.RoleID = item.ID
			err := a.roleMenuRepo.Create(ctx, rmItem)
			if err != nil {
				return err
			}
		}
		return a.roleRepo.Create(ctx, item)
	})
	if err != nil {
		return "", err
	}
	a.rbacAdapter.AddPolicyItemToChan(ctx, a.enforcer)
	return item.ID, nil
}

func (a *roleApp) checkName(ctx context.Context, item *role.Role) error {
	_, pr, err := a.roleRepo.Query(ctx, role.QueryParam{
		PaginationParam: pagination.Param{OnlyCount: true},
		Name:            item.Name,
	})
	if err != nil {
		return err
	}
	if pr.Total > 0 {
		return errors.New400Response("The role name already exists")
	}
	return nil
}

func (a *roleApp) Update(ctx context.Context, id string, item *role.Role) error {
	oldItem, err := a.Get(ctx, id)
	if err != nil {
		return err
	}
	if oldItem == nil {
		return errors.ErrNotFound
	}
	if oldItem.Name != item.Name {
		err := a.checkName(ctx, item)
		if err != nil {
			return err
		}
	}

	item.ID = oldItem.ID
	item.Creator = oldItem.Creator
	item.CreatedAt = oldItem.CreatedAt
	err = a.transRepo.Exec(ctx, func(ctx context.Context) error {
		addRoleMenus, delRoleMenus := a.compareRoleMenus(oldItem.RoleMenus, item.RoleMenus)
		for _, rmitem := range addRoleMenus {
			rmitem.ID = uuid.MustString()
			rmitem.RoleID = id
			err := a.roleMenuRepo.Create(ctx, rmitem)
			if err != nil {
				return err
			}
		}

		for _, rmitem := range delRoleMenus {
			err := a.roleMenuRepo.Delete(ctx, rmitem.ID)
			if err != nil {
				return err
			}
		}

		return a.roleRepo.Update(ctx, id, item)
	})
	if err != nil {
		return err
	}
	a.rbacAdapter.AddPolicyItemToChan(ctx, a.enforcer)
	return nil
}

func (a *roleApp) compareRoleMenus(oldRoleMenus, newRoleMenus rolemenu.RoleMenus) (addList, delList rolemenu.RoleMenus) {
	mOldRoleMenus := oldRoleMenus.ToMap()
	mNewRoleMenus := newRoleMenus.ToMap()

	for k, item := range mNewRoleMenus {
		if _, ok := mOldRoleMenus[k]; ok {
			delete(mOldRoleMenus, k)
			continue
		}
		addList = append(addList, item)
	}

	for _, item := range mOldRoleMenus {
		delList = append(delList, item)
	}
	return
}

func (a *roleApp) Delete(ctx context.Context, id string) error {
	oldItem, err := a.roleRepo.Get(ctx, id)
	if err != nil {
		return err
	}
	if oldItem == nil {
		return errors.ErrNotFound
	}

	_, pr, err := a.userRepo.Query(ctx, user.QueryParams{
		PaginationParam: pagination.Param{OnlyCount: true},
		RoleIDs:         []string{id},
	})
	if err != nil {
		return err
	}
	if pr.Total > 0 {
		return errors.New400Response("The role has been assigned to the user and cannot be deleted")
	}

	err = a.transRepo.Exec(ctx, func(ctx context.Context) error {
		err := a.roleMenuRepo.DeleteByRoleID(ctx, id)
		if err != nil {
			return err
		}

		return a.roleRepo.Delete(ctx, id)
	})
	if err != nil {
		return err
	}

	a.rbacAdapter.AddPolicyItemToChan(ctx, a.enforcer)
	return nil
}

func (a *roleApp) UpdateStatus(ctx context.Context, id string, status int) error {
	oldItem, err := a.roleRepo.Get(ctx, id)
	if err != nil {
		return err
	}
	if oldItem == nil {
		return errors.ErrNotFound
	}

	err = a.roleRepo.UpdateStatus(ctx, id, status)
	if err != nil {
		return err
	}
	a.rbacAdapter.AddPolicyItemToChan(ctx, a.enforcer)
	return nil
}
