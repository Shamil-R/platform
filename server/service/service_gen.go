// Code generated by platform, DO NOT EDIT.

package service

import (
	model "gitlab/nefco/platform/server/model"
)

type Service interface {
	CreatePost(data model.PostCreateInput) (model.Post, error)
	UpdatePost(data model.PostUpdateInput, where model.PostWhereUniqueInput) (*model.Post, error)
	DeletePost(where model.PostWhereUniqueInput) (*model.Post, error)
	CreateUser(data model.UserCreateInput) (model.User, error)
	UpdateUser(data model.UserUpdateInput, where model.UserWhereUniqueInput) (*model.User, error)
	DeleteUser(where model.UserWhereUniqueInput) (*model.User, error)
}
