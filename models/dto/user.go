package dto

import (
	"time"
)

type User struct {
	ID                     int64 `gorm:"primaryKey"`
	FirstName              string `gorm:"type:varchar(50)"`
	LastName               string `gorm:"type:varchar(50)"`
	Name                   string `gorm:"type:varchar(100)"`
	Email                  string `gorm:"uniqueIndex;type:varchar(100);not null"`
	Phone                  string `gorm:"type:varchar(20)"`
	DateOfBirth            *time.Time
	Gender                 string `gorm:"type:gender_enum;default:'other'"`
	AvatarURL              string `gorm:"type:varchar(255)"`
	Bio                    string `gorm:"type:text"`
	PasswordHash           string `gorm:"type:text;not null"`
	Salt                   string `gorm:"type:varchar(64)"`
	Status                 string `gorm:"type:user_status_enum;default:'inactive'"`
	EmailVerifiedAt        *time.Time
	LastLoginAt            *time.Time
	ResetPasswordToken     string `gorm:"type:varchar(255)"`
	ResetPasswordExpiresAt *time.Time
	UserType               string `gorm:"type:user_type_enum;default:'renter'"`
	IdentityType           string `gorm:"type:identity_type_enum;default:'personal'"`
	CompanyName            string `gorm:"type:varchar(255)"`
	BusinessLicenseNumber  string `gorm:"type:varchar(100)"`
	AgentLicenseNumber     string `gorm:"type:varchar(100)"`
	VerifiedOwner          bool   `gorm:"default:false"`
	Role                   string `gorm:"type:role_enum;default:'user'"`
	CreatedAt              time.Time
	UpdatedAt              time.Time
}