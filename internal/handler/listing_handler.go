package handler

import (
	"homemie/internal/service"
	"homemie/models/dto"
	"homemie/models/request"
	"homemie/models/response"
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
		c.JSON(http.StatusBadRequest, response.BaseResponse{
			Success: false,
			Error:   err.Error(),
		})
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
		c.JSON(http.StatusInternalServerError, response.BaseResponse{
			Success: false,
			Error:   "Create listing failed",
		})
		return
	}

	c.JSON(http.StatusCreated, response.BaseResponse{Success: true, Data: listing})
}

func (h *ListingHandler) SearchAndFilter(c *gin.Context) {
	var filter dto.SearchFilterListing
	if err := c.ShouldBindQuery(&filter); err != nil {
		c.JSON(http.StatusBadRequest, response.BaseResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	listings, pagination, err := h.svc.SearchAndFilter(&filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.BaseResponse{
			Success: false,
			Error:   "Get list failed",
		})
		return
	}
	c.JSON(http.StatusOK, response.BaseResponse{Success: true, Data: gin.H{
		"listings":   listings,
		"pagination": pagination,
	}})
}

func (h *ListingHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.Atoi(idStr)

	listing, err := h.svc.GetByID(int64(id))
	if err != nil {
		c.JSON(http.StatusNotFound, response.BaseResponse{
			Success: false,
			Error:   "Not found",
		})
		return
	}
	c.JSON(http.StatusOK, response.BaseResponse{
		Success: true,
		Data:    listing,
	})
}

func (h *ListingHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.Atoi(idStr)
	userID := c.GetInt64("user_id")

	var req request.CreateListingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.BaseResponse{
			Success: false,
			Error:   err.Error(),
		})
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
			c.JSON(http.StatusForbidden, response.BaseResponse{
				Success: false,
				Error:   "Unauthorized to edit this listing",
			})
		} else {
			c.JSON(http.StatusInternalServerError, response.BaseResponse{
				Success: false,
				Error:   "Update failed",
			})
		}
		return
	}

	c.JSON(http.StatusOK, response.BaseResponse{
		Success: true,
		Data:    listing,
	})
}

func (h *ListingHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.Atoi(idStr)
	userID := c.GetInt64("user_id")

	err := h.svc.Delete(int64(id), userID)
	if err != nil {
		if err.Error() == "unauthorized" {
			c.JSON(http.StatusForbidden, response.BaseResponse{
				Success: false,
				Error:   "Unauthorized to delete this listing",
			})
		} else {
			c.JSON(http.StatusInternalServerError, response.BaseResponse{
				Success: false,
				Error:   "Delete failed",
			})
		}
		return
	}

	c.JSON(http.StatusOK, response.BaseResponse{
		Success: true,
		Message: "Delete listing successfully",
	})
}