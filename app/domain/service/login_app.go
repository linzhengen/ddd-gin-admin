package service

import (
	"context"
	"net/http"
	"sort"

	errors2 "github.com/linzhengen/ddd-gin-admin/app/domain/valueobject/errors"

	"github.com/linzhengen/ddd-gin-admin/pkg/util/structure"

	"github.com/linzhengen/ddd-gin-admin/app/domain/valueobject/request"

	"github.com/linzhengen/ddd-gin-admin/app/domain/valueobject/response"

	"github.com/linzhengen/ddd-gin-admin/app/domain/repository"

	"github.com/LyricTian/captcha"
	"github.com/linzhengen/ddd-gin-admin/pkg/auth"
	"github.com/linzhengen/ddd-gin-admin/pkg/util/hash"
)

type Login interface {
	GetCaptcha(ctx context.Context, length int) (*response.LoginCaptcha, error)
	ResCaptcha(ctx context.Context, w http.ResponseWriter, captchaID string, width, height int) error
	Verify(ctx context.Context, userName, password string) (*response.User, error)
	GenerateToken(ctx context.Context, userID string) (*response.LoginToken, error)
	DestroyToken(ctx context.Context, tokenString string) error
	GetLoginInfo(ctx context.Context, userID string) (*response.UserLogin, error)
	QueryUserMenuTree(ctx context.Context, userID string) (response.MenuTrees, error)
	UpdatePassword(ctx context.Context, userID string, params response.UpdatePassword) error
}

func NewLogin(
	auth auth.Author,
	userRepo repository.UserRepository,
	userRoleRepo repository.UserRoleRepository,
	roleRepo repository.RoleRepository,
	roleMenuRepo repository.RoleMenuRepository,
	menuRepo repository.MenuRepository,
	menuActionRepo repository.MenuActionRepository,
) Login {
	return &login{
		auth:           auth,
		userRepo:       userRepo,
		userRoleRepo:   userRoleRepo,
		roleRepo:       roleRepo,
		roleMenuRepo:   roleMenuRepo,
		menuRepo:       menuRepo,
		menuActionRepo: menuActionRepo,
	}
}

type login struct {
	auth           auth.Author
	userRepo       repository.UserRepository
	userRoleRepo   repository.UserRoleRepository
	roleRepo       repository.RoleRepository
	roleMenuRepo   repository.RoleMenuRepository
	menuRepo       repository.MenuRepository
	menuActionRepo repository.MenuActionRepository
}

func (a *login) GetCaptcha(ctx context.Context, length int) (*response.LoginCaptcha, error) {
	captchaID := captcha.NewLen(length)
	item := &response.LoginCaptcha{
		CaptchaID: captchaID,
	}
	return item, nil
}

func (a *login) ResCaptcha(ctx context.Context, w http.ResponseWriter, captchaID string, width, height int) error {
	err := captcha.WriteImage(w, captchaID, width, height)
	if err != nil {
		if err == captcha.ErrNotFound {
			return errors2.ErrNotFound
		}
		return errors2.WithStack(err)
	}

	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")
	w.Header().Set("Content-Type", "image/png")
	return nil
}

func (a *login) Verify(ctx context.Context, userName, password string) (*response.User, error) {
	// is root user
	root := getRootUser()
	if userName == root.UserName && root.Password == password {
		return root, nil
	}

	result, _, err := a.userRepo.Query(ctx, request.UserQuery{
		UserName: userName,
	})
	if err != nil {
		return nil, err
	}
	if len(result) == 0 {
		return nil, errors2.ErrInvalidUserName
	}

	item := result[0]
	if item.Password != hash.SHA1String(password) {
		return nil, errors2.ErrInvalidPassword
	}
	if item.Status != 1 {
		return nil, errors2.ErrUserDisable
	}
	res := new(response.User)
	structure.Copy(item, res)
	return res, nil
}

func (a *login) GenerateToken(ctx context.Context, userID string) (*response.LoginToken, error) {
	tokenInfo, err := a.auth.GenerateToken(ctx, userID)
	if err != nil {
		return nil, errors2.WithStack(err)
	}

	item := &response.LoginToken{
		AccessToken: tokenInfo.GetAccessToken(),
		TokenType:   tokenInfo.GetTokenType(),
		ExpiresAt:   tokenInfo.GetExpiresAt(),
	}
	return item, nil
}

func (a *login) DestroyToken(ctx context.Context, tokenString string) error {
	err := a.auth.DestroyToken(ctx, tokenString)
	if err != nil {
		return errors2.WithStack(err)
	}
	return nil
}

func (a *login) checkAndGetUser(ctx context.Context, userID string) (*response.User, error) {
	user, err := a.userRepo.Get(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors2.ErrInvalidUser
	}
	if user.Status != 1 {
		return nil, errors2.ErrUserDisable
	}
	res := new(response.User)
	structure.Copy(user, res)
	return res, nil
}

