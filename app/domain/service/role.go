package service

import (
	"context"

	"github.com/linzhengen/ddd-gin-admin/app/domain/factory"

	"github.com/linzhengen/ddd-gin-admin/app/domain/repository"
	"github.com/linzhengen/ddd-gin-admin/app/domain/valueobject/errors"
	"github.com/linzhengen/ddd-gin-admin/app/domain/valueobject/schema"

	"github.com/casbin/casbin/v2"
	"github.com/linzhengen/ddd-gin-admin/pkg/util/uuid"
)

type Role interface {
	Query(ctx context.Context, params schema.RoleQueryParam) (*schema.RoleQueryResult, error)
	Get(ctx context.Context, id string) (*schema.Role, error)
	QueryRoleMenus(ctx context.Context, roleID string) (schema.RoleMenus, error)
	Create(ctx context.Context, item schema.Role) (*schema.IDResult, error)
	Update(ctx context.Context, id string, item schema.Role) error
	Delete(ctx context.Context, id string) error
	UpdateStatus(ctx context.Context, id string, status int) error
}

func NewRole(
	casbinAdapter repository.CasbinAdapter,
	enforcer *casbin.SyncedEnforcer,
	transRepo repository.TransRepository,
	roleRepo repository.RoleRepository,
	roleMenuRepo repository.RoleMenuRepository,
	userRepo repository.UserRepository,
	roleFactory factory.Role,
	roleMenuFactory factory.RoleMenu,
	userFactory factory.User,
) Role {
	return &role{
		casbinAdapter:   casbinAdapter,
		enforcer:        enforcer,
		transRepo:       transRepo,
		roleRepo:        roleRepo,
		roleMenuRepo:    roleMenuRepo,
		userRepo:        userRepo,
		roleFactory:     roleFactory,
		roleMenuFactory: roleMenuFactory,
		userFactory:     userFactory,
	}
}

type role struct {
	casbinAdapter   repository.CasbinAdapter
	enforcer        *casbin.SyncedEnforcer
	transRepo       repository.TransRepository
	roleRepo        repository.RoleRepository
	roleMenuRepo    repository.RoleMenuRepository
	userRepo        repository.UserRepository
	roleFactory     factory.Role
	roleMenuFactory factory.RoleMenu
	userFactory     factory.User
}

func (a *role) Query(ctx context.Context, params schema.RoleQueryParam) (*schema.RoleQueryResult, error) {
	result, pr, err := a.roleRepo.Query(ctx, params)
	if err != nil {
		return nil, err
	}
	return &schema.RoleQueryResult{
		Data:       a.roleFactory.ToSchemaList(result),
		PageResult: pr,
	}, nil
}

func (a *role) Get(ctx context.Context, id string) (*schema.Role, error) {
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
	roleSchema := a.roleFactory.ToSchema(item)
	roleSchema.RoleMenus = roleMenus

	return roleSchema, nil
}

func (a *role) QueryRoleMenus(ctx context.Context, roleID string) (schema.RoleMenus, error) {
	result, _, err := a.roleMenuRepo.Query(ctx, schema.RoleMenuQueryParam{
		RoleID: roleID,
	})
	if err != nil {
		return nil, err
	}
	return a.roleMenuFactory.ToSchemaList(result), nil
}

func (a *role) Create(ctx context.Context, item schema.Role) (*schema.IDResult, error) {
	err := a.checkName(ctx, item)
	if err != nil {
		return nil, err
	}

	item.ID = uuid.MustString()
	err = a.transRepo.Exec(ctx, func(ctx context.Context) error {
		for _, rmItem := range item.RoleMenus {
			rmItem.ID = uuid.MustString()
			rmItem.RoleID = item.ID
			err := a.roleMenuRepo.Create(ctx, *a.roleMenuFactory.ToEntity(rmItem))
			if err != nil {
				return err
			}
		}
		return a.roleRepo.Create(ctx, *a.roleFactory.ToEntity(&item))
	})
	if err != nil {
		return nil, err
	}
	a.casbinAdapter.AddCasbinPolicyItemToChan(ctx, a.enforcer)
	return schema.NewIDResult(item.ID), nil
}

func (a *role) checkName(ctx context.Context, item schema.Role) error {
	_, pr, err := a.roleRepo.Query(ctx, schema.RoleQueryParam{
		PaginationParam: schema.PaginationParam{OnlyCount: true},
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

func (a *role) Update(ctx context.Context, id string, item schema.Role) error {
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
			err := a.roleMenuRepo.Create(ctx, *a.roleMenuFactory.ToEntity(rmitem))
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

		return a.roleRepo.Update(ctx, id, *a.roleFactory.ToEntity(&item))
	})
	if err != nil {
		return err
	}
	a.casbinAdapter.AddCasbinPolicyItemToChan(ctx, a.enforcer)
	return nil
}

func (a *role) compareRoleMenus(oldRoleMenus, newRoleMenus schema.RoleMenus) (addList, delList schema.RoleMenus) {
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

func (a *role) Delete(ctx context.Context, id string) error {
	oldItem, err := a.roleRepo.Get(ctx, id)
	if err != nil {
		return err
	}
	if oldItem == nil {
		return errors.ErrNotFound
	}

	_, pr, err := a.userRepo.Query(ctx, schema.UserQueryParam{
		PaginationParam: schema.PaginationParam{OnlyCount: true},
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

	a.casbinAdapter.AddCasbinPolicyItemToChan(ctx, a.enforcer)
	return nil
}

func (a *role) UpdateStatus(ctx context.Context, id string, status int) error {
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
	a.casbinAdapter.AddCasbinPolicyItemToChan(ctx, a.enforcer)
	return nil
}
