package repo

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/ssukanmi/webservice/entity"
	"github.com/ssukanmi/webservice/service"
	"gorm.io/gorm"
)

var (
	s3BucketName = os.Getenv("S3_BUCKETNAME")
)

type UserRepository interface {
	InsertUser(user entity.User) (entity.User, error)
	FindByUsername(username string) (entity.User, error)
	UpdateUser(username string, user entity.User) (entity.User, error)
	GetUserProfilePic(username string) (entity.UserImage, error)
	UpdateUserProfilePic(username, filename string) (entity.UserImage, error)
	DeleteUserProfilePic(username string) error
}

type userRepo struct {
	connection *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Failed to load env file")
	}
	s3BucketName = os.Getenv("S3_BUCKETNAME")
	return &userRepo{
		connection: db,
	}
}

func (ur *userRepo) InsertUser(user entity.User) (entity.User, error) {
	user.Password = service.HashPassword(user.Password)
	res := ur.connection.Save(&user)
	return user, res.Error
}

func (ur *userRepo) FindByUsername(username string) (entity.User, error) {
	user := entity.User{}
	res := ur.connection.Where("username = ?", username).Take(&user)
	return user, res.Error
}

func (ur *userRepo) UpdateUser(username string, user entity.User) (entity.User, error) {
	currentUser, err := ur.FindByUsername(username)
	if err != nil {
		return currentUser, err
	}
	ur.connection.Model(&entity.User{}).Where("username = ?", username)
	if user.FirstName != "" {
		currentUser.FirstName = user.FirstName
	}
	if user.LastName != "" {
		currentUser.LastName = user.LastName
	}
	if user.Password != "" {
		currentUser.Password = service.HashPassword(user.Password)
	}
	res := ur.connection.Save(&currentUser)
	return currentUser, res.Error
}

func (ur *userRepo) UpdateUserProfilePic(username, filename string) (entity.UserImage, error) {
	user, err := ur.FindByUsername(username)
	userImage := entity.UserImage{}
	if err != nil {
		return userImage, err
	}
	res := ur.connection.Model(&entity.UserImage{}).Where("user_id = ?", user.ID).UpdateColumn("url", s3BucketName+"/"+username+"/"+filename).Take(&userImage)
	if res.Error != nil {
		if res.Error == gorm.ErrRecordNotFound {
			userImage.UserID = user.ID
			userImage.FileName = filename
			userImage.URL = s3BucketName + "/" + username + "/" + filename
			res = ur.connection.Create(&userImage)
			return userImage, res.Error
		}
		return userImage, res.Error
	}
	return userImage, res.Error
}

func (ur *userRepo) GetUserProfilePic(username string) (entity.UserImage, error) {
	user, err := ur.FindByUsername(username)
	userImage := entity.UserImage{}
	if err != nil {
		return userImage, err
	}
	res := ur.connection.Where("user_id = ?", user.ID).Take(&userImage)
	return userImage, res.Error
}

func (ur *userRepo) DeleteUserProfilePic(username string) error {
	user, err := ur.FindByUsername(username)
	if err != nil {
		return err
	}
	res := ur.connection.Where("user_id = ?", user.ID).Delete(&entity.UserImage{})
	return res.Error
}
