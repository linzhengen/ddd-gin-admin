package service

import (
	"context"

	"github.com/linzhengen/ddd-gin-admin/app/domain/factory"

	"github.com/linzhengen/ddd-gin-admin/app/domain/valueobject/errors"

	"github.com/linzhengen/ddd-gin-admin/app/domain/repository"
	"github.com/linzhengen/ddd-gin-admin/app/domain/valueobject/schema"

	"github.com/casbin/casbin/v2"
	"github.com/linzhengen/ddd-gin-admin/pkg/util/hash"
	"github.com/linzhengen/ddd-gin-admin/pkg/util/uuid"
)

type User interface {
	Query(ctx context.Context, params schema.UserQueryParam) (*schema.UserQueryResult, error)
	QueryShow(ctx context.Context, params schema.UserQueryParam) (*schema.UserShowQueryResult, error)
	Get(ctx context.Context, id string) (*schema.User, error)
	Create(ctx context.Context, item schema.User) (*schema.IDResult, error)
	Update(ctx context.Context, id string, item schema.User) error
	Delete(ctx context.Context, id string) error
	UpdateStatus(ctx context.Context, id string, status int) error
}

func NewUser(
	casbinAdapter repository.CasbinAdapter,
	enforcer *casbin.SyncedEnforcer,
	transRepo repository.TransRepository,
	userRepo repository.UserRepository,
	userRoleRepo repository.UserRoleRepository,
	roleRepo repository.RoleRepository,
	userFactory factory.User,
	userRoleFactory factory.UserRole,
	roleFactory factory.Role,
) User {
	return &user{
		casbinAdapter:   casbinAdapter,
		enforcer:        enforcer,
		transRepo:       transRepo,
		userRepo:        userRepo,
		userRoleRepo:    userRoleRepo,
		roleRepo:        roleRepo,
		userFactory:     userFactory,
		userRoleFactory: userRoleFactory,
		roleFactory:     roleFactory,
	}
}

type user struct {
	casbinAdapter   repository.CasbinAdapter
	enforcer        *casbin.SyncedEnforcer
	transRepo       repository.TransRepository
	userRepo        repository.UserRepository
	userRoleRepo    repository.UserRoleRepository
	roleRepo        repository.RoleRepository
	userFactory     factory.User
	userRoleFactory factory.UserRole
	roleFactory     factory.Role
}

func (a *user) Query(ctx context.Context, params schema.UserQueryParam) (*schema.UserQueryResult, error) {
	result, pr, err := a.userRepo.Query(ctx, params)
	if err != nil {
		return nil, err
	}
	return &schema.UserQueryResult{
		Data:       a.userFactory.ToSchemaList(result),
		PageResult: pr,
	}, nil
}

func (a *user) QueryShow(ctx context.Context, params schema.UserQueryParam) (*schema.UserShowQueryResult, error) {
	result, pr, err := a.userRepo.Query(ctx, params)
	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, nil
	}

	users := a.userFactory.ToSchemaList(result)
	userRoleResult, _, err := a.userRoleRepo.Query(ctx, schema.UserRoleQueryParam{
		UserIDs: users.ToIDs(),
	})
	if err != nil {
		return nil, err
	}

	roleResult, _, err := a.roleRepo.Query(ctx, schema.RoleQueryParam{
		IDs: a.userRoleFactory.ToSchemaList(userRoleResult).ToRoleIDs(),
	})
	if err != nil {
		return nil, err
	}

	return schema.UserQueryResult{
		Data:       users,
		PageResult: pr,
	}.ToShowResult(a.userRoleFactory.ToSchemaList(userRoleResult).ToUserIDMap(), a.roleFactory.ToSchemaList(roleResult).ToMap()), nil
}

func (a *user) Get(ctx context.Context, id string) (*schema.User, error) {
	item, err := a.userRepo.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	if item == nil {
		return nil, errors.ErrNotFound
	}

	userRoleResult, _, err := a.userRoleRepo.Query(ctx, schema.UserRoleQueryParam{
		UserID: id,
	})
	if err != nil {
		return nil, err
	}
	user := a.userFactory.ToSchema(item)
	user.UserRoles = a.userRoleFactory.ToSchemaList(userRoleResult)

	return user, nil
}

func (a *user) Create(ctx context.Context, item schema.User) (*schema.IDResult, error) {
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

func (a *user) checkUserName(ctx context.Context, item schema.User) error {
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

func (a *user) Update(ctx context.Context, id string, item schema.User) error {
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

func (a *user) compareUserRoles(oldUserRoles, newUserRoles schema.UserRoles) (addList, delList schema.UserRoles) {
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

func (a *user) Delete(ctx context.Context, id string) error {
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

func (a *user) UpdateStatus(ctx context.Context, id string, status int) error {
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
