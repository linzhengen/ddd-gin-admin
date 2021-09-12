package handler

import (
	"time"

	"github.com/linzhengen/ddd-gin-admin/app/domain/schema"
	"github.com/linzhengen/ddd-gin-admin/app/infrastructure/ginx"

	"github.com/gin-gonic/gin"
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
