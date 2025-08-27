package models

import (
	"time"
)

type TokenType string

const (
	EmailVerification TokenType = "email_verification"
	PasswordReset     TokenType = "password_reset"
)

type Token struct {
	ID        uint      `gorm:"primaryKey"`
	UserID    uint      `gorm:"not null"`
	User      User      `gorm:"foreignKey:UserID"`
	TokenType TokenType `gorm:"type:token_type_enum;not null"`
	Token     string    `gorm:"type:varchar(255);not null"`
	ExpiresAt time.Time `gorm:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}