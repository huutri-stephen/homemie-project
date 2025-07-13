package utils

import (
    // "errors"
    "os"
    "time"

    "github.com/golang-jwt/jwt/v4"
    "mihome/models"
)

var jwtSecret = []byte(getJWTSecret())

func getJWTSecret() string {
    secret := os.Getenv("JWT_SECRET")
    if secret == "" {
        secret = "default_secret" // fallback
    }
    return secret
}

type JWTClaims struct {
    UserID uint   `json:"user_id"`
    Email  string `json:"email"`
    Role   string `json:"role"`
    jwt.RegisteredClaims
}

// Tạo JWT
func GenerateJWT(user models.User) (string, error) {
    claims := JWTClaims{
        UserID: user.ID,
        Email:  user.Email,
        Role:   user.Role,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(72 * time.Hour)),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(jwtSecret)
}

// Giải mã token
func ParseJWT(tokenStr string) (*JWTClaims, error) {
    token, err := jwt.ParseWithClaims(tokenStr, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
        return jwtSecret, nil
    })

    if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
        return claims, nil
    }

    return nil, err
}
