package mock

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

var MenuSet = wire.NewSet(wire.Struct(new(Menu), "*"))

type Menu struct{}

// Query Search
// @Tags Menu
// @Summary Search
// @Security ApiKeyAuth
// @Param current query int true "Current page" default(1)
// @Param pageSize query int true "Page size" default(10)
// @Param queryValue query string false "Search value"
// @Param status query int false "Status(1:enable 2:disable)"
// @Param showStatus query int false "Show status(1:show 2:hide)"
// @Param parentID query string false "Parent ID"
// @Success 200 {object} schema.ListResult{list=[]schema.Menu} "Search Result"
// @Failure 401 {object} schema.ErrorResult "{error:{code:0,message:Unauthorized}}"
// @Failure 500 {object} schema.ErrorResult "{error:{code:0,message:SystemError}}"
// @Router /api/v1/menus [get]
func (a *Menu) Query(c *gin.Context) {
}

// QueryTree Search menu tree
// @Tags Menu
// @Summary Search menu tree
// @Security ApiKeyAuth
// @Param status query int false "Status(1:enable 2:disable)"
// @Param parentID query string false "Parent ID"
// @Success 200 {object} schema.ListResult{list=[]schema.MenuTree} "Search Result"
// @Failure 401 {object} schema.ErrorResult "{error:{code:0,message:Unauthorized}}"
// @Failure 500 {object} schema.ErrorResult "{error:{code:0,message:SystemError}}"
// @Router /api/v1/menus.tree [get]
func (a *Menu) QueryTree(c *gin.Context) {
}

// Get Search menu
// @Tags Menu
// @Summary Get by ID
// @Security ApiKeyAuth
// @Param id path string true "UUID"
// @Success 200 {object} schema.Menu
// @Failure 401 {object} schema.ErrorResult "{error:{code:0,message:Unauthorized}}"
// @Failure 404 {object} schema.ErrorResult "{error:{code:0,message:NotFound}}"
// @Failure 500 {object} schema.ErrorResult "{error:{code:0,message:SystemError}}"
// @Router /api/v1/menus/{id} [get]
func (a *Menu) Get(c *gin.Context) {
}

// Create Create
// @Tags Menu
// @Summary Create
// @Security ApiKeyAuth
// @Param body body schema.Menu true "Create"
// @Success 200 {object} schema.IDResult
// @Failure 400 {object} schema.ErrorResult "{error:{code:0,message:BadRequest}}"
// @Failure 401 {object} schema.ErrorResult "{error:{code:0,message:Unauthorized}}"
// @Failure 500 {object} schema.ErrorResult "{error:{code:0,message:SystemError}}"
// @Router /api/v1/menus [post]
func (a *Menu) Create(c *gin.Context) {
}

// Update Update
// @Tags Menu
// @Summary Update
// @Security ApiKeyAuth
// @Param id path string true "UUID"
// @Param body body schema.Menu true "Update"
// @Success 200 {object} schema.StatusResult "{status:OK}"
// @Failure 400 {object} schema.ErrorResult "{error:{code:0,message:BadRequest}}"
// @Failure 401 {object} schema.ErrorResult "{error:{code:0,message:Unauthorized}}"
// @Failure 500 {object} schema.ErrorResult "{error:{code:0,message:SystemError}}"
// @Router /api/v1/menus/{id} [put]
func (a *Menu) Update(c *gin.Context) {
}

// Delete Delete
// @Tags Menu
// @Summary Delete
// @Security ApiKeyAuth
// @Param id path string true "UUID"
// @Success 200 {object} schema.StatusResult "{status:OK}"
// @Failure 401 {object} schema.ErrorResult "{error:{code:0,message:Unauthorized}}"
// @Failure 500 {object} schema.ErrorResult "{error:{code:0,message:SystemError}}"
// @Router /api/v1/menus/{id} [delete]
func (a *Menu) Delete(c *gin.Context) {
}

// Enable Enable
// @Tags Menu
// @Summary Enable
// @Security ApiKeyAuth
// @Param id path string true "UUID"
// @Success 200 {object} schema.StatusResult "{status:OK}"
// @Failure 401 {object} schema.ErrorResult "{error:{code:0,message:Unauthorized}}"
// @Failure 500 {object} schema.ErrorResult "{error:{code:0,message:SystemError}}"
// @Router /api/v1/menus/{id}/enable [patch]
func (a *Menu) Enable(c *gin.Context) {
}

// Disable Disable
// @Tags Menu
// @Summary Disable
// @Security ApiKeyAuth
// @Param id path string true "UUID"
// @Success 200 {object} schema.StatusResult "{status:OK}"
// @Failure 401 {object} schema.ErrorResult "{error:{code:0,message:Unauthorized}}"
// @Failure 500 {object} schema.ErrorResult "{error:{code:0,message:SystemError}}"
// @Router /api/v1/menus/{id}/disable [patch]
func (a *Menu) Disable(c *gin.Context) {
}
