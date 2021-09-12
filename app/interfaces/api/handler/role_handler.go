package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/linzhengen/ddd-gin-admin/app/application"
	"github.com/linzhengen/ddd-gin-admin/app/domain/schema"
	"github.com/linzhengen/ddd-gin-admin/app/interfaces/api"
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
	return &role{roleApp: roleApp}
}

type role struct {
	roleApp application.Role
}

func (a *role) Query(c *gin.Context) {
	ctx := c.Request.Context()
	var params schema.RoleQueryParam
	if err := api.ParseQuery(c, &params); err != nil {
		api.ResError(c, err)
		return
	}

	params.Pagination = true
	result, err := a.roleApp.Query(ctx, params, schema.RoleQueryOptions{
		OrderFields: schema.NewOrderFields(schema.NewOrderField("sequence", schema.OrderByDESC)),
	})
	if err != nil {
		api.ResError(c, err)
		return
	}
	api.ResPage(c, result.Data, result.PageResult)
}

func (a *role) QuerySelect(c *gin.Context) {
	ctx := c.Request.Context()
	var params schema.RoleQueryParam
	if err := api.ParseQuery(c, &params); err != nil {
		api.ResError(c, err)
		return
	}

	result, err := a.roleApp.Query(ctx, params, schema.RoleQueryOptions{
		OrderFields: schema.NewOrderFields(schema.NewOrderField("sequence", schema.OrderByDESC)),
	})
	if err != nil {
		api.ResError(c, err)
		return
	}
	api.ResList(c, result.Data)
}

func (a *role) Get(c *gin.Context) {
	ctx := c.Request.Context()
	item, err := a.roleApp.Get(ctx, c.Param("id"))
	if err != nil {
		api.ResError(c, err)
		return
	}
	api.ResSuccess(c, item)
}

func (a *role) Create(c *gin.Context) {
	ctx := c.Request.Context()
	var item schema.Role
	if err := api.ParseJSON(c, &item); err != nil {
		api.ResError(c, err)
		return
	}

	item.Creator = api.GetUserID(c)
	result, err := a.roleApp.Create(ctx, item)
	if err != nil {
		api.ResError(c, err)
		return
	}
	api.ResSuccess(c, result)
}

func (a *role) Update(c *gin.Context) {
	ctx := c.Request.Context()
	var item schema.Role
	if err := api.ParseJSON(c, &item); err != nil {
		api.ResError(c, err)
		return
	}

	err := a.roleApp.Update(ctx, c.Param("id"), item)
	if err != nil {
		api.ResError(c, err)
		return
	}
	api.ResOK(c)
}

func (a *role) Delete(c *gin.Context) {
	ctx := c.Request.Context()
	err := a.roleApp.Delete(ctx, c.Param("id"))
	if err != nil {
		api.ResError(c, err)
		return
	}
	api.ResOK(c)
}

func (a *role) Enable(c *gin.Context) {
	ctx := c.Request.Context()
	err := a.roleApp.UpdateStatus(ctx, c.Param("id"), 1)
	if err != nil {
		api.ResError(c, err)
		return
	}
	api.ResOK(c)
}

func (a *role) Disable(c *gin.Context) {
	ctx := c.Request.Context()
	err := a.roleApp.UpdateStatus(ctx, c.Param("id"), 2)
	if err != nil {
		api.ResError(c, err)
		return
	}
	api.ResOK(c)
}
