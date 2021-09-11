package handler

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/linzhengen/ddd-gin-admin/domain/schema"
	"github.com/linzhengen/ddd-gin-admin/infrastructure/ginx"
)

type HealthCheck interface {
	Get(c *gin.Context)
}

func NewHealthCheck() HealthCheck {
	return &healthCheck{}
}

type healthCheck struct {
}

func (a *healthCheck) Get(c *gin.Context) {
	ginx.ResSuccess(c, &schema.HealthCheck{
		Status:    "OK",
		CheckedAt: time.Now(),
	})
}
