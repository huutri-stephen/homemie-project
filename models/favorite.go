package models

type Favorite struct {
    ID        uint `gorm:"primaryKey"`
    UserID    uint `gorm:"not null"`
    ListingID uint `gorm:"not null"`

    User    User    `gorm:"foreignKey:UserID"`
    Listing Listing `gorm:"foreignKey:ListingID"`
}
