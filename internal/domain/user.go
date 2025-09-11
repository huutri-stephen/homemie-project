package domain

import "homemie/models/dto"

type UserRepository interface {
	UpdateUser(user *dto.User) error
	CreateUser(user *dto.User) error
	GetUserByEmail(email string) (*dto.User, error)
	GetUserByID(id int64) (*dto.User, error)
}
