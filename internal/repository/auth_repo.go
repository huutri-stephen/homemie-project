package repository

import (
	"homemie/internal/domain"
	"homemie/models/dto"
	"strings"

	"gorm.io/gorm"
)

type authRepo struct {
	db *gorm.DB
}

func NewAuthRepo(db *gorm.DB) domain.AuthRepository {
	return &authRepo{db: db}
}

func (r *authRepo) CreateUser(user *dto.User) error {
	return r.db.Create(user).Error
}

func (r *authRepo) GetUserByEmail(email string) (*dto.User, error) {
	var user dto.User
	if err := r.db.Where("email = ?", strings.ToLower(email)).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *authRepo) GetUserByID(id int64) (*dto.User, error) {
	var user dto.User
	if err := r.db.Where("id = ?", id).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *authRepo) CreateToken(token *dto.Token) error {
	return r.db.Create(token).Error
}

func (r *authRepo) GetToken(token string, userID int64, tokenType dto.TokenType) (*dto.Token, error) {
	var t dto.Token
	if err := r.db.Where("token = ? AND user_id = ? AND token_type = ?", token, userID, tokenType).First(&t).Error; err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *authRepo) UpdateUser(user *dto.User) error {
	return r.db.Save(user).Error
}

func (r *authRepo) DeleteToken(token *dto.Token) error {
	return r.db.Delete(token).Error
}
