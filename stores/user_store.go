package stores

import "gitlab/nefco/platform/models"

type UserStore interface {
	CreateUser(user *models.User) error
	UpdateUser(user *models.User) error
	DeleteUser(id int) error
	GetUserByID(id int) (*models.User, error)
	GetUsers() ([]models.User, error)
}
