package domain

import "github.com/google/uuid"
import "homemie/models/dto"

type UserRepository interface {
	GetUserByID(id uuid.UUID) (*dto.User, error)
	UpdateUser(user *dto.User) error
}
