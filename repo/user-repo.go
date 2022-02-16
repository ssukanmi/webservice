package repo

import (
	"github.com/ssukanmi/webservice/entity"
	"github.com/ssukanmi/webservice/service"
	"gorm.io/gorm"
)

type UserRepository interface {
	InsertUser(user entity.User) (entity.User, error)
	FindByUsername(username string) (entity.User, error)
	UpdateUser(username string, user entity.User) (entity.User, error)
}

type userRepo struct {
	connection *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
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
	currentUser, _ := ur.FindByUsername(username)
	currentUser.FirstName = user.FirstName
	currentUser.LastName = user.LastName
	currentUser.Username = user.Username
	res := ur.connection.Save(&currentUser)
	return user, res.Error
}
