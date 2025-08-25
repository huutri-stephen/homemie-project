package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"homemie/internal/service"
)

type ListingHandler struct {
	svc *service.ListingService
}

func NewListingHandler(svc *service.ListingService) *ListingHandler {
	return &ListingHandler{svc}
}

type CreateListingRequest struct {
	Title       string  `json:"title" binding:"required"`
	Description string  `json:"description"`
	Price       float64 `json:"price" binding:"required"`
	Address     string  `json:"address"`
	City        string  `json:"city"`
}

func (h *ListingHandler) Create(c *gin.Context) {
	var req CreateListingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.GetUint("user_id")
	listing, err := h.svc.Create(service.CreateListingInput{
		Title:       req.Title,
		Description: req.Description,
		Price:       req.Price,
		Address:     req.Address,
		City:        req.City,
		OwnerID:     userID,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Tạo bài đăng thất bại"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": listing})
}

func (h *ListingHandler) GetAll(c *gin.Context) {
	listings, err := h.svc.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Lỗi lấy danh sách"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": listings})
}

func (h *ListingHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.Atoi(idStr)

	listing, err := h.svc.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Không tìm thấy"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": listing})
}

func (h *ListingHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.Atoi(idStr)
	userID := c.GetUint("user_id")

	var req CreateListingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	listing, err := h.svc.Update(uint(id), userID, service.CreateListingInput{
		Title:       req.Title,
		Description: req.Description,
		Price:       req.Price,
		Address:     req.Address,
		City:        req.City,
	})
	if err != nil {
		if err.Error() == "unauthorized" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Không có quyền sửa bài đăng này"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Cập nhật thất bại"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": listing})
}

func (h *ListingHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.Atoi(idStr)
	userID := c.GetUint("user_id")

	err := h.svc.Delete(uint(id), userID)
	if err != nil {
		if err.Error() == "unauthorized" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Không có quyền xóa bài đăng này"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Xóa thất bại"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Xóa bài đăng thành công"})
}
