package application

import (
	"context"

	"github.com/linzhengen/ddd-gin-admin/app/domain/user/role"

	"github.com/linzhengen/ddd-gin-admin/app/domain/user/userrole"

	"github.com/linzhengen/ddd-gin-admin/app/domain/errors"
	"github.com/linzhengen/ddd-gin-admin/pkg/util/hash"

	"github.com/linzhengen/ddd-gin-admin/app/domain/auth"
	"github.com/linzhengen/ddd-gin-admin/app/domain/user"
)

type Login interface {
	Verify(ctx context.Context, userName, password string) (*user.User, error)
}

func NewLogin(
	authRepo auth.Repository,
	userRepo user.Repository,
	roleRepo role.Repository,
	userRoleRepo userrole.Repository,
) Login {
	return &login{
		authRepo:     authRepo,
		userRepo:     userRepo,
		roleRepo:     roleRepo,
		userRoleRepo: userRoleRepo,
	}
}

type login struct {
	authRepo     auth.Repository
	userRepo     user.Repository
	roleRepo     role.Repository
	userRoleRepo userrole.Repository
}

func (l *login) Verify(ctx context.Context, userName, password string) (*user.User, error) {
	if rootUser := l.authRepo.FindRootUser(ctx, userName); rootUser != nil {
		if password == rootUser.Password {
			return &user.User{
				UserName: rootUser.UserName,
				Password: rootUser.Password,
			}, nil
		}
	}
	result, _, err := l.userRepo.Query(ctx, user.QueryParams{
		UserName: userName,
	})
	if err != nil {
		return nil, err
	}
	if len(result) == 0 {
		return nil, errors.ErrInvalidUserName
	}
	item := result[0]
	if item.Password != hash.SHA1String(password) {
		return nil, errors.ErrInvalidPassword
	}
	if item.Status != 1 {
		return nil, errors.ErrUserDisable
	}
	return item, nil
}

func (l *login) GenerateToken(ctx context.Context, userID string) (*auth.Auth, error) {
	auth, err := l.authRepo.GenerateToken(ctx, userID)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return auth, nil
}

func (l *login) DestroyToken(ctx context.Context, tokenString string) error {
	err := l.authRepo.DestroyToken(ctx, tokenString)
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (l *login) checkAndGetUser(ctx context.Context, userID string) (*user.User, error) {
	user, err := l.userRepo.Get(ctx, userID)
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

func (l *login) GetLoginInfo(ctx context.Context, userID string) (*user.User, error) {
	if rootUser := l.authRepo.FindRootUser(ctx, userID); rootUser != nil {
		return &user.User{
			UserName: rootUser.UserName,
		}, nil
	}

	user, err := l.checkAndGetUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	userRoles, _, err := l.userRoleRepo.Query(ctx, userrole.QueryParam{
		UserID: userID,
	})
	if err != nil {
		return nil, err
	}

	roleIDs := userRoles.ToRoleIDs()
	if len(roleIDs) > 0 {
		roles, _, err := l.roleRepo.Query(ctx, role.QueryParam{
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

func (l *login) UpdatePassword(ctx context.Context, userID string, oldPassword, newPassword string) error {
	if rootUser := l.authRepo.FindRootUser(ctx, userID); rootUser != nil {
		return errors.New400Response("The root user is not allowed to update the password")
	}

	user, err := l.checkAndGetUser(ctx, userID)
	if err != nil {
		return err
	} else if hash.SHA1String(oldPassword) != user.Password {
		return errors.New400Response("The old password is invalid")
	}

	return l.userRepo.UpdatePassword(ctx, userID, hash.SHA1String(newPassword))
}
