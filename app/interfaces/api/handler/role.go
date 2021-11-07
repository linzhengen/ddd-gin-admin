package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/linzhengen/ddd-gin-admin/app/application"
	"github.com/linzhengen/ddd-gin-admin/app/domain/pagination"
	"github.com/linzhengen/ddd-gin-admin/app/domain/user/role"
	"github.com/linzhengen/ddd-gin-admin/app/interfaces/api"
	"github.com/linzhengen/ddd-gin-admin/app/interfaces/api/request"
	"github.com/linzhengen/ddd-gin-admin/app/interfaces/api/response"
)

type Role interface {
	Query(c *gin.Context)
	QuerySelect(c *gin.Context)
	Get(c *gin.Context)
	Create(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
	Enable(c *gin.Context)
	Disable(c *gin.Context)
}

func NewRole(roleApp application.Role) Role {
	return &roleHandler{roleApp: roleApp}
}

type roleHandler struct {
	roleApp application.Role
}

func (a *roleHandler) Query(c *gin.Context) {
	ctx := c.Request.Context()
	var params request.RoleQueryParam
	if err := api.ParseQuery(c, &params); err != nil {
		api.ResError(c, err)
		return
	}

	domainParams := role.QueryParam{
		PaginationParam: pagination.Param{Pagination: true},
		OrderFields:     pagination.NewOrderFields(pagination.NewOrderField("sequence", pagination.OrderByDESC)),
		IDs:             params.IDs,
		Name:            params.Name,
		QueryValue:      params.QueryValue,
		UserID:          params.UserID,
		Status:          params.Status,
	}
	result, p, err := a.roleApp.Query(ctx, domainParams)
	if err != nil {
		api.ResError(c, err)
		return
	}
	api.ResPage(c, response.RolesFromDomain(result), p)
}

func (a *roleHandler) QuerySelect(c *gin.Context) {
	ctx := c.Request.Context()
	var params request.RoleQueryParam
	if err := api.ParseQuery(c, &params); err != nil {
		api.ResError(c, err)
		return
	}

	domainParams := role.QueryParam{
		PaginationParam: pagination.Param{Pagination: true},
		OrderFields:     pagination.NewOrderFields(pagination.NewOrderField("sequence", pagination.OrderByDESC)),
		IDs:             params.IDs,
		Name:            params.Name,
		QueryValue:      params.QueryValue,
		UserID:          params.UserID,
		Status:          params.Status,
	}
	result, _, err := a.roleApp.Query(ctx, domainParams)
	if err != nil {
		api.ResError(c, err)
		return
	}
	api.ResList(c, response.RolesFromDomain(result))
}

func (a *roleHandler) Get(c *gin.Context) {
	ctx := c.Request.Context()
	item, err := a.roleApp.Get(ctx, c.Param("id"))
	if err != nil {
		api.ResError(c, err)
		return
	}
	api.ResSuccess(c, response.RoleFromDomain(item))
}

func (a *roleHandler) Create(c *gin.Context) {
	ctx := c.Request.Context()
	var item request.Role
	if err := api.ParseJSON(c, &item); err != nil {
		api.ResError(c, err)
		return
	}

	item.Creator = api.GetUserID(c)
	result, err := a.roleApp.Create(ctx, item.ToDomain())
	if err != nil {
		api.ResError(c, err)
		return
	}
	api.ResSuccess(c, response.NewIDResult(result))
}

func (a *roleHandler) Update(c *gin.Context) {
	ctx := c.Request.Context()
	var item request.Role
	if err := api.ParseJSON(c, &item); err != nil {
		api.ResError(c, err)
		return
	}

	err := a.roleApp.Update(ctx, c.Param("id"), item.ToDomain())
	if err != nil {
		api.ResError(c, err)
		return
	}
	api.ResOK(c)
}

func (a *roleHandler) Delete(c *gin.Context) {
	ctx := c.Request.Context()
	err := a.roleApp.Delete(ctx, c.Param("id"))
	if err != nil {
		api.ResError(c, err)
		return
	}
	api.ResOK(c)
}

func (a *roleHandler) Enable(c *gin.Context) {
	ctx := c.Request.Context()
	err := a.roleApp.UpdateStatus(ctx, c.Param("id"), 1)
	if err != nil {
		api.ResError(c, err)
		return
	}
	api.ResOK(c)
}

func (a *roleHandler) Disable(c *gin.Context) {
	ctx := c.Request.Context()
	err := a.roleApp.UpdateStatus(ctx, c.Param("id"), 2)
	if err != nil {
		api.ResError(c, err)
		return
	}
	api.ResOK(c)
}
