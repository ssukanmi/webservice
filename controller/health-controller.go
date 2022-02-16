package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type HealthController interface {
	GetHealthStatus(c *gin.Context)
}

type healthController struct {
}

func NewHealthController() HealthController {
	return &healthController{}
}

func (hc *healthController) GetHealthStatus(c *gin.Context) {
	c.Status(http.StatusOK)
}
