package auth

import (
    "net/http"
    "strings"
    // "time"

    "github.com/gin-gonic/gin"
    "golang.org/x/crypto/bcrypt"
    "gorm.io/gorm"
    "mihome/models"
)

type AuthHandler struct {
    DB *gorm.DB
}

func NewAuthHandler(db *gorm.DB) *AuthHandler {
    return &AuthHandler{DB: db}
}

type SignUpInput struct {
    Name     string `json:"name" binding:"required"`
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required,min=6"`
    Phone    string `json:"phone"`
}

func (h *AuthHandler) SignUp(c *gin.Context) {
    var input SignUpInput
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error hashing password"})
        return
    }

    user := models.User{
        Name:     input.Name,
        Email:    strings.ToLower(input.Email),
        Password: string(hashedPassword),
        Phone:    input.Phone,
        Role:     "renter", // default role
    }

    if err := h.DB.Create(&user).Error; err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Email đã tồn tại hoặc lỗi DB"})
        return
    }

    c.JSON(http.StatusCreated, gin.H{"message": "Đăng ký thành công"})
}

type LoginInput struct {
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required"`
}

func (h *AuthHandler) Login(c *gin.Context) {
    var input LoginInput
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    var user models.User
    if err := h.DB.Where("email = ?", strings.ToLower(input.Email)).First(&user).Error; err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Tài khoản không tồn tại"})
        return
    }

    if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Mật khẩu không đúng"})
        return
    }

    // Tạo token JWT
    token, err := GenerateJWT(user)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Không thể tạo token"})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "token": token,
        "user": gin.H{
            "id":    user.ID,
            "name":  user.Name,
            "email": user.Email,
            "role":  user.Role,
        },
    })
}
