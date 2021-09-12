package application

import (
	"context"
	"net/http"
	"sort"

	"github.com/linzhengen/ddd-gin-admin/app/domain/errors"

	"github.com/linzhengen/ddd-gin-admin/app/domain/repository"
	"github.com/linzhengen/ddd-gin-admin/app/domain/schema"

	"github.com/LyricTian/captcha"
	"github.com/linzhengen/ddd-gin-admin/pkg/auth"
	"github.com/linzhengen/ddd-gin-admin/pkg/util/hash"
)

type Login interface {
	GetCaptcha(ctx context.Context, length int) (*schema.LoginCaptcha, error)
	ResCaptcha(ctx context.Context, w http.ResponseWriter, captchaID string, width, height int) error
	Verify(ctx context.Context, userName, password string) (*schema.User, error)
	GenerateToken(ctx context.Context, userID string) (*schema.LoginTokenInfo, error)
	DestroyToken(ctx context.Context, tokenString string) error
	GetLoginInfo(ctx context.Context, userID string) (*schema.UserLoginInfo, error)
	QueryUserMenuTree(ctx context.Context, userID string) (schema.MenuTrees, error)
	UpdatePassword(ctx context.Context, userID string, params schema.UpdatePasswordParam) error
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

func (a *login) GetCaptcha(ctx context.Context, length int) (*schema.LoginCaptcha, error) {
	captchaID := captcha.NewLen(length)
	item := &schema.LoginCaptcha{
		CaptchaID: captchaID,
	}
	return item, nil
}

func (a *login) ResCaptcha(ctx context.Context, w http.ResponseWriter, captchaID string, width, height int) error {
	err := captcha.WriteImage(w, captchaID, width, height)
	if err != nil {
		if err == captcha.ErrNotFound {
			return errors.ErrNotFound
		}
		return errors.WithStack(err)
	}

	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")
	w.Header().Set("Content-Type", "image/png")
	return nil
}

func (a *login) Verify(ctx context.Context, userName, password string) (*schema.User, error) {
	// is root user
	root := schema.GetRootUser()
	if userName == root.UserName && root.Password == password {
		return root, nil
	}

	result, err := a.userRepo.Query(ctx, schema.UserQueryParam{
		UserName: userName,
	})
	if err != nil {
		return nil, err
	}
	if len(result.Data) == 0 {
		return nil, errors.ErrInvalidUserName
	}

	item := result.Data[0]
	if item.Password != hash.SHA1String(password) {
		return nil, errors.ErrInvalidPassword
	}
	if item.Status != 1 {
		return nil, errors.ErrUserDisable
	}

	return item, nil
}

func (a *login) GenerateToken(ctx context.Context, userID string) (*schema.LoginTokenInfo, error) {
	tokenInfo, err := a.auth.GenerateToken(ctx, userID)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	item := &schema.LoginTokenInfo{
		AccessToken: tokenInfo.GetAccessToken(),
		TokenType:   tokenInfo.GetTokenType(),
		ExpiresAt:   tokenInfo.GetExpiresAt(),
	}
	return item, nil
}

func (a *login) DestroyToken(ctx context.Context, tokenString string) error {
	err := a.auth.DestroyToken(ctx, tokenString)
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (a *login) checkAndGetUser(ctx context.Context, userID string) (*schema.User, error) {
	user, err := a.userRepo.Get(ctx, userID)
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

func (a *login) GetLoginInfo(ctx context.Context, userID string) (*schema.UserLoginInfo, error) {
	if isRoot := schema.CheckIsRootUser(ctx, userID); isRoot {
		root := schema.GetRootUser()
		loginInfo := &schema.UserLoginInfo{
			UserName: root.UserName,
			RealName: root.RealName,
		}
		return loginInfo, nil
	}

	user, err := a.checkAndGetUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	info := &schema.UserLoginInfo{
		UserID:   user.ID,
		UserName: user.UserName,
		RealName: user.RealName,
	}

	userRoleResult, err := a.userRoleRepo.Query(ctx, schema.UserRoleQueryParam{
		UserID: userID,
	})
	if err != nil {
		return nil, err
	}

	if roleIDs := userRoleResult.Data.ToRoleIDs(); len(roleIDs) > 0 {
		roleResult, err := a.roleRepo.Query(ctx, schema.RoleQueryParam{
			IDs:    roleIDs,
			Status: 1,
		})
		if err != nil {
			return nil, err
		}
		info.Roles = roleResult.Data
	}

	return info, nil
}

func (a *login) QueryUserMenuTree(ctx context.Context, userID string) (schema.MenuTrees, error) {
	isRoot := schema.CheckIsRootUser(ctx, userID)
	// show all menu when root user
	if isRoot {
		result, err := a.menuRepo.Query(ctx, schema.MenuQueryParam{
			Status: 1,
		}, schema.MenuQueryOptions{
			OrderFields: schema.NewOrderFields(schema.NewOrderField("sequence", schema.OrderByDESC)),
		})
		if err != nil {
			return nil, err
		}

		menuActionResult, err := a.menuActionRepo.Query(ctx, schema.MenuActionQueryParam{})
		if err != nil {
			return nil, err
		}
		return result.Data.FillMenuAction(menuActionResult.Data.ToMenuIDMap()).ToTree(), nil
	}

	userRoleResult, err := a.userRoleRepo.Query(ctx, schema.UserRoleQueryParam{
		UserID: userID,
	})
	if err != nil {
		return nil, err
	}
	if len(userRoleResult.Data) == 0 {
		return nil, errors.ErrNoPerm
	}

	roleMenuResult, err := a.roleMenuRepo.Query(ctx, schema.RoleMenuQueryParam{
		RoleIDs: userRoleResult.Data.ToRoleIDs(),
	})
	if err != nil {
		return nil, err
	}
	if len(roleMenuResult.Data) == 0 {
		return nil, errors.ErrNoPerm
	}

	menuResult, err := a.menuRepo.Query(ctx, schema.MenuQueryParam{
		IDs:    roleMenuResult.Data.ToMenuIDs(),
		Status: 1,
	})
	if err != nil {
		return nil, err
	}
	if len(menuResult.Data) == 0 {
		return nil, errors.ErrNoPerm
	}

	mData := menuResult.Data.ToMap()
	var qIDs []string
	for _, pid := range menuResult.Data.SplitParentIDs() {
		if _, ok := mData[pid]; !ok {
			qIDs = append(qIDs, pid)
		}
	}
	if len(qIDs) > 0 {
		pmenuResult, err := a.menuRepo.Query(ctx, schema.MenuQueryParam{
			IDs: qIDs,
		})
		if err != nil {
			return nil, err
		}
		menuResult.Data = append(menuResult.Data, pmenuResult.Data...)
	}

	sort.Sort(menuResult.Data)
	menuActionResult, err := a.menuActionRepo.Query(ctx, schema.MenuActionQueryParam{
		IDs: roleMenuResult.Data.ToActionIDs(),
	})
	if err != nil {
		return nil, err
	}
	return menuResult.Data.FillMenuAction(menuActionResult.Data.ToMenuIDMap()).ToTree(), nil
}

func (a *login) UpdatePassword(ctx context.Context, userID string, params schema.UpdatePasswordParam) error {
	if schema.CheckIsRootUser(ctx, userID) {
		return errors.New400Response("The root user is not allowed to update the password")
	}

	user, err := a.checkAndGetUser(ctx, userID)
	if err != nil {
		return err
	} else if hash.SHA1String(params.OldPassword) != user.Password {
		return errors.New400Response("The old password is invalid")
	}

	params.NewPassword = hash.SHA1String(params.NewPassword)
	return a.userRepo.UpdatePassword(ctx, userID, params.NewPassword)
}
