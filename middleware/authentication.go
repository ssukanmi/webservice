package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ssukanmi/webservice/repo"
	"github.com/ssukanmi/webservice/service"
	"gorm.io/gorm"
)

func BasicAuth(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		username, password, hasAuth := c.Request.BasicAuth()
		if !(hasAuth && authenticateUser(username, password, db)) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "Unable to authenticate",
			})
			return
		}
		c.Next()
	}
}

func authenticateUser(username string, password string, db *gorm.DB) bool {
	userRepo := repo.NewUserRepository(db)
	user, err := userRepo.FindByUsername(username)
	if err != nil {
		return false
	}
	return service.CheckPasswordHash(password, user.Password)
}
