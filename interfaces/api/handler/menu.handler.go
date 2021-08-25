package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/linzhengen/ddd-gin-admin/application"
	"github.com/linzhengen/ddd-gin-admin/domain/schema"
	"github.com/linzhengen/ddd-gin-admin/infrastructure/ginx"
)

var MenuSet = wire.NewSet(wire.Struct(new(Menu), "*"))

type Menu struct {
	MenuSrv *application.Menu
}

func (a *Menu) Query(c *gin.Context) {
	ctx := c.Request.Context()
	var params schema.MenuQueryParam
	if err := ginx.ParseQuery(c, &params); err != nil {
		ginx.ResError(c, err)
		return
	}

	params.Pagination = true
	result, err := a.MenuSrv.Query(ctx, params, schema.MenuQueryOptions{
		OrderFields: schema.NewOrderFields(schema.NewOrderField("sequence", schema.OrderByDESC)),
	})
	if err != nil {
		ginx.ResError(c, err)
		return
	}
	ginx.ResPage(c, result.Data, result.PageResult)
}

func (a *Menu) QueryTree(c *gin.Context) {
	ctx := c.Request.Context()
	var params schema.MenuQueryParam
	if err := ginx.ParseQuery(c, &params); err != nil {
		ginx.ResError(c, err)
		return
	}

	result, err := a.MenuSrv.Query(ctx, params, schema.MenuQueryOptions{
		OrderFields: schema.NewOrderFields(schema.NewOrderField("sequence", schema.OrderByDESC)),
	})
	if err != nil {
		ginx.ResError(c, err)
		return
	}
	ginx.ResList(c, result.Data.ToTree())
}

func (a *Menu) Get(c *gin.Context) {
	ctx := c.Request.Context()
	item, err := a.MenuSrv.Get(ctx, c.Param("id"))
	if err != nil {
		ginx.ResError(c, err)
		return
	}
	ginx.ResSuccess(c, item)
}

func (a *Menu) Create(c *gin.Context) {
	ctx := c.Request.Context()
	var item schema.Menu
	if err := ginx.ParseJSON(c, &item); err != nil {
		ginx.ResError(c, err)
		return
	}

	item.Creator = ginx.GetUserID(c)
	result, err := a.MenuSrv.Create(ctx, item)
	if err != nil {
		ginx.ResError(c, err)
		return
	}
	ginx.ResSuccess(c, result)
}

func (a *Menu) Update(c *gin.Context) {
	ctx := c.Request.Context()
	var item schema.Menu
	if err := ginx.ParseJSON(c, &item); err != nil {
		ginx.ResError(c, err)
		return
	}

	err := a.MenuSrv.Update(ctx, c.Param("id"), item)
	if err != nil {
		ginx.ResError(c, err)
		return
	}
	ginx.ResOK(c)
}

func (a *Menu) Delete(c *gin.Context) {
	ctx := c.Request.Context()
	err := a.MenuSrv.Delete(ctx, c.Param("id"))
	if err != nil {
		ginx.ResError(c, err)
		return
	}
	ginx.ResOK(c)
}

func (a *Menu) Enable(c *gin.Context) {
	ctx := c.Request.Context()
	err := a.MenuSrv.UpdateStatus(ctx, c.Param("id"), 1)
	if err != nil {
		ginx.ResError(c, err)
		return
	}
	ginx.ResOK(c)
}

func (a *Menu) Disable(c *gin.Context) {
	ctx := c.Request.Context()
	err := a.MenuSrv.UpdateStatus(ctx, c.Param("id"), 2)
	if err != nil {
		ginx.ResError(c, err)
		return
	}
	ginx.ResOK(c)
}
