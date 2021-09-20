package service

import (
	"context"
	"net/http"
	"sort"

	"github.com/linzhengen/ddd-gin-admin/app/domain/factory"

	"github.com/linzhengen/ddd-gin-admin/app/domain/valueobject/errors"

	"github.com/linzhengen/ddd-gin-admin/app/domain/repository"
	"github.com/linzhengen/ddd-gin-admin/app/domain/valueobject/schema"

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
	userFactory factory.User,
	userRoleFactory factory.UserRole,
	roleFactory factory.Role,
	roleMenuFactory factory.RoleMenu,
	menuActionFactory factory.MenuAction,
) Login {
	return &login{
		auth:              auth,
		userRepo:          userRepo,
		userRoleRepo:      userRoleRepo,
		roleRepo:          roleRepo,
		roleMenuRepo:      roleMenuRepo,
		menuRepo:          menuRepo,
		menuActionRepo:    menuActionRepo,
		userFactory:       userFactory,
		userRoleFactory:   userRoleFactory,
		roleFactory:       roleFactory,
		roleMenuFactory:   roleMenuFactory,
		menuActionFactory: menuActionFactory,
	}
}

type login struct {
	auth              auth.Author
	userRepo          repository.UserRepository
	userRoleRepo      repository.UserRoleRepository
	roleRepo          repository.RoleRepository
	roleMenuRepo      repository.RoleMenuRepository
	menuRepo          repository.MenuRepository
	menuActionRepo    repository.MenuActionRepository
	userFactory       factory.User
	userRoleFactory   factory.UserRole
	roleFactory       factory.Role
	roleMenuFactory   factory.RoleMenu
	menuActionFactory factory.MenuAction
	menuFactory       factory.Menu
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

	result, _, err := a.userRepo.Query(ctx, schema.UserQueryParam{
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

	return a.userFactory.ToSchema(item), nil
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
	return a.userFactory.ToSchema(user), nil
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

	userRoleResult, _, err := a.userRoleRepo.Query(ctx, schema.UserRoleQueryParam{
		UserID: userID,
	})
	if err != nil {
		return nil, err
	}

	if roleIDs := a.userRoleFactory.ToSchemaList(userRoleResult).ToRoleIDs(); len(roleIDs) > 0 {
		roleResult, _, err := a.roleRepo.Query(ctx, schema.RoleQueryParam{
			IDs:    roleIDs,
			Status: 1,
		})
		if err != nil {
			return nil, err
		}
		info.Roles = a.roleFactory.ToSchemaList(roleResult)
	}

	return info, nil
}

func (a *login) QueryUserMenuTree(ctx context.Context, userID string) (schema.MenuTrees, error) {
	isRoot := schema.CheckIsRootUser(ctx, userID)
	// show all menu when root user
	if isRoot {
		result, _, err := a.menuRepo.Query(ctx, schema.MenuQueryParam{
			OrderFields: schema.NewOrderFields(schema.NewOrderField("sequence", schema.OrderByDESC)),
			Status:      1,
		})
		if err != nil {
			return nil, err
		}

		menuActionResult, _, err := a.menuActionRepo.Query(ctx, schema.MenuActionQueryParam{})
		if err != nil {
			return nil, err
		}
		return a.menuFactory.ToSchemaList(result).FillMenuAction(
			a.menuActionFactory.ToSchemaList(menuActionResult).ToMenuIDMap(),
		).ToTree(), nil
	}

	userRoleResult, _, err := a.userRoleRepo.Query(ctx, schema.UserRoleQueryParam{
		UserID: userID,
	})
	if err != nil {
		return nil, err
	}
	if len(userRoleResult) == 0 {
		return nil, errors.ErrNoPerm
	}

	roleMenuResult, _, err := a.roleMenuRepo.Query(ctx, schema.RoleMenuQueryParam{
		RoleIDs: a.userRoleFactory.ToSchemaList(userRoleResult).ToRoleIDs(),
	})
	if err != nil {
		return nil, err
	}
	if len(roleMenuResult) == 0 {
		return nil, errors.ErrNoPerm
	}

	menuResult, _, err := a.menuRepo.Query(ctx, schema.MenuQueryParam{
		IDs:    a.roleMenuFactory.ToSchemaList(roleMenuResult).ToMenuIDs(),
		Status: 1,
	})
	if err != nil {
		return nil, err
	}
	if len(menuResult) == 0 {
		return nil, errors.ErrNoPerm
	}

	menusSchema := a.menuFactory.ToSchemaList(menuResult)
	mData := menusSchema.ToMap()
	var qIDs []string
	for _, pid := range menusSchema.SplitParentIDs() {
		if _, ok := mData[pid]; !ok {
			qIDs = append(qIDs, pid)
		}
	}
	if len(qIDs) > 0 {
		pmenuResult, _, err := a.menuRepo.Query(ctx, schema.MenuQueryParam{
			IDs: qIDs,
		})
		if err != nil {
			return nil, err
		}
		menusSchema = append(menusSchema, a.menuFactory.ToSchemaList(pmenuResult)...)
	}

	sort.Sort(menusSchema)
	menuActionResult, _, err := a.menuActionRepo.Query(ctx, schema.MenuActionQueryParam{
		IDs: a.roleMenuFactory.ToSchemaList(roleMenuResult).ToActionIDs(),
	})
	if err != nil {
		return nil, err
	}
	return menusSchema.FillMenuAction(a.menuActionFactory.ToSchemaList(menuActionResult).ToMenuIDMap()).ToTree(), nil
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
