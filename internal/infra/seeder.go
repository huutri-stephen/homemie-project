package infra

import (
    "log"
    "time"

    "homemie/models"
    "gorm.io/gorm"
)

func SeedData(db *gorm.DB) {
    log.Println("Seeding mock data...")

    // Check nếu có user rồi thì không seed lại
    var count int64
    db.Model(&models.User{}).Count(&count)
    if count > 0 {
        log.Println("Data already exists, skip seeding.")
        return
    }

    // 1. Tạo users
    users := []models.User{
        {Name: "Alice", Email: "alice@example.com", Password: "hashed123", Phone: "0912345678", Role: "owner"},
        {Name: "Bob", Email: "bob@example.com", Password: "hashed123", Phone: "0987654321", Role: "renter"},
    }
    db.Create(&users)

    // 2. Tạo listing
    listing := models.Listing{
        Title:       "Phòng trọ gần ĐH Bách Khoa",
        Description: "Phòng sạch sẽ, có máy lạnh, wifi miễn phí",
        Price:       2_500_000,
        Address:     "123 Nguyễn Văn Cừ, Quận 5, TP.HCM",
        City:        "Hồ Chí Minh",
        OwnerID:     users[0].ID,
    }
    db.Create(&listing)

    // 3. Ảnh bài đăng
    images := []models.ListingImage{
        {ListingID: listing.ID, URL: "https://example.com/images/phong1.jpg"},
        {ListingID: listing.ID, URL: "https://example.com/images/phong2.jpg"},
    }
    db.Create(&images)

    // 4. Yêu thích
    fav := models.Favorite{
        UserID:    users[1].ID,
        ListingID: listing.ID,
    }
    db.Create(&fav)

    // // 5. Booking
    // booking := models.Booking{
    //     UserID:    users[1].ID,
    //     ListingID: listing.ID,
    //     StartDate: time.Now().Add(48 * time.Hour),
    //     EndDate:   time.Now().Add(72 * time.Hour),
    //     Status:    "pending",
    // }
    // db.Create(&booking)


    // 6. Tin nhắn
    msg := models.Message{
        SenderID:   users[1].ID,
        ReceiverID: users[0].ID,
        Content:    "Phòng còn trống không ạ?",
        SentAt:     time.Now(),
    }
    db.Create(&msg)

    log.Println("✅ Seeded successfully.")
}
