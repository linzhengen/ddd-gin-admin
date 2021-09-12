package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/linzhengen/ddd-gin-admin/app/application"
	"github.com/linzhengen/ddd-gin-admin/app/domain/schema"
	"github.com/linzhengen/ddd-gin-admin/app/interfaces/api"
)

type Menu interface {
	Query(c *gin.Context)
	QueryTree(c *gin.Context)
	Get(c *gin.Context)
	Create(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
	Enable(c *gin.Context)
	Disable(c *gin.Context)
}

func NewMenu(menuApp application.Menu) Menu {
	return &menu{
		menuApp: menuApp,
	}
}

type menu struct {
	menuApp application.Menu
}

func (a *menu) Query(c *gin.Context) {
	ctx := c.Request.Context()
	var params schema.MenuQueryParam
	if err := api.ParseQuery(c, &params); err != nil {
		api.ResError(c, err)
		return
	}

	params.Pagination = true
	result, err := a.menuApp.Query(ctx, params, schema.MenuQueryOptions{
		OrderFields: schema.NewOrderFields(schema.NewOrderField("sequence", schema.OrderByDESC)),
	})
	if err != nil {
		api.ResError(c, err)
		return
	}
	api.ResPage(c, result.Data, result.PageResult)
}

func (a *menu) QueryTree(c *gin.Context) {
	ctx := c.Request.Context()
	var params schema.MenuQueryParam
	if err := api.ParseQuery(c, &params); err != nil {
		api.ResError(c, err)
		return
	}

	result, err := a.menuApp.Query(ctx, params, schema.MenuQueryOptions{
		OrderFields: schema.NewOrderFields(schema.NewOrderField("sequence", schema.OrderByDESC)),
	})
	if err != nil {
		api.ResError(c, err)
		return
	}
	api.ResList(c, result.Data.ToTree())
}

func (a *menu) Get(c *gin.Context) {
	ctx := c.Request.Context()
	item, err := a.menuApp.Get(ctx, c.Param("id"))
	if err != nil {
		api.ResError(c, err)
		return
	}
	api.ResSuccess(c, item)
}

func (a *menu) Create(c *gin.Context) {
	ctx := c.Request.Context()
	var item schema.Menu
	if err := api.ParseJSON(c, &item); err != nil {
		api.ResError(c, err)
		return
	}

	item.Creator = api.GetUserID(c)
	result, err := a.menuApp.Create(ctx, item)
	if err != nil {
		api.ResError(c, err)
		return
	}
	api.ResSuccess(c, result)
}

func (a *menu) Update(c *gin.Context) {
	ctx := c.Request.Context()
	var item schema.Menu
	if err := api.ParseJSON(c, &item); err != nil {
		api.ResError(c, err)
		return
	}

	err := a.menuApp.Update(ctx, c.Param("id"), item)
	if err != nil {
		api.ResError(c, err)
		return
	}
	api.ResOK(c)
}

func (a *menu) Delete(c *gin.Context) {
	ctx := c.Request.Context()
	err := a.menuApp.Delete(ctx, c.Param("id"))
	if err != nil {
		api.ResError(c, err)
		return
	}
	api.ResOK(c)
}

func (a *menu) Enable(c *gin.Context) {
	ctx := c.Request.Context()
	err := a.menuApp.UpdateStatus(ctx, c.Param("id"), 1)
	if err != nil {
		api.ResError(c, err)
		return
	}
	api.ResOK(c)
}

func (a *menu) Disable(c *gin.Context) {
	ctx := c.Request.Context()
	err := a.menuApp.UpdateStatus(ctx, c.Param("id"), 2)
	if err != nil {
		api.ResError(c, err)
		return
	}
	api.ResOK(c)
}
