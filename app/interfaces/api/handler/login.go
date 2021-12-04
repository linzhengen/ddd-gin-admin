package handler

import (
	"github.com/LyricTian/captcha"
	"github.com/gin-gonic/gin"
	"github.com/linzhengen/ddd-gin-admin/app/application"
	"github.com/linzhengen/ddd-gin-admin/app/domain/errors"
	"github.com/linzhengen/ddd-gin-admin/app/interfaces/api"
	"github.com/linzhengen/ddd-gin-admin/app/interfaces/api/request"
	"github.com/linzhengen/ddd-gin-admin/app/interfaces/api/response"
	"github.com/linzhengen/ddd-gin-admin/configs"
	"github.com/linzhengen/ddd-gin-admin/pkg/logger"
	"github.com/linzhengen/ddd-gin-admin/pkg/util/structure"
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
	captchaID := captcha.NewLen(configs.C.Captcha.Length)
	item := &response.LoginCaptcha{
		CaptchaID: captchaID,
	}
	api.ResSuccess(c, item)
}

func (a *login) ResCaptcha(c *gin.Context) {
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
	err := captcha.WriteImage(c.Writer, captchaID, cfg.Width, cfg.Height)
	if err != nil {
		if err == captcha.ErrNotFound {
			err = errors.ErrNotFound
		}
		api.ResError(c, err)
		return
	}
	c.Writer.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	c.Writer.Header().Set("Pragma", "no-cache")
	c.Writer.Header().Set("Expires", "0")
	c.Writer.Header().Set("Content-Type", "image/png")
}

func (a *login) Login(c *gin.Context) {
	ctx := c.Request.Context()
	var item request.LoginParam
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
	respTokenInfo := new(response.LoginTokenInfo)
	structure.Copy(tokenInfo, respTokenInfo)
	ctx = logger.NewUserIDContext(ctx, userID)
	ctx = logger.NewTagContext(ctx, "__login__")
	logger.WithContext(ctx).Infof("logged in")
	api.ResSuccess(c, respTokenInfo)
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
	schemaTokenInfo := new(response.LoginTokenInfo)
	structure.Copy(tokenInfo, schemaTokenInfo)
	api.ResSuccess(c, schemaTokenInfo)
}

func (a *login) GetUserInfo(c *gin.Context) {
	ctx := c.Request.Context()
	info, err := a.loginApp.GetLoginInfo(ctx, api.GetUserID(c))
	if err != nil {
		api.ResError(c, err)
		return
	}
	schemaInfo := &response.UserLoginInfo{
		UserID:   info.ID,
		UserName: info.UserName,
		RealName: info.RealName,
	}
	roles := make([]*response.Role, len(info.Roles))
	for i, item := range info.Roles {
		ts := new(response.Role)
		structure.Copy(item, ts)
		roles[i] = ts
	}
	schemaInfo.Roles = roles
	api.ResSuccess(c, schemaInfo)
}

func (a *login) QueryUserMenuTree(c *gin.Context) {
	ctx := c.Request.Context()
	menus, err := a.loginApp.QueryUserMenuTree(ctx, api.GetUserID(c))
	if err != nil {
		api.ResError(c, err)
		return
	}
	schemaMenus := make([]*response.MenuTrees, len(menus))
	for i, item := range menus {
		ts := new(response.MenuTrees)
		structure.Copy(item, ts)
		schemaMenus[i] = ts
	}
	api.ResList(c, schemaMenus)
}

func (a *login) UpdatePassword(c *gin.Context) {
	ctx := c.Request.Context()
	var item request.UpdatePasswordParam
	if err := api.ParseJSON(c, &item); err != nil {
		api.ResError(c, err)
		return
	}

	err := a.loginApp.UpdatePassword(ctx, api.GetUserID(c), item.OldPassword, item.NewPassword)
	if err != nil {
		api.ResError(c, err)
		return
	}
	api.ResOK(c)
}
