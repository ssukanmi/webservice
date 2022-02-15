package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ssukanmi/webservice/models"
)

// var err error

func CreateUser(c *gin.Context) {
	user := models.User{}
	c.BindJSON(&user)

	// c.SecureJSON(http.StatusCreated, user)
	c.JSON(http.StatusCreated, user)
}

func GetUser(c *gin.Context) {
}

// func CreateUser(c *gin.Context) {
// 	user := models.User{}
// 	err = c.BindJSON(&user)
// 	if err != nil {
// 		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
// 			"message": "Invalid input json object",
// 		})
// 		return
// 	}

// 	if !services.ValidateEmail(user.Username) {
// 		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
// 			"message": "Invalid email address " + user.Username,
// 		})
// 		return
// 	}

// 	if !services.ValidatePassword(user.Password) {
// 		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
// 			"message": "Empty password field!!!",
// 		})
// 		return
// 	}

// 	user.Password, err = services.HashPassword(user.Password)
// 	if err != nil {
// 		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
// 			"message": "Error hashing password!!!",
// 		})
// 		return
// 	}

// 	// Add DB code here

// 	c.JSON(http.StatusCreated, user)
// }

func UpdateUser(c *gin.Context) {
}
