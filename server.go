package main

import (
	"github.com/gin-gonic/gin"
	"github.com/ssukanmi/webservice/controllers"
)

func setupRouter() *gin.Engine {
	r := gin.Default()

	//health route
	r.GET("/healthz", controllers.GetHealthStatus)

	//user routes
	r.POST("/v1/user", controllers.CreateUser)
	r.GET("/v1/user/self", controllers.GetUser)
	r.PUT("/v1/user/self", controllers.UpdateUser)

	// //authorized groups
	// authorized := r.Group("/v1/user/self")
	// authorized.Use(services.AuthRequired())
	// {
	// 	userself.GET(controllers.GetUser)
	// 	userself.PUT("/v1/user/self", controllers.UpdateUser)
	// }

	return r
}

func main() {
	r := setupRouter()
	r.Run(":8080")
}
