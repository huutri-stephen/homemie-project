package listing

import (
    "net/http"
    "strconv"

    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
    "mihome/models"
)

type Handler struct {
    DB *gorm.DB
}

func NewHandler(db *gorm.DB) *Handler {
    return &Handler{DB: db}
}

type CreateListingInput struct {
    Title       string  `json:"title" binding:"required"`
    Description string  `json:"description"`
    Price       float64 `json:"price" binding:"required"`
    Address     string  `json:"address"`
    City        string  `json:"city"`
}

func (h *Handler) CreateListing(c *gin.Context) {
    // Lấy user ID từ JWT context
    userIDVal, exists := c.Get("user_id")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
        return
    }

    userID, ok := userIDVal.(uint)
    if !ok {
        // Trường hợp convert lỗi
        switch v := userIDVal.(type) {
        case string:
            id, err := strconv.Atoi(v)
            if err == nil {
                userID = uint(id)
            } else {
                c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID"})
                return
            }
        default:
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user context"})
            return
        }
    }

    var input CreateListingInput
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

	listing := models.Listing{
		Title:       input.Title,
		Description: input.Description,
		Price:       input.Price,
		Address:     input.Address,
		City:        input.City,
		OwnerID:     userID,
	}

    if err := h.DB.Create(&listing).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create listing"})
        return
    }

    c.JSON(http.StatusCreated, gin.H{
        "message": "Listing created successfully",
        "listing": listing,
    })
}

func (h *Handler) GetAllListings(c *gin.Context) {
    var listings []models.Listing
    if err := h.DB.Preload("Owner").Find(&listings).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Không thể lấy danh sách"})
        return
    }
    c.JSON(http.StatusOK, gin.H{"data": listings})
}

func (h *Handler) GetListingByID(c *gin.Context) {
    id := c.Param("id")

    var listing models.Listing
    if err := h.DB.Preload("Owner").First(&listing, id).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Không tìm thấy bài đăng"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"data": listing})
}

func (h *Handler) UpdateListing(c *gin.Context) {
    id := c.Param("id")
    userID := c.GetUint("user_id")

    var listing models.Listing
    if err := h.DB.First(&listing, id).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Không tìm thấy bài đăng"})
        return
    }

    if listing.OwnerID != userID {
        c.JSON(http.StatusForbidden, gin.H{"error": "Không có quyền cập nhật bài đăng này"})
        return
    }

    var input CreateListingInput
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    listing.Title = input.Title
    listing.Description = input.Description
    listing.Price = input.Price
    listing.Address = input.Address
    listing.City = input.City

    if err := h.DB.Save(&listing).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Cập nhật thất bại"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Cập nhật thành công", "data": listing})
}

func (h *Handler) DeleteListing(c *gin.Context) {
    id := c.Param("id")
    userID := c.GetUint("user_id")

    var listing models.Listing
    if err := h.DB.First(&listing, id).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Không tìm thấy bài đăng"})
        return
    }

    if listing.OwnerID != userID {
        c.JSON(http.StatusForbidden, gin.H{"error": "Không có quyền xóa bài đăng này"})
        return
    }

    if err := h.DB.Delete(&listing).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Xóa thất bại"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Xóa bài đăng thành công"})
}

