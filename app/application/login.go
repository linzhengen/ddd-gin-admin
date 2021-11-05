package application

import (
	"context"

	"github.com/linzhengen/ddd-gin-admin/app/domain/auth"
	"github.com/linzhengen/ddd-gin-admin/app/domain/errors"
	"github.com/linzhengen/ddd-gin-admin/app/domain/menu"
	"github.com/linzhengen/ddd-gin-admin/app/domain/menu/menuaction"
	"github.com/linzhengen/ddd-gin-admin/app/domain/pagination"
	"github.com/linzhengen/ddd-gin-admin/app/domain/user"
	"github.com/linzhengen/ddd-gin-admin/app/domain/user/role"
	"github.com/linzhengen/ddd-gin-admin/app/domain/user/rolemenu"
	"github.com/linzhengen/ddd-gin-admin/app/domain/user/userrole"
	"github.com/linzhengen/ddd-gin-admin/pkg/util/hash"
)

type Login interface {
	Verify(ctx context.Context, userName, password string) (*user.User, error)
	GenerateToken(ctx context.Context, userID string) (*auth.Auth, error)
	DestroyToken(ctx context.Context, tokenString string) error
	GetLoginInfo(ctx context.Context, userID string) (*user.User, error)
	UpdatePassword(ctx context.Context, userID string, oldPassword, newPassword string) error
	QueryUserMenuTree(ctx context.Context, userID string) (menu.Menus, error)
}

func NewLogin(
	authRepo auth.Repository,
	userRepo user.Repository,
	roleRepo role.Repository,
	userRoleRepo userrole.Repository,
	userSvc user.Service,
	menuRepo menu.Repository,
	menuActionRepo menuaction.Repository,
	roleMenuRepo rolemenu.Repository,
) Login {
	return &loginApp{
		authRepo:       authRepo,
		userRepo:       userRepo,
		roleRepo:       roleRepo,
		userRoleRepo:   userRoleRepo,
		userSvc:        userSvc,
		menuRepo:       menuRepo,
		menuActionRepo: menuActionRepo,
		roleMenuRepo:   roleMenuRepo,
	}
}

type loginApp struct {
	authRepo       auth.Repository
	userRepo       user.Repository
	roleRepo       role.Repository
	userRoleRepo   userrole.Repository
	userSvc        user.Service
	menuRepo       menu.Repository
	menuActionRepo menuaction.Repository
	roleMenuRepo   rolemenu.Repository
}

func (l loginApp) Verify(ctx context.Context, userName, password string) (*user.User, error) {
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

func (l loginApp) GenerateToken(ctx context.Context, userID string) (*auth.Auth, error) {
	auth, err := l.authRepo.GenerateToken(ctx, userID)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return auth, nil
}

func (l loginApp) DestroyToken(ctx context.Context, tokenString string) error {
	err := l.authRepo.DestroyToken(ctx, tokenString)
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (l loginApp) GetLoginInfo(ctx context.Context, userID string) (*user.User, error) {
	return l.userSvc.GetActiveUserWithRole(ctx, userID)
}

func (l loginApp) UpdatePassword(ctx context.Context, userID string, oldPassword, newPassword string) error {
	if rootUser := l.authRepo.FindRootUser(ctx, userID); rootUser != nil {
		return errors.New400Response("The root user is not allowed to update the password")
	}

	user, err := l.userSvc.GetActiveUser(ctx, userID)
	if err != nil {
		return err
	} else if hash.SHA1String(oldPassword) != user.Password {
		return errors.New400Response("The old password is invalid")
	}

	return l.userRepo.UpdatePassword(ctx, userID, hash.SHA1String(newPassword))
}

func (l loginApp) QueryUserMenuTree(ctx context.Context, userID string) (menu.Menus, error) {
	isRoot := false
	if rootUser := l.authRepo.FindRootUser(ctx, userID); rootUser != nil {
		isRoot = true
	}
	// show all menu when root user
	if isRoot {
		menuResult, _, err := l.menuRepo.Query(ctx, menu.QueryParam{
			OrderFields: pagination.NewOrderFields(pagination.NewOrderField("sequence", pagination.OrderByDESC)),
			Status:      1,
		})
		if err != nil {
			return nil, err
		}

		menuActionResult, _, err := l.menuActionRepo.Query(ctx, menuaction.QueryParam{})
		if err != nil {
			return nil, err
		}
		return menuResult.FillMenuAction(menuActionResult.ToMenuIDMap()), nil
	}

	userRoleResult, _, err := l.userRoleRepo.Query(ctx, userrole.QueryParam{
		UserID: userID,
	})
	if err != nil {
		return nil, err
	}
	if len(userRoleResult) == 0 {
		return nil, errors.ErrNoPerm
	}

	roleMenuResult, _, err := l.roleMenuRepo.Query(ctx, rolemenu.QueryParam{
		RoleIDs: userRoleResult.ToRoleIDs(),
	})
	if err != nil {
		return nil, err
	}
	if len(roleMenuResult) == 0 {
		return nil, errors.ErrNoPerm
	}

	menuResult, _, err := l.menuRepo.Query(ctx, menu.QueryParam{
		IDs:    roleMenuResult.ToMenuIDs(),
		Status: 1,
	})
	if err != nil {
		return nil, err
	}
	if len(menuResult) == 0 {
		return nil, errors.ErrNoPerm
	}

	qIDs := menuResult.SplitParentIDs()
	if len(qIDs) > 0 {
		pmenuResult, _, err := l.menuRepo.Query(ctx, menu.QueryParam{
			IDs: qIDs,
		})
		if err != nil {
			return nil, err
		}
		menuResult = append(menuResult, pmenuResult...)
	}

	menuActionResult, _, err := l.menuActionRepo.Query(ctx, menuaction.QueryParam{
		IDs: roleMenuResult.ToActionIDs(),
	})
	if err != nil {
		return nil, err
	}

	return menuResult.FillMenuAction(menuActionResult.ToMenuIDMap()), nil
}
