package main

import (
	"io"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/ssukanmi/webservice/config"
	"github.com/ssukanmi/webservice/controller"
	"github.com/ssukanmi/webservice/middleware"
	"github.com/ssukanmi/webservice/repo"
	"gorm.io/gorm"
)

var (
	db               *gorm.DB                    = config.SetupDatabaseConnection()
	healthController controller.HealthController = controller.NewHealthController()
	userRepo         repo.UserRepository         = repo.NewUserRepository(db)
	userController   controller.UserController   = controller.NewUserController(userRepo)
)

func setupRouter() *gin.Engine {
	r := gin.Default()

	//health route
	r.GET("/healthz", healthController.GetHealthStatus)

	//user routes
	r.POST("/v1/user", userController.CreateUser)
	//authenticated user routes
	authRoutes := r.Group("/v1/user/self", middleware.BasicAuth(db))
	{
		authRoutes.GET("", userController.GetUser)
		authRoutes.PUT("", userController.UpdateUser)
	}

	return r
}

func main() {
	defer config.CloseDatabaseConnection(db)
	f, _ := os.Create("gin.log")
	gin.DefaultWriter = io.MultiWriter(f, os.Stdout)
	r := setupRouter()
	r.Run(":8080")
}
