package handler

import (
	"time"

	"github.com/linzhengen/ddd-gin-admin/app/interfaces/api"

	"github.com/gin-gonic/gin"
	"github.com/linzhengen/ddd-gin-admin/app/domain/schema"
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
	api.ResSuccess(c, &schema.HealthCheck{
		Status:    "OK",
		CheckedAt: time.Now(),
	})
}
