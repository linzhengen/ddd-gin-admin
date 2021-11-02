package user

import (
	"context"

	"github.com/linzhengen/ddd-gin-admin/app/domain/errors"

	"github.com/linzhengen/ddd-gin-admin/app/domain/auth"
	"github.com/linzhengen/ddd-gin-admin/app/domain/user/role"
	"github.com/linzhengen/ddd-gin-admin/app/domain/user/userrole"
)

func NewService(
	authRepo auth.Repository,
	userRepo Repository,
	roleRepo role.Repository,
	userRoleRepo userrole.Repository,
) Service {
	return &service{
		authRepo:     authRepo,
		userRepo:     userRepo,
		roleRepo:     roleRepo,
		userRoleRepo: userRoleRepo,
	}
}

type Service interface {
	GetActiveUser(ctx context.Context, userID string) (*User, error)
	GetActiveUserWithRole(ctx context.Context, userID string) (*User, error)
}

type service struct {
	authRepo     auth.Repository
	userRepo     Repository
	roleRepo     role.Repository
	userRoleRepo userrole.Repository
}

func (s service) GetActiveUser(ctx context.Context, userID string) (*User, error) {
	user, err := s.userRepo.Get(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.ErrInvalidUser
	}
	if user.Status != 1 {
		return nil, errors.ErrUserDisable
	}
	return user, nil
}

func (s service) GetActiveUserWithRole(ctx context.Context, userID string) (*User, error) {
	if rootUser := s.authRepo.FindRootUser(ctx, userID); rootUser != nil {
		return &User{
			UserName: rootUser.UserName,
		}, nil
	}

	user, err := s.GetActiveUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	userRoles, _, err := s.userRoleRepo.Query(ctx, userrole.QueryParam{
		UserID: userID,
	})
	if err != nil {
		return nil, err
	}

	roleIDs := userRoles.ToRoleIDs()
	if len(roleIDs) > 0 {
		roles, _, err := s.roleRepo.Query(ctx, role.QueryParam{
			IDs:    roleIDs,
			Status: 1,
		})
		if err != nil {
			return nil, err
		}
		user.Roles = roles
	}

	return user, nil
}
