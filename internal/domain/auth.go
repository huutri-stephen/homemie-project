package domain

import "homemie/models/dto"

type AuthRepository interface {
	CreateToken(token *dto.Token) error
	GetToken(token string, userID int64, tokenType dto.TokenType) (*dto.Token, error)
	DeleteToken(token *dto.Token) error
}
