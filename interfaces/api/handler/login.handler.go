package handler

import (
	"github.com/LyricTian/captcha"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/linzhengen/ddd-gin-admin/application"
	"github.com/linzhengen/ddd-gin-admin/domain/schema"
	"github.com/linzhengen/ddd-gin-admin/infrastructure/config"
	"github.com/linzhengen/ddd-gin-admin/infrastructure/ginx"
	"github.com/linzhengen/ddd-gin-admin/pkg/errors"
	"github.com/linzhengen/ddd-gin-admin/pkg/logger"
)

var LoginSet = wire.NewSet(wire.Struct(new(Login), "*"))

type Login struct {
	LoginSrv *application.Login
}

func (a *Login) GetCaptcha(c *gin.Context) {
	ctx := c.Request.Context()
	item, err := a.LoginSrv.GetCaptcha(ctx, config.C.Captcha.Length)
	if err != nil {
		ginx.ResError(c, err)
		return
	}
	ginx.ResSuccess(c, item)
}

func (a *Login) ResCaptcha(c *gin.Context) {
	ctx := c.Request.Context()
	captchaID := c.Query("id")
	if captchaID == "" {
		ginx.ResError(c, errors.New400Response("Captcha is required"))
		return
	}

	if c.Query("reload") != "" {
		if !captcha.Reload(captchaID) {
			ginx.ResError(c, errors.New400Response("Failed get Captcha"))
			return
		}
	}

	cfg := config.C.Captcha
	err := a.LoginSrv.ResCaptcha(ctx, c.Writer, captchaID, cfg.Width, cfg.Height)
	if err != nil {
		ginx.ResError(c, err)
	}
}

func (a *Login) Login(c *gin.Context) {
	ctx := c.Request.Context()
	var item schema.LoginParam
	if err := ginx.ParseJSON(c, &item); err != nil {
		ginx.ResError(c, err)
		return
	}

	if !captcha.VerifyString(item.CaptchaID, item.CaptchaCode) {
		ginx.ResError(c, errors.New400Response("Invalid Captcha"))
		return
	}

	user, err := a.LoginSrv.Verify(ctx, item.UserName, item.Password)
	if err != nil {
		ginx.ResError(c, err)
		return
	}

	userID := user.ID
	ginx.SetUserID(c, userID)

	tokenInfo, err := a.LoginSrv.GenerateToken(ctx, userID)
	if err != nil {
		ginx.ResError(c, err)
		return
	}

	ctx = logger.NewUserIDContext(ctx, userID)
	ctx = logger.NewTagContext(ctx, "__login__")
	logger.WithContext(ctx).Infof("logged in")
	ginx.ResSuccess(c, tokenInfo)
}

func (a *Login) Logout(c *gin.Context) {
	ctx := c.Request.Context()

	userID := ginx.GetUserID(c)
	if userID != "" {
		ctx = logger.NewTagContext(ctx, "__logout__")
		err := a.LoginSrv.DestroyToken(ctx, ginx.GetToken(c))
		if err != nil {
			logger.WithContext(ctx).Errorf(err.Error())
		}
		logger.WithContext(ctx).Infof("lougged out")
	}
	ginx.ResOK(c)
}

func (a *Login) RefreshToken(c *gin.Context) {
	ctx := c.Request.Context()
	tokenInfo, err := a.LoginSrv.GenerateToken(ctx, ginx.GetUserID(c))
	if err != nil {
		ginx.ResError(c, err)
		return
	}
	ginx.ResSuccess(c, tokenInfo)
}

func (a *Login) GetUserInfo(c *gin.Context) {
	ctx := c.Request.Context()
	info, err := a.LoginSrv.GetLoginInfo(ctx, ginx.GetUserID(c))
	if err != nil {
		ginx.ResError(c, err)
		return
	}
	ginx.ResSuccess(c, info)
}

func (a *Login) QueryUserMenuTree(c *gin.Context) {
	ctx := c.Request.Context()
	menus, err := a.LoginSrv.QueryUserMenuTree(ctx, ginx.GetUserID(c))
	if err != nil {
		ginx.ResError(c, err)
		return
	}
	ginx.ResList(c, menus)
}

func (a *Login) UpdatePassword(c *gin.Context) {
	ctx := c.Request.Context()
	var item schema.UpdatePasswordParam
	if err := ginx.ParseJSON(c, &item); err != nil {
		ginx.ResError(c, err)
		return
	}

	err := a.LoginSrv.UpdatePassword(ctx, ginx.GetUserID(c), item)
	if err != nil {
		ginx.ResError(c, err)
		return
	}
	ginx.ResOK(c)
}
