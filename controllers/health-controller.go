package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetHealthStatus(c *gin.Context) {
	c.Status(http.StatusOK)
}