func (a *login) GetLoginInfo(ctx context.Context, userID string) (*response.UserLogin, error) {
	if isRoot := checkIsRootUser(userID); isRoot {
		root := getRootUser()
		loginInfo := &response.UserLogin{
			UserName: root.UserName,
			RealName: root.RealName,
		}
		return loginInfo, nil
	}

	user, err := a.checkAndGetUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	info := &response.UserLogin{
		UserID:   user.ID,
		UserName: user.UserName,
		RealName: user.RealName,
	}

	userRole, _, err := a.userRoleRepo.Query(ctx, request.UserRoleQuery{
		UserID: userID,
	})
	if err != nil {
		return nil, err
	}
	userRoleRes := new(response.UserRoles)
	structure.Copy(userRole, userRoleRes)
	if roleIDs := userRoleRes.ToRoleIDs(); len(roleIDs) > 0 {
		roles, _, err := a.roleRepo.Query(ctx, request.RoleQuery{
			IDs:    roleIDs,
			Status: 1,
		})
		if err != nil {
			return nil, err
		}
		rolesRes := new(response.Roles)
		structure.Copy(roles, rolesRes)
		info.Roles = *rolesRes
	}

	return info, nil
}

func (a *login) QueryUserMenuTree(ctx context.Context, userID string) (response.MenuTrees, error) {
	isRoot := checkIsRootUser(userID)
	// show all menu when root user
	if isRoot {
		menus, _, err := a.menuRepo.Query(ctx, request.MenuQuery{
			OrderFields: request.NewOrderFields(request.NewOrderField("sequence", request.OrderByDESC)),
			Status:      1,
		})
		if err != nil {
			return nil, err
		}
		menuAction, _, err := a.menuActionRepo.Query(ctx, request.MenuActionQuery{})
		if err != nil {
			return nil, err
		}
		menuRes := new(response.Menus)
		structure.Copy(menus, menuRes)
		menuActionsRes := new(response.MenuActions)
		structure.Copy(menuAction, menuActionsRes)
		return menuRes.FillMenuAction(menuActionsRes.ToMenuIDMap()).ToTree(), nil
	}

	userRoles, _, err := a.userRoleRepo.Query(ctx, request.UserRoleQuery{
		UserID: userID,
	})
	if err != nil {
		return nil, err
	}
	if len(userRoles) == 0 {
		return nil, errors2.ErrNoPerm
	}
	userRolesRes := new(response.UserRoles)
	structure.Copy(userRoles, userRolesRes)
	roleMenus, _, err := a.roleMenuRepo.Query(ctx, request.RoleMenuQuery{
		RoleIDs: userRolesRes.ToRoleIDs(),
	})
	if err != nil {
		return nil, err
	}
	if len(roleMenus) == 0 {
		return nil, errors2.ErrNoPerm
	}
	roleMenusRes := new(response.RoleMenus)
	structure.Copy(roleMenus, roleMenusRes)
	menus, _, err := a.menuRepo.Query(ctx, request.MenuQuery{
		IDs:    roleMenusRes.ToMenuIDs(),
		Status: 1,
	})
	if err != nil {
		return nil, err
	}
	if len(menus) == 0 {
		return nil, errors2.ErrNoPerm
	}
	menusRes := response.Menus{}
	structure.Copy(menus, menusRes)
	mData := menusRes.ToMap()
	var qIDs []string
	for _, pid := range menusRes.SplitParentIDs() {
		if _, ok := mData[pid]; !ok {
			qIDs = append(qIDs, pid)
		}
	}
	if len(qIDs) > 0 {
		pmenu, _, err := a.menuRepo.Query(ctx, request.MenuQuery{
			IDs: qIDs,
		})
		if err != nil {
			return nil, err
		}
		pmenuRes := response.Menus{}
		structure.Copy(pmenu, pmenuRes)
		menusRes = append(menusRes, pmenuRes...)
	}

	sort.Sort(menusRes)
	menuActions, _, err := a.menuActionRepo.Query(ctx, request.MenuActionQuery{
		IDs: roleMenusRes.ToActionIDs(),
	})
	if err != nil {
		return nil, err
	}
	menuActionsRes := new(response.MenuActions)
	structure.Copy(menuActions, menuActionsRes)
	return menusRes.FillMenuAction(menuActionsRes.ToMenuIDMap()).ToTree(), nil
}

func (a *login) UpdatePassword(ctx context.Context, userID string, params response.UpdatePassword) error {
	if checkIsRootUser(userID) {
		return errors2.New400Response("The root user is not allowed to update the password")
	}

	user, err := a.checkAndGetUser(ctx, userID)
	if err != nil {
		return err
	} else if hash.SHA1String(params.OldPassword) != user.Password {
		return errors2.New400Response("The old password is invalid")
	}

	params.NewPassword = hash.SHA1String(params.NewPassword)
	return a.userRepo.UpdatePassword(ctx, userID, params.NewPassword)
}
