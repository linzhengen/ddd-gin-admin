package application

import (
	"context"

	"github.com/linzhengen/ddd-gin-admin/app/domain/user/userrole"

	"github.com/linzhengen/ddd-gin-admin/app/domain/pagination"

	"github.com/linzhengen/ddd-gin-admin/app/domain/user"
	"github.com/linzhengen/ddd-gin-admin/app/domain/user/role"

	"github.com/linzhengen/ddd-gin-admin/app/domain/rbac"
	"github.com/linzhengen/ddd-gin-admin/app/domain/trans"

	"github.com/casbin/casbin/v2"
	"github.com/linzhengen/ddd-gin-admin/app/domain/errors"
	"github.com/linzhengen/ddd-gin-admin/pkg/util/hash"
	"github.com/linzhengen/ddd-gin-admin/pkg/util/uuid"
)

type User interface {
	Query(ctx context.Context, params user.QueryParams) (user.Users, *pagination.Pagination, error)
	QueryShow(ctx context.Context, params user.QueryParams) (user.Users, *pagination.Pagination, error)
	Get(ctx context.Context, id string) (*user.User, error)
	Create(ctx context.Context, item *user.User) (string, error)
	Update(ctx context.Context, id string, item *user.User) error
	Delete(ctx context.Context, id string) error
	UpdateStatus(ctx context.Context, id string, status int) error
}

func NewUser(
	rbacRepo rbac.Repository,
	enforcer *casbin.SyncedEnforcer,
	transRepo trans.Repository,
	userRepo user.Repository,
	userRoleRepo user.Repository,
	roleRepo role.Repository,
) User {
	return &userApp{
		rbacRepo:     rbacRepo,
		enforcer:     enforcer,
		transRepo:    transRepo,
		userRepo:     userRepo,
		userRoleRepo: userRoleRepo,
		roleRepo:     roleRepo,
	}
}

type userApp struct {
	rbacRepo     rbac.Repository
	enforcer     *casbin.SyncedEnforcer
	transRepo    trans.Repository
	userRepo     user.Repository
	userRoleRepo userrole.Repository
	roleRepo     role.Repository
}

func (a *userApp) Query(ctx context.Context, params user.QueryParams) (user.Users, *pagination.Pagination, error) {
	result, pr, err := a.userRepo.Query(ctx, params)
	if err != nil {
		return nil, nil, err
	}
	return result, pr, nil
}

func (a *userApp) QueryShow(ctx context.Context, params user.QueryParams) (user.Users, *pagination.Pagination, error) {
	result, pr, err := a.userRepo.Query(ctx, params)
	if err != nil {
		return nil, nil, err
	}
	if result == nil {
		return nil, nil, nil
	}

	userRoleResult, _, err := a.userRoleRepo.Query(ctx, userrole.QueryParam{
		UserIDs: result.ToIDs(),
	})
	if err != nil {
		return nil, nil, err
	}

	roleResult, _, err := a.roleRepo.Query(ctx, role.QueryParam{
		IDs: userRoleResult.ToRoleIDs(),
	})
	if err != nil {
		return nil, nil, err
	}

	return result, pr, nil
}

func (a *userApp) Get(ctx context.Context, id string) (*user.User, error) {
	item, err := a.userRepo.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	if item == nil {
		return nil, errors.ErrNotFound
	}

	userRoleResult, _, err := a.userRoleRepo.Query(ctx, userrole.QueryParam{
		UserID: id,
	})
	if err != nil {
		return nil, err
	}
	item.Roles = a.userRoleFactory.ToSchemaList(userRoleResult)

	return user, nil
}

func (a *userApp) Create(ctx context.Context, item schema.User) (*schema.IDResult, error) {
	err := a.checkUserName(ctx, item)
	if err != nil {
		return nil, err
	}

	item.Password = hash.SHA1String(item.Password)
	item.ID = uuid.MustString()
	err = a.transRepo.Exec(ctx, func(ctx context.Context) error {
		for _, urItem := range item.UserRoles {
			urItem.ID = uuid.MustString()
			urItem.UserID = item.ID
			err := a.userRoleRepo.Create(ctx, *a.userRoleFactory.ToEntity(urItem))
			if err != nil {
				return err
			}
		}

		return a.userRepo.Create(ctx, *a.userFactory.ToEntity(&item))
	})
	if err != nil {
		return nil, err
	}

	a.casbinAdapter.AddCasbinPolicyItemToChan(ctx, a.enforcer)
	return schema.NewIDResult(item.ID), nil
}

