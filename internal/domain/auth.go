package domain

import "homemie/models/dto"

type AuthRepository interface {
	CreateUser(user *dto.User) error
	GetUserByEmail(email string) (*dto.User, error)
	GetUserByID(id int64) (*dto.User, error)
	CreateToken(token *dto.Token) error
	GetToken(token string, userID int64, tokenType dto.TokenType) (*dto.Token, error)
	UpdateUser(user *dto.User) error
	DeleteToken(token *dto.Token) error
}
