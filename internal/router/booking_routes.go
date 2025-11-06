package router

import (
	"database/sql"

	"github.com/gin-gonic/gin"

	"homemie/internal/handler"
	"homemie/internal/repo"
	"homemie/internal/service"
)

type ListingDeps struct {
	Handler *handler.ListingHandler
}

// buildListingDeps khởi tạo đầy đủ repository, service, handler cho Listing
func buildListingDeps(db *sql.DB) *ListingDeps {
	listingRepo := repo.NewListingRepo(db)
	addressRepo := repo.NewAddressRepository(db)
	listingImageRepo := repo.NewListingImageRepository(db)

	svc := service.NewListingService(
		listingRepo,
		addressRepo,
		listingImageRepo,
	)

	h := handler.NewListingHandler(svc)

	return &ListingDeps{Handler: h}
}

// InitBookingRoutes đăng ký các route liên quan đến booking
func InitBookingRoutes(rg *gin.RouterGroup, db *sql.DB) {
	bookingRepo := repo.NewBookingRepo(db)
	listingRepo := repo.NewListingRepo(db)
	svc := service.NewBookingService(bookingRepo, listingRepo)
	h := handler.NewBookingHandler(svc)

	bookings := rg.Group("/bookings")
	{
		bookings.POST("", h.CreateBooking)
		bookings.GET("/my-bookings", h.GetMyBookings)
		bookings.GET("/owner-bookings", h.GetOwnerBookings)

		bookings.PUT("/:id/respond", h.RespondToBooking)
		bookings.PUT("/:id/cancel", h.CancelBooking)
	}
}