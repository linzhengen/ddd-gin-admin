package application

import (
	"context"

	"github.com/casbin/casbin/v2"

	"github.com/linzhengen/ddd-gin-admin/app/domain/auth"

	"github.com/linzhengen/ddd-gin-admin/app/domain/user/userrole"

	"github.com/linzhengen/ddd-gin-admin/app/domain/pagination"

	"github.com/linzhengen/ddd-gin-admin/app/domain/user"
	"github.com/linzhengen/ddd-gin-admin/app/domain/user/role"

	"github.com/linzhengen/ddd-gin-admin/app/domain/trans"

	"github.com/linzhengen/ddd-gin-admin/app/domain/errors"
	"github.com/linzhengen/ddd-gin-admin/pkg/util/hash"
	"github.com/linzhengen/ddd-gin-admin/pkg/util/uuid"
)

type User interface {
	Query(ctx context.Context, params user.QueryParams) (user.Users, *pagination.Pagination, error)
	QueryShow(ctx context.Context, params user.QueryParams) (user.Users, *pagination.Pagination, error)
	Get(ctx context.Context, id string) (*user.User, error)
	Create(ctx context.Context, item *user.User, roleIDs []string) (string, error)
	Update(ctx context.Context, id string, item *user.User, roleIDs []string) error
	Delete(ctx context.Context, id string) error
	UpdateStatus(ctx context.Context, id string, status int) error
}

func NewUser(
	authRepo auth.Repository,
	rbacAdapter RbacAdapter,
	enforcer *casbin.SyncedEnforcer,
	transRepo trans.Repository,
	userRepo user.Repository,
	userRoleRepo userrole.Repository,
	roleRepo role.Repository,
) User {
	return &userApp{
		authRepo:     authRepo,
		rbacAdapter:  rbacAdapter,
		enforcer:     enforcer,
		transRepo:    transRepo,
		userRepo:     userRepo,
		userRoleRepo: userRoleRepo,
		roleRepo:     roleRepo,
	}
}

type userApp struct {
	authRepo     auth.Repository
	rbacAdapter  RbacAdapter
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

	return result.FillRoles(userRoleResult.ToUserIDMap(), roleResult.ToMap()), pr, nil
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
	roleResult, _, err := a.roleRepo.Query(ctx, role.QueryParam{
		IDs: userRoleResult.ToRoleIDs(),
	})
	if err != nil {
		return nil, err
	}
	return item.FillRoles(userRoleResult.ToUserIDMap(), roleResult.ToMap()), nil
}

func (a *userApp) Create(ctx context.Context, item *user.User, roleIDs []string) (string, error) {
	err := a.checkUserName(ctx, item)
	if err != nil {
		return "", err
	}

	item.Password = hash.SHA1String(item.Password)
	item.ID = uuid.MustString()
	err = a.transRepo.Exec(ctx, func(ctx context.Context) error {
		for _, roleID := range roleIDs {
			urItem := new(userrole.UserRole)
			urItem.ID = uuid.MustString()
			urItem.UserID = item.ID
			urItem.RoleID = roleID
			err := a.userRoleRepo.Create(ctx, urItem)
			if err != nil {
				return err
			}
		}

		return a.userRepo.Create(ctx, item)
	})
	if err != nil {
		return "", err
	}

	a.rbacAdapter.AddPolicyItemToChan(ctx, a.enforcer)
	return item.ID, nil
}

func (a *userApp) checkUserName(ctx context.Context, item *user.User) error {
	if rootUser := a.authRepo.FindRootUser(ctx, item.UserName); rootUser != nil {
		return errors.New400Response("The user name is invalid")
	}
	_, pr, err := a.userRepo.Query(ctx, user.QueryParams{
		PaginationParam: pagination.Param{OnlyCount: true},
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

func (a *userApp) Update(ctx context.Context, id string, item *user.User, roleIDs []string) error {
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

	userRoleResult, _, err := a.userRoleRepo.Query(ctx, userrole.QueryParam{
		UserID: id,
	})
	if err != nil {
		return err
	}

	err = a.transRepo.Exec(ctx, func(ctx context.Context) error {
		addRoleIDs, delUserRoles := a.compareUserRoles(userRoleResult.ToMap(), roleIDs)
		for _, roleID := range addRoleIDs {
			urItem := new(userrole.UserRole)
			urItem.ID = uuid.MustString()
			urItem.UserID = item.ID
			urItem.RoleID = roleID
			err := a.userRoleRepo.Create(ctx, urItem)
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

		return a.userRepo.Update(ctx, id, item)
	})
	if err != nil {
		return err
	}

	a.rbacAdapter.AddPolicyItemToChan(ctx, a.enforcer)
	return nil
}

func (a *userApp) compareUserRoles(oldUserRoles map[string]*userrole.UserRole, roleIDs []string) (addList []string, delList []*userrole.UserRole) {
	for _, roleID := range roleIDs {
		if _, ok := oldUserRoles[roleID]; ok {
			delete(oldUserRoles, roleID)
			continue
		}
		addList = append(addList, roleID)
	}

	for _, item := range oldUserRoles {
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

	a.rbacAdapter.AddPolicyItemToChan(ctx, a.enforcer)
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

	a.rbacAdapter.AddPolicyItemToChan(ctx, a.enforcer)
	return nil
}
