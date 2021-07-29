package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/linzhengen/ddd-gin-admin/domain/schema"
	"github.com/linzhengen/ddd-gin-admin/infrastructure/ginx"
	"time"
)

var HealthCheckSet = wire.NewSet(wire.Struct(new(HealthCheck), "*"))

type HealthCheck struct {
}

func (a *HealthCheck) Get(c *gin.Context) {
	ginx.ResSuccess(c, &schema.HealthCheck{
		Status: "OK",
		CheckedAt: time.Now(),
	})
}
