package repo

import (
	"context"
	"fmt"
	"mime/multipart"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/joho/godotenv"
	"github.com/ssukanmi/webservice/entity"
	"github.com/ssukanmi/webservice/service"
	"gorm.io/gorm"
)

var (
	s3BucketName = os.Getenv("S3_BUCKETNAME")
	ctx          = context.Background()
)

type UserRepository interface {
	InsertUser(user entity.User) (entity.User, error)
	FindByUsername(username string) (entity.User, error)
	UpdateUser(username string, user entity.User) (entity.User, error)
	GetUserProfilePic(username string) (entity.UserImage, error)
	UpdateUserProfilePic(username string, file *multipart.FileHeader) (entity.UserImage, error)
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

func (ur *userRepo) UpdateUserProfilePic(username string, file *multipart.FileHeader) (entity.UserImage, error) {
	user, err := ur.FindByUsername(username)
	userImage := entity.UserImage{}
	if err != nil {
		return userImage, err
	}
	sess, err := session.NewSession()
	if err != nil {
		return userImage, err
	}
	uploader := s3manager.NewUploader(sess)
	f, err := file.Open()
	if err != nil {
		fmt.Println("Unable to open file -- " + err.Error())
	}
	defer f.Close()
	// res := ur.connection.Model(&entity.UserImage{}).Where("user_id = ?", user.ID).UpdateColumns(entity.UserImage{URL: s3BucketName + "/" + username + "/" + file.Filename, FileName: file.Filename}).Take(&userImage)
	res := ur.connection.Where("user_id = ?", user.ID).Take(&userImage)
	if res.Error != nil {
		if res.Error == gorm.ErrRecordNotFound {
			userImage.UserID = user.ID
			userImage.FileName = file.Filename
			userImage.URL = s3BucketName + "/" + username + "/" + file.Filename
			res = ur.connection.Create(&userImage)
			_, err = uploader.Upload(&s3manager.UploadInput{
				Bucket: aws.String(s3BucketName),
				Key:    aws.String(user.ID.String() + "/" + file.Filename),
				Body:   f,
			})
			if err != nil {
				return userImage, err
			}
			return userImage, res.Error
		}
		return userImage, res.Error
	}

	svc := s3.New(sess)
	if _, err := svc.DeleteObjectWithContext(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s3BucketName),
		Key:    aws.String(userImage.UserID.String() + "/" + userImage.FileName),
	}); err != nil {
		return userImage, err
	}
	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(s3BucketName),
		Key:    aws.String(user.ID.String() + "/" + file.Filename),
		Body:   f,
	})
	if err != nil {
		return userImage, err
	}

	userImage.FileName = file.Filename
	userImage.URL = s3BucketName + "/" + username + "/" + file.Filename
	ur.connection.Save(&userImage)

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
	userImage, err := ur.GetUserProfilePic(username)
	if err != nil {
		return err
	}
	res := ur.connection.Where("user_id = ?", user.ID).Delete(&entity.UserImage{})
	sess, err := session.NewSession()
	if err != nil {
		return err
	}
	svc := s3.New(sess)
	if _, err := svc.DeleteObjectWithContext(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s3BucketName),
		Key:    aws.String(userImage.UserID.String() + "/" + userImage.FileName),
	}); err != nil {
		return err
	}
	return res.Error
}
