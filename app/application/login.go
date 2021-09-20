package application

import (
	"context"
	"net/http"

	"github.com/linzhengen/ddd-gin-admin/app/domain/service"
	"github.com/linzhengen/ddd-gin-admin/app/domain/valueobject/schema"
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
	loginSvc service.Login,
) Login {
	return &login{
		loginSvc: loginSvc,
	}
}

type login struct {
	loginSvc service.Login
}

func (l login) GetCaptcha(ctx context.Context, length int) (*schema.LoginCaptcha, error) {
	return l.loginSvc.GetCaptcha(ctx, length)
}

func (l login) ResCaptcha(ctx context.Context, w http.ResponseWriter, captchaID string, width, height int) error {
	return l.loginSvc.ResCaptcha(ctx, w, captchaID, width, height)
}

func (l login) Verify(ctx context.Context, userName, password string) (*schema.User, error) {
	return l.loginSvc.Verify(ctx, userName, password)
}

func (l login) GenerateToken(ctx context.Context, userID string) (*schema.LoginTokenInfo, error) {
	return l.loginSvc.GenerateToken(ctx, userID)
}

func (l login) DestroyToken(ctx context.Context, tokenString string) error {
	return l.loginSvc.DestroyToken(ctx, tokenString)
}

func (l login) GetLoginInfo(ctx context.Context, userID string) (*schema.UserLoginInfo, error) {
	return l.loginSvc.GetLoginInfo(ctx, userID)
}

func (l login) QueryUserMenuTree(ctx context.Context, userID string) (schema.MenuTrees, error) {
	return l.loginSvc.QueryUserMenuTree(ctx, userID)
}

func (l login) UpdatePassword(ctx context.Context, userID string, params schema.UpdatePasswordParam) error {
	return l.loginSvc.UpdatePassword(ctx, userID, params)
}
