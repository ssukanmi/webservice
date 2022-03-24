package controller

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/joho/godotenv"
	"github.com/mashingan/smapping"
	"github.com/ssukanmi/webservice/dto"
	"github.com/ssukanmi/webservice/entity"
	"github.com/ssukanmi/webservice/repo"
	"gorm.io/gorm"
)

var (
	s3BucketName = os.Getenv("S3_BUCKETNAME")
)

type UserController interface {
	CreateUser(c *gin.Context)
	GetUser(c *gin.Context)
	UpdateUser(c *gin.Context)
	AddOrUpdateProfilePic(c *gin.Context)
	GetProfilePic(c *gin.Context)
	DeleteProfilePic(c *gin.Context)
}

type userController struct {
	userRepo repo.UserRepository
}

func NewUserController(userRepo repo.UserRepository) UserController {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Failed to load env file")
	}
	s3BucketName = os.Getenv("S3_BUCKETNAME")
	return &userController{
		userRepo: userRepo,
	}
}

func (uc *userController) CreateUser(c *gin.Context) {
	userCreateDTO := dto.UserCreateDTO{}
	err := c.ShouldBindJSON(&userCreateDTO)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Unable to bind json body -- " + err.Error(),
		})
		return
	}

	err = userCreateDTO.Validate()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Unable to validate json value -- " + err.Error(),
		})
		return
	}

	user := entity.User{}
	err = smapping.FillStruct(&user, smapping.MapFields(&userCreateDTO))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Failed to swap maps -- " + err.Error(),
		})
		return
	}

	user, err = uc.userRepo.InsertUser(user)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Unable to insert user to database -- " + err.Error(),
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
			"message": "Unable to get user from the database -- " + err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, user)
}

func (uc *userController) UpdateUser(c *gin.Context) {
	userUpdateDTO := dto.UserUpdateDTO{}
	binding.EnableDecoderDisallowUnknownFields = true
	err := c.ShouldBindJSON(&userUpdateDTO)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Unable to bind json body -- " + err.Error(),
		})
		return
	}

	user := entity.User{}
	err = smapping.FillStruct(&user, smapping.MapFields(&userUpdateDTO))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Failed to swap maps -- " + err.Error(),
		})
		return
	}

	username, _, _ := c.Request.BasicAuth()
	user, err = uc.userRepo.UpdateUser(username, user)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Unable to update user in the database -- " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusNoContent, user)
}

func (uc *userController) AddOrUpdateProfilePic(c *gin.Context) {
	userImage := entity.UserImage{}
	username, _, _ := c.Request.BasicAuth()
	file, err := c.FormFile("file")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Unable to upload profile pic -- " + err.Error(),
		})
		return
	}
	fileType := file.Header.Get("Content-Type")
	if !(strings.HasPrefix(fileType, "image/")) {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Wrong file type",
		})
		return
	}
	os.MkdirAll(s3BucketName+"/"+username, os.ModePerm)
	err = c.SaveUploadedFile(file, s3BucketName+"/"+username+"/"+file.Filename)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Unable to upload profile pic -- " + err.Error(),
		})
		return
	}
	userImage, err = uc.userRepo.UpdateUserProfilePic(username, file)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Unable to add/update user profile picture -- " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, userImage)
}

func (uc *userController) GetProfilePic(c *gin.Context) {
	username, _, _ := c.Request.BasicAuth()
	userImage, err := uc.userRepo.GetUserProfilePic(username)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"message": "Unable to get user profile picture -- " + err.Error(),
			})
			return
		}
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Unable to get user profile picture -- " + err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, userImage)
}

func (uc *userController) DeleteProfilePic(c *gin.Context) {
	username, _, _ := c.Request.BasicAuth()
	os.RemoveAll(s3BucketName + "/" + username)
	err := uc.userRepo.DeleteUserProfilePic(username)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Unable to delete user profile picture -- " + err.Error(),
		})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
