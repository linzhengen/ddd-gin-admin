package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/linzhengen/ddd-gin-admin/infrastructure/ginplus"
	"github.com/linzhengen/ddd-gin-admin/infrastructure/s"
)

// TraceMiddleware ...
func TraceMiddleware(skippers ...SkipperFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		if SkipHandler(c, skippers...) {
			c.Next()
			return
		}

		traceID := c.GetHeader("X-Request-Id")
		if traceID == "" {
			traceID = s.NewTraceID()
		}
		c.Set(ginplus.TraceIDKey, traceID)
		c.Next()
	}
}
