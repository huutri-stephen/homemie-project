package repository

import (
	"github.com/google/uuid"
	"gorm.io/gorm"

	"homemie/internal/domain"
)

type favoriteRepository struct {
	db *gorm.DB
}

func NewFavoriteRepository(db *gorm.DB) *favoriteRepository {
	return &favoriteRepository{db}
}

func (r *favoriteRepository) Create(favorite *domain.Favorite) error {
	return r.db.Create(favorite).Error
}

func (r *favoriteRepository) GetByUserID(userID uuid.UUID, page int, pageSize int) ([]domain.Favorite, error) {
	var favorites []domain.Favorite
	offset := (page - 1) * pageSize
	err := r.db.Where("user_id = ?", userID).Offset(offset).Limit(pageSize).Find(&favorites).Error
	return favorites, err
}

func (r *favoriteRepository) FindByUserIDAndListingID(userID, listingID uuid.UUID) (*domain.Favorite, error) {
    var favorite domain.Favorite
    err := r.db.Where("user_id = ? AND listing_id = ?", userID, listingID).First(&favorite).Error
    if err != nil {
        if err == gorm.ErrRecordNotFound {
            return nil, nil
        }
        return nil, err
    }
    return &favorite, nil
}
