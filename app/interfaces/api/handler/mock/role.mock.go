package mock

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

var RoleSet = wire.NewSet(wire.Struct(new(Role), "*"))

type Role struct {
}

// Query Search role
// @Tags Role
// @Summary Search role
// @Security ApiKeyAuth
// @Param current query int true "Current page" default(1)
// @Param pageSize query int true "Page size" default(10)
// @Param queryValue query string false "Search value"
// @Param status query int false "Status(1:enable 2:disable)"
// @Success 200 {object} response.ListResult{list=[]schema.Role} "Search Result"
// @Failure 401 {object} response.ErrorResult "{error:{code:0,message:Unauthorized}}"
// @Failure 500 {object} response.ErrorResult "{error:{code:0,message:SystemError}}"
// @Router /api/v1/roles [get]
func (a *Role) Query(c *gin.Context) {
}

// QuerySelect Search selected data
// @Tags Role
// @Summary Search selected data
// @Security ApiKeyAuth
// @Param queryValue query string false "Search value"
// @Param status query int false "Status(1:enable 2:disable)"
// @Success 200 {object} response.ListResult{list=[]schema.Role} "Search Result"
// @Failure 400 {object} response.ErrorResult "{error:{code:0,message:未知的查询类型}}"
// @Failure 401 {object} response.ErrorResult "{error:{code:0,message:Unauthorized}}"
// @Failure 500 {object} response.ErrorResult "{error:{code:0,message:SystemError}}"
// @Router /api/v1/roles.select [get]
func (a *Role) QuerySelect(c *gin.Context) {
}

// Get Get by ID
// @Tags Role
// @Summary Get by ID
// @Security ApiKeyAuth
// @Param id path string true "UUID"
// @Success 200 {object} response.Role
// @Failure 401 {object} response.ErrorResult "{error:{code:0,message:Unauthorized}}"
// @Failure 404 {object} response.ErrorResult "{error:{code:0,message:NotFound}}"
// @Failure 500 {object} response.ErrorResult "{error:{code:0,message:SystemError}}"
// @Router /api/v1/roles/{id} [get]
func (a *Role) Get(c *gin.Context) {
}

// Create Create
// @Tags Role
// @Summary Create
// @Security ApiKeyAuth
// @Param body body request.Role true "Create"
// @Success 200 {object} response.IDResult
// @Failure 400 {object} response.ErrorResult "{error:{code:0,message:BadRequest}}"
// @Failure 401 {object} response.ErrorResult "{error:{code:0,message:Unauthorized}}"
// @Failure 500 {object} response.ErrorResult "{error:{code:0,message:SystemError}}"
// @Router /api/v1/roles [post]
func (a *Role) Create(c *gin.Context) {
}

// Update Update
// @Tags Role
// @Summary Update
// @Security ApiKeyAuth
// @Param id path string true "UUID"
// @Param body body request.Role true "Update"
// @Success 200 {object} response.Role
// @Failure 400 {object} response.ErrorResult "{error:{code:0,message:BadRequest}}"
// @Failure 401 {object} response.ErrorResult "{error:{code:0,message:Unauthorized}}"
// @Failure 500 {object} response.ErrorResult "{error:{code:0,message:SystemError}}"
// @Router /api/v1/roles/{id} [put]
func (a *Role) Update(c *gin.Context) {
}

// Delete Delete
// @Tags Role
// @Summary Delete
// @Security ApiKeyAuth
// @Param id path string true "UUID"
// @Success 200 {object} response.StatusResult "{status:OK}"
// @Failure 401 {object} response.ErrorResult "{error:{code:0,message:Unauthorized}}"
// @Failure 500 {object} response.ErrorResult "{error:{code:0,message:SystemError}}"
// @Router /api/v1/roles/{id} [delete]
func (a *Role) Delete(c *gin.Context) {
}

// Enable Enable
// @Tags Role
// @Summary Enable
// @Security ApiKeyAuth
// @Param id path string true "UUID"
// @Success 200 {object} response.StatusResult "{status:OK}"
// @Failure 401 {object} response.ErrorResult "{error:{code:0,message:Unauthorized}}"
// @Failure 500 {object} response.ErrorResult "{error:{code:0,message:SystemError}}"
// @Router /api/v1/roles/{id}/enable [patch]
func (a *Role) Enable(c *gin.Context) {
}

// Disable Disable
// @Tags Role
// @Summary Disable
// @Security ApiKeyAuth
// @Param id path string true "UUID"
// @Success 200 {object} response.StatusResult "{status:OK}"
// @Failure 401 {object} response.ErrorResult "{error:{code:0,message:Unauthorized}}"
// @Failure 500 {object} response.ErrorResult "{error:{code:0,message:SystemError}}"
// @Router /api/v1/roles/{id}/disable [patch]
func (a *Role) Disable(c *gin.Context) {
}