func (a *userApp) checkUserName(ctx context.Context, item schema.User) error {
	if item.UserName == schema.GetRootUser().UserName {
		return errors.New400Response("The user name is invalid")
	}

	_, pr, err := a.userRepo.Query(ctx, schema.UserQueryParam{
		PaginationParam: schema.PaginationParam{OnlyCount: true},
		UserName:        item.UserName,
	})
	if err != nil {
		return err
	}
	if pr.Total > 0 {
		return errors.New400Response("The user name already exists")
	}
	return nil
}

func (a *userApp) Update(ctx context.Context, id string, item schema.User) error {
	oldItem, err := a.Get(ctx, id)
	if err != nil {
		return err
	}
	if oldItem == nil {
		return errors.ErrNotFound
	}
	if oldItem.UserName != item.UserName {
		err := a.checkUserName(ctx, item)
		if err != nil {
			return err
		}
	}

	if item.Password != "" {
		item.Password = hash.SHA1String(item.Password)
	} else {
		item.Password = oldItem.Password
	}

	item.ID = oldItem.ID
	item.Creator = oldItem.Creator
	item.CreatedAt = oldItem.CreatedAt
	err = a.transRepo.Exec(ctx, func(ctx context.Context) error {
		addUserRoles, delUserRoles := a.compareUserRoles(oldItem.UserRoles, item.UserRoles)
		for _, rmitem := range addUserRoles {
			rmitem.ID = uuid.MustString()
			rmitem.UserID = id
			err := a.userRoleRepo.Create(ctx, *a.userRoleFactory.ToEntity(rmitem))
			if err != nil {
				return err
			}
		}

		for _, rmitem := range delUserRoles {
			err := a.userRoleRepo.Delete(ctx, rmitem.ID)
			if err != nil {
				return err
			}
		}

		return a.userRepo.Update(ctx, id, *a.userFactory.ToEntity(&item))
	})
	if err != nil {
		return err
	}

	a.casbinAdapter.AddCasbinPolicyItemToChan(ctx, a.enforcer)
	return nil
}

func (a *userApp) compareUserRoles(oldUserRoles, newUserRoles schema.UserRoles) (addList, delList schema.UserRoles) {
	mOldUserRoles := oldUserRoles.ToMap()
	mNewUserRoles := newUserRoles.ToMap()

	for k, item := range mNewUserRoles {
		if _, ok := mOldUserRoles[k]; ok {
			delete(mOldUserRoles, k)
			continue
		}
		addList = append(addList, item)
	}

	for _, item := range mOldUserRoles {
		delList = append(delList, item)
	}
	return
}

func (a *userApp) Delete(ctx context.Context, id string) error {
	oldItem, err := a.userRepo.Get(ctx, id)
	if err != nil {
		return err
	}
	if oldItem == nil {
		return errors.ErrNotFound
	}

	err = a.transRepo.Exec(ctx, func(ctx context.Context) error {
		err := a.userRoleRepo.DeleteByUserID(ctx, id)
		if err != nil {
			return err
		}

		return a.userRepo.Delete(ctx, id)
	})
	if err != nil {
		return err
	}

	a.casbinAdapter.AddCasbinPolicyItemToChan(ctx, a.enforcer)
	return nil
}

func (a *userApp) UpdateStatus(ctx context.Context, id string, status int) error {
	oldItem, err := a.userRepo.Get(ctx, id)
	if err != nil {
		return err
	}
	if oldItem == nil {
		return errors.ErrNotFound
	}
	oldItem.Status = status

	err = a.userRepo.UpdateStatus(ctx, id, status)
	if err != nil {
		return err
	}

	a.casbinAdapter.AddCasbinPolicyItemToChan(ctx, a.enforcer)
	return nil
}
