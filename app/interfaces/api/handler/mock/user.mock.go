package mock

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

var UserSet = wire.NewSet(wire.Struct(new(User), "*"))

type User struct {
}

// Query Search user
// @Tags User
// @Summary Search user
// @Security ApiKeyAuth
// @Param current query int true "Current page" default(1)
// @Param pageSize query int true "Page size" default(10)
// @Param queryValue query string false "Search value"
// @Param roleIDs query string false "Role IDs(Comma division)"
// @Param status query int false "Status (1: Enabled 2: Disabled)"
// @Success 200 {object} response.ListResult{list=[]response.UserShow} "Search Result"
// @Failure 401 {object} response.ErrorResult "{error:{code:0,message:Unauthorized}}"
// @Failure 500 {object} response.ErrorResult "{error:{code:0,message:SystemError}}"
// @Router /api/v1/users [get]
func (a *User) Query(c *gin.Context) {
}

// Get Get by ID
// @Tags User
// @Summary Get by ID
// @Security ApiKeyAuth
// @Param id path string true "UUID"
// @Success 200 {object} response.User
// @Failure 401 {object} response.ErrorResult "{error:{code:0,message:Unauthorized}}"
// @Failure 404 {object} response.ErrorResult "{error:{code:0,message:NotFound}}"
// @Failure 500 {object} response.ErrorResult "{error:{code:0,message:SystemError}}"
// @Router /api/v1/users/{id} [get]
func (a *User) Get(c *gin.Context) {
}

// Create Create
// @Tags User
// @Summary Create
// @Security ApiKeyAuth
// @Param body body request.User true "Create"
// @Success 200 {object} response.IDResult
// @Failure 400 {object} response.ErrorResult "{error:{code:0,message:BadRequest}}"
// @Failure 401 {object} response.ErrorResult "{error:{code:0,message:Unauthorized}}"
// @Failure 500 {object} response.ErrorResult "{error:{code:0,message:SystemError}}"
// @Router /api/v1/users [post]
func (a *User) Create(c *gin.Context) {
}

// Update Update
// @Tags User
// @Summary Update
// @Security ApiKeyAuth
// @Param id path string true "UUID"
// @Param body body request.User true "Update"
// @Success 200 {object} response.User
// @Failure 400 {object} response.ErrorResult "{error:{code:0,message:BadRequest}}"
// @Failure 401 {object} response.ErrorResult "{error:{code:0,message:Unauthorized}}"
// @Failure 500 {object} response.ErrorResult "{error:{code:0,message:SystemError}}"
// @Router /api/v1/users/{id} [put]
func (a *User) Update(c *gin.Context) {
}

// Delete Delete
// @Tags User
// @Summary Delete
// @Security ApiKeyAuth
// @Param id path string true "UUID"
// @Success 200 {object} response.StatusResult "{status:OK}"
// @Failure 401 {object} response.ErrorResult "{error:{code:0,message:Unauthorized}}"
// @Failure 500 {object} response.ErrorResult "{error:{code:0,message:SystemError}}"
// @Router /api/v1/users/{id} [delete]
func (a *User) Delete(c *gin.Context) {
}

// Enable Enable
// @Tags User
// @Summary Enable
// @Security ApiKeyAuth
// @Param id path string true "UUID"
// @Success 200 {object} response.StatusResult "{status:OK}"
// @Failure 401 {object} response.ErrorResult "{error:{code:0,message:Unauthorized}}"
// @Failure 500 {object} response.ErrorResult "{error:{code:0,message:SystemError}}"
// @Router /api/v1/users/{id}/enable [patch]
func (a *User) Enable(c *gin.Context) {
}

// Disable Disable
// @Tags User
// @Summary Disable
// @Security ApiKeyAuth
// @Param id path string true "UUID"
// @Success 200 {object} response.StatusResult "{status:OK}"
// @Failure 401 {object} response.ErrorResult "{error:{code:0,message:Unauthorized}}"
// @Failure 500 {object} response.ErrorResult "{error:{code:0,message:SystemError}}"
// @Router /api/v1/users/{id}/disable [patch]
func (a *User) Disable(c *gin.Context) {
}
