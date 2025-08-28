package handler

import (
	"homemie/internal/service"
	"homemie/models"
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

type CreateListingRequest struct {
	Title           string               `json:"title" binding:"required"`
	Description     string               `json:"description"`
	PropertyType    string               `json:"property_type"`
	IsShared        bool                 `json:"is_shared"`
	Price           float64              `json:"price" binding:"required"`
	AreaM2          float64              `json:"area_m2"`
	ContactPhone    string               `json:"contact_phone"`
	ContactEmail    string               `json:"contact_email"`
	ContactName     string               `json:"contact_name"`
	NumBedrooms     int                  `json:"num_bedrooms"`
	NumBathrooms    int                  `json:"num_bathrooms"`
	NumFloors       int                  `json:"num_floors"`
	HasBalcony      bool                 `json:"has_balcony"`
	HasParking      bool                 `json:"has_parking"`
	Amenities       []string             `json:"amenities"`
	PetAllowed      bool                 `json:"pet_allowed"`
	AllowedPetTypes []string             `json:"allowed_pet_types"`
	Latitude        float64              `json:"latitude"`
	Longitude       float64              `json:"longitude"`
	ListingType     string               `json:"listing_type"`
	DepositAmount   float64              `json:"deposit_amount"`
	Address         models.Address       `json:"address"`
	Images          []models.ListingImage `json:"images"`
}

func (h *ListingHandler) Create(c *gin.Context) {
	var req CreateListingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.GetUint("user_id")
	listing, err := h.svc.Create(service.CreateListingInput{
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
