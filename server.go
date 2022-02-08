package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func getHealthStatus(c *gin.Context) {
	c.Status(http.StatusOK)
}

func setupRouter() *gin.Engine {
	r := gin.Default()
	r.GET("/healthz", getHealthStatus)

	return r
}

func main() {
	r := setupRouter()
	r.Run(":8080")
}
