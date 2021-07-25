package mock

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

var DemoSet = wire.NewSet(wire.Struct(new(Demo), "*"))

type Demo struct {
}

// Query Search
// @Tags Demo
// @Security ApiKeyAuth
// @Summary Search Summary
// @Param current query int true "Current page" default(1)
// @Param pageSize query int true "Page size" default(10)
// @Param queryValue query string false "Search value"
// @Success 200 {object} schema.ListResult{list=[]schema.Demo} "Search result"
// @Failure 401 {object} schema.ErrorResult "{error:{code:0,message:Unauthorized}}"
// @Failure 500 {object} schema.ErrorResult "{error:{code:0,message:SystemError}}"
// @Router /api/v1/demos [get]
func (a *Demo) Query(c *gin.Context) {
}

// Get Search
// @Tags Demo
// @Security ApiKeyAuth
// @Summary Search
// @Param id path string true "UUID"
// @Success 200 {object} schema.Demo
// @Failure 401 {object} schema.ErrorResult "{error:{code:0,message:Unauthorized}}"
// @Failure 404 {object} schema.ErrorResult "{error:{code:0,message:NotFound}}"
// @Failure 500 {object} schema.ErrorResult "{error:{code:0,message:SystemError}}"
// @Router /api/v1/demos/{id} [get]
func (a *Demo) Get(c *gin.Context) {
}

// Create Create
// @Tags Demo
// @Security ApiKeyAuth
// @Summary Create
// @Param body body schema.Demo true "Create"
// @Success 200 {object} schema.IDResult
// @Failure 400 {object} schema.ErrorResult "{error:{code:0,message:BadRequest}}"
// @Failure 401 {object} schema.ErrorResult "{error:{code:0,message:Unauthorized}}"
// @Failure 500 {object} schema.ErrorResult "{error:{code:0,message:SystemError}}"
// @Router /api/v1/demos [post]
func (a *Demo) Create(c *gin.Context) {
}

// Update Update
// @Tags Demo
// @Security ApiKeyAuth
// @Summary Update
// @Param id path string true "UUID"
// @Param body body schema.Demo true "Update"
// @Success 200 {object} schema.StatusResult "{status:OK}"
// @Failure 400 {object} schema.ErrorResult "{error:{code:0,message:BadRequest}}"
// @Failure 401 {object} schema.ErrorResult "{error:{code:0,message:Unauthorized}}"
// @Failure 500 {object} schema.ErrorResult "{error:{code:0,message:SystemError}}"
// @Router /api/v1/demos/{id} [put]
func (a *Demo) Update(c *gin.Context) {
}

// Delete Delete
// @Tags Demo
// @Security ApiKeyAuth
// @Summary Delete
// @Param id path string true "UUID"
// @Success 200 {object} schema.StatusResult "{status:OK}"
// @Failure 401 {object} schema.ErrorResult "{error:{code:0,message:Unauthorized}}"
// @Failure 500 {object} schema.ErrorResult "{error:{code:0,message:SystemError}}"
// @Router /api/v1/demos/{id} [delete]
func (a *Demo) Delete(c *gin.Context) {
}

// Enable Enable
// @Tags Demo
// @Security ApiKeyAuth
// @Summary Enable
// @Param id path string true "UUID"
// @Success 200 {object} schema.StatusResult "{status:OK}"
// @Failure 401 {object} schema.ErrorResult "{error:{code:0,message:Unauthorized}}"
// @Failure 500 {object} schema.ErrorResult "{error:{code:0,message:SystemError}}"
// @Router /api/v1/demos/{id}/enable [patch]
func (a *Demo) Enable(c *gin.Context) {
}

// Disable Disable
// @Tags Demo
// @Security ApiKeyAuth
// @Summary Disable
// @Param id path string true "UUID"
// @Success 200 {object} schema.StatusResult "{status:OK}"
// @Failure 401 {object} schema.ErrorResult "{error:{code:0,message:Unauthorized}}"
// @Failure 500 {object} schema.ErrorResult "{error:{code:0,message:SystemError}}"
// @Router /api/v1/demos/{id}/disable [patch]
func (a *Demo) Disable(c *gin.Context) {
}
