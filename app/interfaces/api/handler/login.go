package handler

import (
	"github.com/LyricTian/captcha"
	"github.com/gin-gonic/gin"
	"github.com/linzhengen/ddd-gin-admin/app/application"
	"github.com/linzhengen/ddd-gin-admin/app/domain/errors"
	"github.com/linzhengen/ddd-gin-admin/app/domain/valueobject/schema"
	"github.com/linzhengen/ddd-gin-admin/app/interfaces/api"
	"github.com/linzhengen/ddd-gin-admin/configs"
	"github.com/linzhengen/ddd-gin-admin/pkg/logger"
)

type Login interface {
	GetCaptcha(c *gin.Context)
	ResCaptcha(c *gin.Context)
	Login(c *gin.Context)
	Logout(c *gin.Context)
	RefreshToken(c *gin.Context)
	GetUserInfo(c *gin.Context)
	QueryUserMenuTree(c *gin.Context)
	UpdatePassword(c *gin.Context)
}

func NewLogin(loginApp application.Login) Login {
	return &login{
		loginApp: loginApp,
	}
}

type login struct {
	loginApp application.Login
}

func (a *login) GetCaptcha(c *gin.Context) {
	ctx := c.Request.Context()
	item, err := a.loginApp.GetCaptcha(ctx, configs.C.Captcha.Length)
	if err != nil {
		api.ResError(c, err)
		return
	}
	api.ResSuccess(c, item)
}

func (a *login) ResCaptcha(c *gin.Context) {
	ctx := c.Request.Context()
	captchaID := c.Query("id")
	if captchaID == "" {
		api.ResError(c, errors.New400Response("Captcha is required"))
		return
	}

	if c.Query("reload") != "" {
		if !captcha.Reload(captchaID) {
			api.ResError(c, errors.New400Response("Failed get Captcha"))
			return
		}
	}

	cfg := configs.C.Captcha
	err := a.loginApp.ResCaptcha(ctx, c.Writer, captchaID, cfg.Width, cfg.Height)
	if err != nil {
		api.ResError(c, err)
	}
}

func (a *login) Login(c *gin.Context) {
	ctx := c.Request.Context()
	var item schema.LoginParam
	if err := api.ParseJSON(c, &item); err != nil {
		api.ResError(c, err)
		return
	}

	if !captcha.VerifyString(item.CaptchaID, item.CaptchaCode) {
		api.ResError(c, errors.New400Response("Invalid Captcha"))
		return
	}

	user, err := a.loginApp.Verify(ctx, item.UserName, item.Password)
	if err != nil {
		api.ResError(c, err)
		return
	}

	userID := user.ID
	api.SetUserID(c, userID)

	tokenInfo, err := a.loginApp.GenerateToken(ctx, userID)
	if err != nil {
		api.ResError(c, err)
		return
	}

	ctx = logger.NewUserIDContext(ctx, userID)
	ctx = logger.NewTagContext(ctx, "__login__")
	logger.WithContext(ctx).Infof("logged in")
	api.ResSuccess(c, tokenInfo)
}

func (a *login) Logout(c *gin.Context) {
	ctx := c.Request.Context()

	userID := api.GetUserID(c)
	if userID != "" {
		ctx = logger.NewTagContext(ctx, "__logout__")
		err := a.loginApp.DestroyToken(ctx, api.GetToken(c))
		if err != nil {
			logger.WithContext(ctx).Errorf(err.Error())
		}
		logger.WithContext(ctx).Infof("lougged out")
	}
	api.ResOK(c)
}

func (a *login) RefreshToken(c *gin.Context) {
	ctx := c.Request.Context()
	tokenInfo, err := a.loginApp.GenerateToken(ctx, api.GetUserID(c))
	if err != nil {
		api.ResError(c, err)
		return
	}
	api.ResSuccess(c, tokenInfo)
}

func (a *login) GetUserInfo(c *gin.Context) {
	ctx := c.Request.Context()
	info, err := a.loginApp.GetLoginInfo(ctx, api.GetUserID(c))
	if err != nil {
		api.ResError(c, err)
		return
	}
	api.ResSuccess(c, info)
}

func (a *login) QueryUserMenuTree(c *gin.Context) {
	ctx := c.Request.Context()
	menus, err := a.loginApp.QueryUserMenuTree(ctx, api.GetUserID(c))
	if err != nil {
		api.ResError(c, err)
		return
	}
	api.ResList(c, menus)
}

func (a *login) UpdatePassword(c *gin.Context) {
	ctx := c.Request.Context()
	var item schema.UpdatePasswordParam
	if err := api.ParseJSON(c, &item); err != nil {
		api.ResError(c, err)
		return
	}

	err := a.loginApp.UpdatePassword(ctx, api.GetUserID(c), item)
	if err != nil {
		api.ResError(c, err)
		return
	}
	api.ResOK(c)
}
