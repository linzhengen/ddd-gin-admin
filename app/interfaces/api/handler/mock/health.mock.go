package mock

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

var HealthCheckSet = wire.NewSet(wire.Struct(new(HealthCheck), "*"))

type HealthCheck struct {
}

// Get HealthCheck
// @Tags HealthCheck
// @Summary HealthCheck
// @Success 200 {object} response.HealthCheck
// @Router /api/health [get]
func (a *HealthCheck) Get(c *gin.Context) {
}
