package handler

import (
	"strings"

	"github.com/linzhengen/ddd-gin-admin/app/interfaces/api/response"

	"github.com/linzhengen/ddd-gin-admin/app/domain/pagination"
	"github.com/linzhengen/ddd-gin-admin/app/domain/user"

	"github.com/linzhengen/ddd-gin-admin/app/interfaces/api/request"

	"github.com/linzhengen/ddd-gin-admin/app/domain/errors"

	"github.com/linzhengen/ddd-gin-admin/app/application"
	"github.com/linzhengen/ddd-gin-admin/app/interfaces/api"

	"github.com/gin-gonic/gin"
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
	return &userHandler{userApp: userApp}
}

type userHandler struct {
	userApp application.User
}

func (a *userHandler) Query(c *gin.Context) {
	ctx := c.Request.Context()
	var params request.UserQueryParam
	if err := api.ParseQuery(c, &params); err != nil {
		api.ResError(c, err)
		return
	}
	if v := c.Query("roleIDs"); v != "" {
		params.RoleIDs = strings.Split(v, ",")
	}

	domainParams := user.QueryParams{
		PaginationParam: pagination.Param{Pagination: true},
		OrderFields:     nil,
		UserName:        params.UserName,
		QueryValue:      params.QueryValue,
		Status:          params.Status,
		RoleIDs:         params.RoleIDs,
	}
	result, p, err := a.userApp.QueryShow(ctx, domainParams)
	if err != nil {
		api.ResError(c, err)
		return
	}
	api.ResPage(c, response.UsersFromDomain(result), p)
}

func (a *userHandler) Get(c *gin.Context) {
	ctx := c.Request.Context()
	item, err := a.userApp.Get(ctx, c.Param("id"))
	if err != nil {
		api.ResError(c, err)
		return
	}
	api.ResSuccess(c, response.UserFromDomain(item).CleanSecure())
}

func (a *userHandler) Create(c *gin.Context) {
	ctx := c.Request.Context()
	var item request.User
	if err := api.ParseJSON(c, &item); err != nil {
		api.ResError(c, err)
		return
	} else if item.Password == "" {
		api.ResError(c, errors.New400Response("password is empty"))
		return
	}

	item.Creator = api.GetUserID(c)
	result, err := a.userApp.Create(ctx, item.ToDomain(), item.RoleIDs)
	if err != nil {
		api.ResError(c, err)
		return
	}
	api.ResSuccess(c, response.NewIDResult(result))
}

func (a *userHandler) Update(c *gin.Context) {
	ctx := c.Request.Context()
	var item request.User
	if err := api.ParseJSON(c, &item); err != nil {
		api.ResError(c, err)
		return
	}

	err := a.userApp.Update(ctx, c.Param("id"), item.ToDomain(), item.RoleIDs)
	if err != nil {
		api.ResError(c, err)
		return
	}
	api.ResOK(c)
}

func (a *userHandler) Delete(c *gin.Context) {
	ctx := c.Request.Context()
	err := a.userApp.Delete(ctx, c.Param("id"))
	if err != nil {
		api.ResError(c, err)
		return
	}
	api.ResOK(c)
}

func (a *userHandler) Enable(c *gin.Context) {
	ctx := c.Request.Context()
	err := a.userApp.UpdateStatus(ctx, c.Param("id"), 1)
	if err != nil {
		api.ResError(c, err)
		return
	}
	api.ResOK(c)
}

func (a *userHandler) Disable(c *gin.Context) {
	ctx := c.Request.Context()
	err := a.userApp.UpdateStatus(ctx, c.Param("id"), 2)
	if err != nil {
		api.ResError(c, err)
		return
	}
	api.ResOK(c)
}
