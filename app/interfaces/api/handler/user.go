package handler

import (
	"strings"

	"github.com/linzhengen/ddd-gin-admin/app/application"
	"github.com/linzhengen/ddd-gin-admin/app/interfaces/api"

	"github.com/linzhengen/ddd-gin-admin/app/domain/valueobject/errors"

	"github.com/gin-gonic/gin"
	"github.com/linzhengen/ddd-gin-admin/app/domain/valueobject/schema"
)

type User interface {
	Query(c *gin.Context)
	Get(c *gin.Context)
	Create(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
	Enable(c *gin.Context)
	Disable(c *gin.Context)
}

func NewUser(userApp application.User) User {
	return &user{userApp: userApp}
}

type user struct {
	userApp application.User
}

func (a *user) Query(c *gin.Context) {
	ctx := c.Request.Context()
	var params schema.UserQueryParam
	if err := api.ParseQuery(c, &params); err != nil {
		api.ResError(c, err)
		return
	}
	if v := c.Query("roleIDs"); v != "" {
		params.RoleIDs = strings.Split(v, ",")
	}

	params.Pagination = true
	result, err := a.userApp.QueryShow(ctx, params)
	if err != nil {
		api.ResError(c, err)
		return
	}
	api.ResPage(c, result.Data, result.PageResult)
}

func (a *user) Get(c *gin.Context) {
	ctx := c.Request.Context()
	item, err := a.userApp.Get(ctx, c.Param("id"))
	if err != nil {
		api.ResError(c, err)
		return
	}
	api.ResSuccess(c, item.CleanSecure())
}

func (a *user) Create(c *gin.Context) {
	ctx := c.Request.Context()
	var item schema.User
	if err := api.ParseJSON(c, &item); err != nil {
		api.ResError(c, err)
		return
	} else if item.Password == "" {
		api.ResError(c, errors.New400Response("密码不能为空"))
		return
	}

	item.Creator = api.GetUserID(c)
	result, err := a.userApp.Create(ctx, item)
	if err != nil {
		api.ResError(c, err)
		return
	}
	api.ResSuccess(c, result)
}

func (a *user) Update(c *gin.Context) {
	ctx := c.Request.Context()
	var item schema.User
	if err := api.ParseJSON(c, &item); err != nil {
		api.ResError(c, err)
		return
	}

	err := a.userApp.Update(ctx, c.Param("id"), item)
	if err != nil {
		api.ResError(c, err)
		return
	}
	api.ResOK(c)
}

func (a *user) Delete(c *gin.Context) {
	ctx := c.Request.Context()
	err := a.userApp.Delete(ctx, c.Param("id"))
	if err != nil {
		api.ResError(c, err)
		return
	}
	api.ResOK(c)
}

func (a *user) Enable(c *gin.Context) {
	ctx := c.Request.Context()
	err := a.userApp.UpdateStatus(ctx, c.Param("id"), 1)
	if err != nil {
		api.ResError(c, err)
		return
	}
	api.ResOK(c)
}

func (a *user) Disable(c *gin.Context) {
	ctx := c.Request.Context()
	err := a.userApp.UpdateStatus(ctx, c.Param("id"), 2)
	if err != nil {
		api.ResError(c, err)
		return
	}
	api.ResOK(c)
}
