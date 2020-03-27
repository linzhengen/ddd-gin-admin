package ginplus

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/linzhengen/ddd-gin-admin/infrastructure/logger"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/linzhengen/ddd-gin-admin/domain/schema"
	icontext "github.com/linzhengen/ddd-gin-admin/infrastructure/context"
	"github.com/linzhengen/ddd-gin-admin/infrastructure/errors"
	"github.com/linzhengen/ddd-gin-admin/infrastructure/json"
	"github.com/linzhengen/ddd-gin-admin/infrastructure/s"
)

// context keys
const (
	prefix = "ddd-gin-admin"
	// UserIDKey for user-id
	UserIDKey = prefix + "/user-id"
	// TraceIDKey for trace-id
	TraceIDKey = prefix + "/trace-id"
	// ResBodyKey for res-body
	ResBodyKey = prefix + "/res-body"
)

// NewContext ...
func NewContext(c *gin.Context) context.Context {
	parent := context.Background()

	if v := GetTraceID(c); v != "" {
		parent = icontext.NewTraceID(parent, v)
		parent = logger.NewTraceIDContext(parent, v)
	}

	if v := GetUserID(c); v != "" {
		parent = icontext.NewUserID(parent, v)
		parent = logger.NewUserIDContext(parent, v)
	}

	return parent
}

// GetToken gets token from context.
func GetToken(c *gin.Context) string {
	var token string
	auth := c.GetHeader("Authorization")
	prefix := "Bearer "
	if auth != "" && strings.HasPrefix(auth, prefix) {
		token = auth[len(prefix):]
	}
	return token
}

// GetPageIndex gets page index.
func GetPageIndex(c *gin.Context) int {
	defaultVal := 1
	if v := c.Query("current"); v != "" {
		if iv := s.S(v).DefaultInt(defaultVal); iv > 0 {
			return iv
		}
	}
	return defaultVal
}

// GetPageSize gets page(max: 50).
func GetPageSize(c *gin.Context) int {
	defaultVal := 10
	if v := c.Query("pageSize"); v != "" {
		if iv := s.S(v).DefaultInt(defaultVal); iv > 0 {
			if iv > 50 {
				iv = 50
			}
			return iv
		}
	}
	return defaultVal
}

// GetPaginationParam gets pagination parameters.
func GetPaginationParam(c *gin.Context) *schema.PaginationParam {
	return &schema.PaginationParam{
		PageIndex: GetPageIndex(c),
		PageSize:  GetPageSize(c),
	}
}

// GetTraceID gets trace id.
func GetTraceID(c *gin.Context) string {
	return c.GetString(TraceIDKey)
}

// GetUserID gets user id.
func GetUserID(c *gin.Context) string {
	return c.GetString(UserIDKey)
}

// SetUserID sets user id.
func SetUserID(c *gin.Context, userID string) {
	c.Set(UserIDKey, userID)
}

// ParseJSON parse json.
func ParseJSON(c *gin.Context, obj interface{}) error {
	if err := c.ShouldBindJSON(obj); err != nil {
		return errors.Wrap400Response(err, "parse json error")
	}
	return nil
}

// ParseQuery parse query parameters.
func ParseQuery(c *gin.Context, obj interface{}) error {
	if err := c.ShouldBindQuery(obj); err != nil {
		return errors.Wrap400Response(err, "parse query error")
	}
	return nil
}

// ParseForm parse form.
func ParseForm(c *gin.Context, obj interface{}) error {
	if err := c.ShouldBindWith(obj, binding.Form); err != nil {
		return errors.Wrap400Response(err, "parse form error")
	}
	return nil
}

// ResPage responses pages.
func ResPage(c *gin.Context, v interface{}, pr *schema.PaginationResult) {
	list := schema.HTTPList{
		List: v,
		Pagination: &schema.HTTPPagination{
			Current:  GetPageIndex(c),
			PageSize: GetPageSize(c),
		},
	}
	if pr != nil {
		list.Pagination.Total = pr.Total
	}

	ResSuccess(c, list)
}

// ResList list for responses.
func ResList(c *gin.Context, v interface{}) {
	ResSuccess(c, schema.HTTPList{List: v})
}

// ResOK ok response.
func ResOK(c *gin.Context) {
	ResSuccess(c, schema.HTTPStatus{Status: schema.OKStatusText.String()})
}

// ResSuccess success json response.
func ResSuccess(c *gin.Context, v interface{}) {
	ResJSON(c, http.StatusOK, v)
}

// ResJSON json response.
func ResJSON(c *gin.Context, status int, v interface{}) {
	buf, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	c.Set(ResBodyKey, buf)
	c.Data(status, "application/json; charset=utf-8", buf)
	c.Abort()
}

// ResError error response.
func ResError(c *gin.Context, err error, status ...int) {
	var res *errors.ResponseError
	if err != nil {
		if e, ok := err.(*errors.ResponseError); ok {
			res = e
		} else {
			res = errors.UnWrapResponse(errors.Wrap500Response(err))
		}
	} else {
		res = errors.UnWrapResponse(errors.ErrInternalServer)
	}

	if len(status) > 0 {
		res.StatusCode = status[0]
	}

	if err := res.ERR; err != nil {
		if status := res.StatusCode; status >= 400 && status < 500 {
			logger.StartSpan(NewContext(c)).Warnf(err.Error())
		} else if status >= 500 {
			span := logger.StartSpan(NewContext(c))
			span = span.WithField("stack", fmt.Sprintf("%+v", err))
			span.Errorf(err.Error())
		}
	}

	eitem := schema.HTTPErrorItem{
		Code:    res.Code,
		Message: res.Message,
	}
	ResJSON(c, res.StatusCode, schema.HTTPError{Error: eitem})
}
