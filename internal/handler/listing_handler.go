package handler

import (
	"homemie/internal/service"
	"homemie/models/request"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ListingHandler struct {
	svc *service.ListingService
}

func NewListingHandler(svc *service.ListingService) *ListingHandler {
	return &ListingHandler{svc}
}

func (h *ListingHandler) Create(c *gin.Context) {
	var req request.CreateListingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.GetInt64("user_id")
	listing, err := h.svc.Create(request.CreateListingRequest{
		OwnerID:         userID,
		Title:           req.Title,
		Description:     req.Description,
		PropertyType:    req.PropertyType,
		IsShared:        req.IsShared,
		Price:           req.Price,
		AreaM2:          req.AreaM2,
		ContactPhone:    req.ContactPhone,
		ContactEmail:    req.ContactEmail,
		ContactName:     req.ContactName,
		NumBedrooms:     req.NumBedrooms,
		NumBathrooms:    req.NumBathrooms,
		NumFloors:       req.NumFloors,
		HasBalcony:      req.HasBalcony,
		HasParking:      req.HasParking,
		Amenities:       req.Amenities,
		PetAllowed:      req.PetAllowed,
		AllowedPetTypes: req.AllowedPetTypes,
		Latitude:        req.Latitude,
		Longitude:       req.Longitude,
		ListingType:     req.ListingType,
		DepositAmount:   req.DepositAmount,
		Address:         req.Address,
		Images:          req.Images,
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

	listing, err := h.svc.GetByID(int64(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Không tìm thấy"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": listing})
}

func (h *ListingHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.Atoi(idStr)
	userID := c.GetInt64("user_id")

	var req request.CreateListingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	listing, err := h.svc.Update(int64(id), userID, request.CreateListingRequest{
		Title:       req.Title,
		Description: req.Description,
		Price:       req.Price,
		// todo: add more param
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
	userID := c.GetInt64("user_id")

	err := h.svc.Delete(int64(id), userID)
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
