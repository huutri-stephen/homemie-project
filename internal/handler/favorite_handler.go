package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"homemie/internal/service"
	"homemie/models/request"
	"homemie/models/response"
)

type FavoriteHandler struct {
	favoriteService service.FavoriteService
}

func NewFavoriteHandler(favoriteService service.FavoriteService) *FavoriteHandler {
	return &FavoriteHandler{favoriteService}
}

func (h *FavoriteHandler) Create(c *gin.Context) {
	var req request.Favorite
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request")
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	if err := h.favoriteService.Create(userID.(uuid.UUID), req); err != nil {
		response.Error(c, http.StatusInternalServerError, "failed to create favorite")
		return
	}

	response.Success(c, http.StatusCreated, "favorite created successfully", nil)
}

func (h *FavoriteHandler) GetByUserID(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	favorites, err := h.favoriteService.GetByUserID(userID.(uuid.UUID), page, pageSize)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "failed to get favorites")
		return
	}

	response.Success(c, http.StatusOK, "favorites retrieved successfully", favorites)
}
