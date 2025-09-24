package router

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"homemie/internal/handler"
	"homemie/internal/repo"
	"homemie/internal/service"
)

type ListingDeps struct {
	Handler *handler.ListingHandler
}

// buildListingDeps khởi tạo đầy đủ repository, service, handler cho Listing
func buildListingDeps(db *gorm.DB, logger *zap.Logger) *ListingDeps {
	listingRepo := repo.NewListingRepo(db, logger.Named("listing_repo"))
	addressRepo := repo.NewAddressRepository(db, logger.Named("address_repo"))
	listingImageRepo := repo.NewListingImageRepository(db, logger.Named("listing_image_repo"))

	svc := service.NewListingService(
		listingRepo,
		addressRepo,
		listingImageRepo,
		logger.Named("listing_service"),
	)

	h := handler.NewListingHandler(svc, logger.Named("listing_handler"))

	return &ListingDeps{Handler: h}
}

// InitBookingRoutes đăng ký các route liên quan đến booking
func InitBookingRoutes(rg *gin.RouterGroup, db *gorm.DB, logger *zap.Logger) {
	bookingRepo := repo.NewBookingRepo(db, logger.Named("booking_repo"))
	listingRepo := repo.NewListingRepo(db, logger.Named("listing_repo"))
	svc := service.NewBookingService(bookingRepo, listingRepo, logger.Named("booking_service"))
	h := handler.NewBookingHandler(svc, logger.Named("booking_handler"))

	bookings := rg.Group("/bookings")
	{
		bookings.POST("", h.CreateBooking)
		bookings.GET("/my-bookings", h.GetMyBookings)
		bookings.GET("/owner-bookings", h.GetOwnerBookings)

		bookings.PUT("/:id/respond", h.RespondToBooking)
		bookings.PUT("/:id/cancel", h.CancelBooking)
	}
}
