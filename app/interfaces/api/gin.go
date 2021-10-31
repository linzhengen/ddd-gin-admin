package api

import (
	"fmt"
	"net/http"
	"strings"

	errors2 "github.com/linzhengen/ddd-gin-admin/app/domain/errors"

	"github.com/linzhengen/ddd-gin-admin/app/domain/valueobject/schema"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/linzhengen/ddd-gin-admin/pkg/logger"
	"github.com/linzhengen/ddd-gin-admin/pkg/util/json"
)

const (
	prefix           = "ddd-gin-admin"
	UserIDKey        = prefix + "/user-id"
	ReqBodyKey       = prefix + "/req-body"
	ResBodyKey       = prefix + "/res-body"
	LoggerReqBodyKey = prefix + "/logger-req-body"
)

func GetToken(c *gin.Context) string {
	var token string
	auth := c.GetHeader("Authorization")
	prefix := "Bearer "
	if auth != "" && strings.HasPrefix(auth, prefix) {
		token = auth[len(prefix):]
	}
	return token
}

func GetUserID(c *gin.Context) string {
	return c.GetString(UserIDKey)
}

func SetUserID(c *gin.Context, userID string) {
	c.Set(UserIDKey, userID)
}

func GetBody(c *gin.Context) []byte {
	if v, ok := c.Get(ReqBodyKey); ok {
		if b, ok := v.([]byte); ok {
			return b
		}
	}
	return nil
}

func ParseJSON(c *gin.Context, obj interface{}) error {
	if err := c.ShouldBindJSON(obj); err != nil {
		return errors2.Wrap400Response(err, fmt.Sprintf("400 Bad Request - %s", err.Error()))
	}
	return nil
}

func ParseQuery(c *gin.Context, obj interface{}) error {
	if err := c.ShouldBindQuery(obj); err != nil {
		return errors2.Wrap400Response(err, fmt.Sprintf("400 Bad Request - %s", err.Error()))
	}
	return nil
}

func ParseForm(c *gin.Context, obj interface{}) error {
	if err := c.ShouldBindWith(obj, binding.Form); err != nil {
		return errors2.Wrap400Response(err, fmt.Sprintf("解析请求参数发生错误 - %s", err.Error()))
	}
	return nil
}

func ResOK(c *gin.Context) {
	ResSuccess(c, schema.StatusResult{Status: schema.OKStatus})
}

func ResList(c *gin.Context, v interface{}) {
	ResSuccess(c, schema.ListResult{List: v})
}

func ResPage(c *gin.Context, v interface{}, pr *schema.PaginationResult) {
	list := schema.ListResult{
		List:       v,
		Pagination: pr,
	}
	ResSuccess(c, list)
}

func ResSuccess(c *gin.Context, v interface{}) {
	ResJSON(c, http.StatusOK, v)
}

func ResJSON(c *gin.Context, status int, v interface{}) {
	buf, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	c.Set(ResBodyKey, buf)
	c.Data(status, "application/json; charset=utf-8", buf)
	c.Abort()
}

func ResError(c *gin.Context, err error, status ...int) {
	ctx := c.Request.Context()
	var res *errors2.ResponseError

	if err != nil {
		if e, ok := err.(*errors2.ResponseError); ok {
			res = e
		} else {
			res = errors2.UnWrapResponse(errors2.ErrInternalServer)
			res.ERR = err
		}
	} else {
		res = errors2.UnWrapResponse(errors2.ErrInternalServer)
	}

	if len(status) > 0 {
		res.StatusCode = status[0]
	}

	if err := res.ERR; err != nil {
		if res.Message == "" {
			res.Message = err.Error()
		}

		if status := res.StatusCode; status >= 400 && status < 500 {
			logger.WithContext(ctx).Warnf(err.Error())
		} else if status >= 500 {
			logger.WithContext(logger.NewStackContext(ctx, err)).Errorf(err.Error())
		}
	}

	eitem := schema.ErrorItem{
		Code:    res.Code,
		Message: res.Message,
	}
	ResJSON(c, res.StatusCode, schema.ErrorResult{Error: eitem})
}
