package domain

import "homemie/models"

type AuthRepository interface {
	CreateUser(user *models.User) error
	GetUserByEmail(email string) (*models.User, error)
	GetUserByID(id uint) (*models.User, error)
	CreateToken(token *models.Token) error
	GetToken(token string, userID uint, tokenType models.TokenType) (*models.Token, error)
	UpdateUser(user *models.User) error
	DeleteToken(token *models.Token) error
}
