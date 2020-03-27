package handler

import (
	"github.com/linzhengen/ddd-gin-admin/infrastructure/s"

	"github.com/gin-gonic/gin"
	"github.com/linzhengen/ddd-gin-admin/application"
	"github.com/linzhengen/ddd-gin-admin/domain/schema"
	"github.com/linzhengen/ddd-gin-admin/infrastructure/ginplus"
)

// NewUser constructor.
func NewUser(us application.UserRepository) *User {
	return &User{
		us: us,
	}
}

// User struct defines the dependencies that will be used.
type User struct {
	us application.UserRepository
}

// Query find user use params.
func (u *User) Query(c *gin.Context) {
	var params schema.UserQueryParam
	params.LikeUserName = c.Query("userName")
	params.LikeRealName = c.Query("realName")
	if v := s.S(c.Query("status")).DefaultInt(0); v > 0 {
		params.Status = v
	}

	result, err := u.us.QueryShow(ginplus.NewContext(c), params, schema.UserQueryOptions{
		PageParam: ginplus.GetPaginationParam(c),
	})
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResPage(c, result.Data, result.PageResult)
}

// Get gets user.
func (u *User) Get(c *gin.Context) {
	item, err := u.us.Get(ginplus.NewContext(c), c.Param("id"))
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, item.CleanSecure())
}

// Create create user.
func (u *User) Create(c *gin.Context) {
	var item schema.User
	if err := ginplus.ParseJSON(c, &item); err != nil {
		ginplus.ResError(c, err)
		return
	}
	item.Creator = ginplus.GetUserID(c)
	nitem, err := u.us.Create(ginplus.NewContext(c), item)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, nitem.CleanSecure())
}

// Update update user.
func (u *User) Update(c *gin.Context) {
	var item schema.User
	if err := ginplus.ParseJSON(c, &item); err != nil {
		ginplus.ResError(c, err)
		return
	}

	nitem, err := u.us.Update(ginplus.NewContext(c), c.Param("id"), item)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, nitem.CleanSecure())
}

// Delete delete user.
func (u *User) Delete(c *gin.Context) {
	err := u.us.Delete(ginplus.NewContext(c), c.Param("id"))
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResOK(c)
}

// Enable enable user.
func (u *User) Enable(c *gin.Context) {
	err := u.us.UpdateStatus(ginplus.NewContext(c), c.Param("id"), 1)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResOK(c)
}

// Disable user.
func (u *User) Disable(c *gin.Context) {
	err := u.us.UpdateStatus(ginplus.NewContext(c), c.Param("id"), 2)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResOK(c)
}
