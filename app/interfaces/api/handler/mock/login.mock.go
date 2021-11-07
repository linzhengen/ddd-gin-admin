package mock

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

var LoginSet = wire.NewSet(wire.Struct(new(Login), "*"))

type Login struct {
}

// GetCaptcha Get Captcha
// @Tags Login
// @Summary Get Captcha
// @Success 200 {object} response.LoginCaptcha
// @Router /api/v1/pub/login/captchaid [get]
func (a *Login) GetCaptcha(c *gin.Context) {
}

// ResCaptcha Response captcha
// @Tags Login
// @Summary Response captcha
// @Param id query string true "Captcha ID"
// @Param reload query string false "Reload"
// @Produce image/png
// @Success 200 "ResCaptcha"
// @Failure 400 {object} response.ErrorResult "{error:{code:0,message:BadRequest}}"
// @Failure 500 {object} response.ErrorResult "{error:{code:0,message:SystemError}}"
// @Router /api/v1/pub/login/captcha [get]
func (a *Login) ResCaptcha(c *gin.Context) {
}

// Login Login
// @Tags Login
// @Summary Login
// @Param body body request.LoginParam true "Request parameters"
// @Success 200 {object} response.LoginTokenInfo
// @Failure 400 {object} response.ErrorResult "{error:{code:0,message:BadRequest}}"
// @Failure 500 {object} response.ErrorResult "{error:{code:0,message:SystemError}}"
// @Router /api/v1/pub/login [post]
func (a *Login) Login(c *gin.Context) {
}

// Logout Logout
// @Tags Login
// @Summary Logout
// @Success 200 {object} response.StatusResult "{status:OK}"
// @Router /api/v1/pub/login/exit [post]
func (a *Login) Logout(c *gin.Context) {
}

// RefreshToken Refresh token
// @Tags Login
// @Summary Refresh token
// @Security ApiKeyAuth
// @Success 200 {object} response.LoginTokenInfo
// @Failure 401 {object} response.ErrorResult "{error:{code:0,message:Unauthorized}}"
// @Failure 500 {object} response.ErrorResult "{error:{code:0,message:SystemError}}"
// @Router /api/v1/pub/refresh-token [post]
func (a *Login) RefreshToken(c *gin.Context) {
}

// GetUserInfo Get current user info
// @Tags Login
// @Summary Get current user info
// @Security ApiKeyAuth
// @Success 200 {object} response.UserLoginInfo
// @Failure 401 {object} response.ErrorResult "{error:{code:0,message:Unauthorized}}"
// @Failure 500 {object} response.ErrorResult "{error:{code:0,message:SystemError}}"
// @Router /api/v1/pub/current/user [get]
func (a *Login) GetUserInfo(c *gin.Context) {
}

// QueryUserMenuTree Get user menu tree
// @Tags Login
// @Summary Get user menu tree
// @Security ApiKeyAuth
// @Success 200 {object} response.ListResult{list=[]schema.MenuTree} "Search Result"
// @Failure 401 {object} response.ErrorResult "{error:{code:0,message:Unauthorized}}"
// @Failure 500 {object} response.ErrorResult "{error:{code:0,message:SystemError}}"
// @Router /api/v1/pub/current/menutree [get]
func (a *Login) QueryUserMenuTree(c *gin.Context) {
}

// UpdatePassword Update password
// @Tags Login
// @Summary Update password
// @Security ApiKeyAuth
// @Param body body request.UpdatePasswordParam true "Request parameters"
// @Success 200 {object} response.StatusResult "{status:OK}"
// @Failure 400 {object} response.ErrorResult "{error:{code:0,message:BadRequest}}"
// @Failure 401 {object} response.ErrorResult "{error:{code:0,message:Unauthorized}}"
// @Failure 500 {object} response.ErrorResult "{error:{code:0,message:SystemError}}"
// @Router /api/v1/pub/current/password [put]
func (a *Login) UpdatePassword(c *gin.Context) {
}
