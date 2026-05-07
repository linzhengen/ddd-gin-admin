package middleware

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/linzhengen/ddd-gin-admin/app/interfaces/api"
	"github.com/linzhengen/ddd-gin-admin/configs"
)

func CopyBodyMiddleware(skippers ...SkipperFunc) gin.HandlerFunc {
	var maxMemory int64 = 64 << 20 // 64 MB
	if v := configs.C.HTTP.MaxContentLength; v > 0 {
		maxMemory = v
	}

	return func(c *gin.Context) {
		if SkipHandler(c, skippers...) || c.Request.Body == nil {
			c.Next()
			return
		}

		safe := &io.LimitedReader{R: c.Request.Body, N: maxMemory}
		rawBody, _ := io.ReadAll(safe)

		var requestBody []byte
		if c.GetHeader("Content-Encoding") == "gzip" {
			reader, err := gzip.NewReader(bytes.NewReader(rawBody))
			if err == nil {
				requestBody, _ = io.ReadAll(reader)
				reader.Close()
			} else {
				requestBody = rawBody
			}
		} else {
			requestBody = rawBody
		}

		c.Request.Body.Close()
		bf := bytes.NewBuffer(requestBody)
		c.Request.Body = http.MaxBytesReader(c.Writer, io.NopCloser(bf), maxMemory)
		c.Set(api.ReqBodyKey, requestBody)

		c.Next()
	}
}
