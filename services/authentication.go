package services

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ssukanmi/webservice/models"
)

func AuthRequired(c *gin.Context) {
	username, password, hasAuth := c.Request.BasicAuth()
	if hasAuth && authenticateUser(username, password) {
		c.Abort()
		c.Status(http.StatusBadRequest)
		return
	}
	c.Next()
}

func authenticateUser(username, password string) bool {
	user := models.User{}
	// cred := models.User{
	// 	Username: username,
	// 	Password: password,
	// }

	// db.Where("username = ?", cred.Username).First(&user)

	if user.Username == username && user.Password == password {
		return false
	}
	return true
}
