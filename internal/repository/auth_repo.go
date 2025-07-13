package repository

import (
	"mihome/internal/domain"
	"mihome/models"

	"gorm.io/gorm"
)

type authRepo struct {
	db *gorm.DB
}

func NewAuthRepo(db *gorm.DB) domain.AuthRepository {
	return &authRepo{db: db}
}

func (r *authRepo) CreateUser(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *authRepo) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
