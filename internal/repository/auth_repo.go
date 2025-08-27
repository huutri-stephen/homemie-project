package repository

import (
	"homemie/internal/domain"
	"homemie/models"
	"strings"

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
	if err := r.db.Where("email = ?", strings.ToLower(email)).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *authRepo) GetUserByID(id uint) (*models.User, error) {
	var user models.User
	if err := r.db.Where("id = ?", id).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *authRepo) CreateToken(token *models.Token) error {
	return r.db.Create(token).Error
}

func (r *authRepo) GetToken(token string, userID uint, tokenType models.TokenType) (*models.Token, error) {
	var t models.Token
	if err := r.db.Where("token = ? AND user_id = ? AND token_type = ?", token, userID, tokenType).First(&t).Error; err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *authRepo) UpdateUser(user *models.User) error {
	return r.db.Save(user).Error
}

func (r *authRepo) DeleteToken(token *models.Token) error {
	return r.db.Delete(token).Error
}
