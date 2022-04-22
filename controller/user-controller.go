package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/joho/godotenv"
	"github.com/mashingan/smapping"
	"github.com/ssukanmi/webservice/dto"
	"github.com/ssukanmi/webservice/entity"
	"github.com/ssukanmi/webservice/repo"
	"github.com/ssukanmi/webservice/service"
	"gorm.io/gorm"
)

var (
	s3BucketName  = os.Getenv("S3_BUCKETNAME")
	dynamobdTable = os.Getenv("DYNAMODB_TABLE")
	accountID     = os.Getenv("ACCOUNT_ID")
)

type UserController interface {
	CreateUser(c *gin.Context)
	VerifyUserEmail(c *gin.Context)
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
	dynamobdTable = os.Getenv("DYNAMODB_TABLE")
	accountID = os.Getenv("ACCOUNT_ID")
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
	sess, err := session.NewSession()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Error creating token (aws Seeion connect) -- " + err.Error(),
		})
		return
	}
	token, err := service.GenerateToken(user.Username)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Error creating token (Token generation) -- " + err.Error(),
		})
		return
	}
	item := entity.EmailToken{
		Email: user.Username,
		Token: token,
		TTL:   time.Now().Add(time.Second * 300).Unix(),
	}
	dydbSVC := dynamodb.New(sess)
	av, err := dynamodbattribute.MarshalMap(item)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Error creating token (dynamodbattribute marshal mapping) -- " + err.Error(),
		})
		return
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(dynamobdTable),
	}
	_, err = dydbSVC.PutItem(input)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Error creating token (putting token Item in dynamo) -- " + err.Error(),
		})
		return
	}

	etMessage := entity.EmailTokenMessage{
		Email:       user.Username,
		Token:       token,
		MessageType: "publish message",
	}
	etMessageStr, _ := json.Marshal(etMessage)
	message := entity.Message{
		Default: string(etMessageStr),
	}
	messageStr, _ := json.Marshal(message)

	snsSVC := sns.New(sess)
	_, err = snsSVC.Publish(&sns.PublishInput{
		Message:          aws.String(string(messageStr)),
		TopicArn:         aws.String("arn:aws:sns:us-east-1:" + accountID + ":EmailVerificationTopic"),
		MessageStructure: aws.String("json"),
	})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Error creating token (Posting message to topic) -- " + err.Error(),
		})
		return
	}
	c.JSON(http.StatusCreated, user)
}

func (uc *userController) VerifyUserEmail(c *gin.Context) {
	userEmail := ""
	userToken := ""
	if email, ok := c.GetQuery("email"); ok {
		userEmail = email
	} else {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Error verifying email (Invalid query string no email)",
		})
		return
	}
	if token, ok := c.GetQuery("token"); ok {
		userToken = token
	} else {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Error verifying email (Invalid query string no token)",
		})
		return
	}

	user, err := uc.userRepo.FindByUsername(userEmail)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Error verifying email (user doesn't exist) -- " + err.Error(),
		})
		return
	}
	if user.Verified {
		c.AbortWithStatusJSON(http.StatusAccepted, gin.H{
			"message": "User has already been verified",
		})
		return
	}

	sess, err := session.NewSession()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Error verifying email (aws Seeion connect) -- " + err.Error(),
		})
		return
	}
	dydbSVC := dynamodb.New(sess)
	result, err := dydbSVC.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(dynamobdTable),
		Key: map[string]*dynamodb.AttributeValue{
			"Email": {
				S: aws.String(userEmail),
			},
		},
	})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Error verifying email (geting item from dynamo) -- " + err.Error(),
		})
		return
	}

	if len(result.Item) == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Error verifying email (item doesn't exist or expired) -- " + err.Error(),
		})
		return
	}

	emailToken := entity.EmailToken{}
	err = dynamodbattribute.UnmarshalMap(result.Item, &emailToken)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Error verifying email (unable to marshal item) -- " + err.Error(),
		})
		return
	}

	if emailToken.Token != userToken {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Error verifying email (Invalid token)",
		})
		return
	}

	if emailToken.TTL-180 < time.Now().Unix() {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Error verifying email (token expired!!)",
		})
		return
	}

	err = uc.userRepo.VerifyUserEmail(userEmail)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Error verifying email (unable to update db) -- " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{
		"message": "User verified successfully!!!",
	})
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
