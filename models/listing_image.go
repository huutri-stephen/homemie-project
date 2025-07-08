package models

type ListingImage struct {
    ID        uint   `gorm:"primaryKey"`
    ListingID uint   `gorm:"not null"`
    URL       string `gorm:"type:text;not null"`

    Listing Listing `gorm:"foreignKey:ListingID"`
}
