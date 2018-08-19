package services

import (
	"gitlab/nefco/platform/models"
	"gitlab/nefco/platform/stores"
)

type UserService interface {
	CreateUser(data models.UserCreateInput) (*models.User, error)
	GetUsers(where *models.UserWhereInput) []*models.User
}

type userService struct {
	store stores.UserStore
}

func (s *userService) CreateUser(data models.UserCreateInput) (*models.User, error) {
	panic("not implemented")
}

func (s *userService) GetUsers(where *models.UserWhereInput) []*models.User {
	panic("not implemented")
}
