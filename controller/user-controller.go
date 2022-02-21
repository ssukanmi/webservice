package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/mashingan/smapping"
	"github.com/ssukanmi/webservice/dto"
	"github.com/ssukanmi/webservice/entity"
	"github.com/ssukanmi/webservice/repo"
)

type UserController interface {
	CreateUser(c *gin.Context)
	GetUser(c *gin.Context)
	UpdateUser(c *gin.Context)
}

type userController struct {
	userRepo repo.UserRepository
}

func NewUserController(userRepo repo.UserRepository) UserController {
	return &userController{
		userRepo: userRepo,
	}
}

func (uc *userController) CreateUser(c *gin.Context) {
	userCreateDTO := dto.UserCreateDTO{}
	err := c.ShouldBindJSON(&userCreateDTO)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Unable to bind json body" + err.Error(),
		})
		return
	}

	err = userCreateDTO.Validate()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Unable to validate json value" + err.Error(),
		})
		return
	}

	user := entity.User{}
	err = smapping.FillStruct(&user, smapping.MapFields(&userCreateDTO))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Failed to swap maps" + err.Error(),
		})
		return
	}

	user, err = uc.userRepo.InsertUser(user)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Unable to insert user to database" + err.Error(),
		})
		return
	}
	c.JSON(http.StatusCreated, user)
}

func (uc *userController) GetUser(c *gin.Context) {
	username, _, _ := c.Request.BasicAuth()
	user, err := uc.userRepo.FindByUsername(username)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Unable to user get from the database" + err.Error(),
		})
		return
	}
	c.JSON(http.StatusCreated, user)
}

func (uc *userController) UpdateUser(c *gin.Context) {
	userUpdateDTO := dto.UserUpdateDTO{}
	binding.EnableDecoderDisallowUnknownFields = true
	err := c.ShouldBindJSON(&userUpdateDTO)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Unable to bind json body" + err.Error(),
		})
		return
	}

	user := entity.User{}
	err = smapping.FillStruct(&user, smapping.MapFields(&userUpdateDTO))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Failed to swap maps" + err.Error(),
		})
		return
	}

	username, _, _ := c.Request.BasicAuth()
	user, err = uc.userRepo.UpdateUser(username, user)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Unable to update user in the database" + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, user)
}
