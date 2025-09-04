# Hướng dẫn chạy với Docker

## Yêu cầu
- Cài đặt [Docker](https://docs.docker.com/get-docker/)  
- Cài đặt [Docker Compose](https://docs.docker.com/compose/install/)  

## Các lệnh Docker Compose

### Run Docker in WSL (without Docker Desktop)
```bash
sudo dockerd
```

### Build và chạy container
```bash
sudo docker compose up --build -d
```

### Dừng và xóa container
```bash
sudo docker compose down -v
```

### Tạm dừng container
```bash
sudo docker compose stop
```

### Khởi động lại container đã dừng
```bash
sudo docker compose start
```

---

## Ghi chú
Analysis:

The current codebase provides a solid foundation for a booking/rental platform. However, to be considered a complete MVP, several key features are missing. The focus should be on
enhancing the user experience, providing more control over bookings, and adding essential functionalities like user profiles and search.

Recommended Features for MVP:

Here are the features I recommend adding to make this a complete MVP, in order of priority:

1. User Profile Management:
    * Get User Profile: An endpoint to retrieve the current user's profile information.
    * Update User Profile: An endpoint for users to update their profile information (name, bio, avatar, etc.).

2. Search and Filtering for Listings:
    * Basic Search: A simple search functionality based on keywords (e.g., title, description).
    * Advanced Filtering: Allow users to filter listings by various criteria like city, price range, number of bedrooms, etc.

3. Booking Management:
    * Approve/Reject Bookings: Allow property owners to approve or reject booking requests.
    * Cancel Bookings: Allow users (both renters and owners) to cancel bookings.

4. Favorites/Wishlist:
    * Add to Favorites: Allow users to mark listings as favorites.
    * Get Favorites: An endpoint to retrieve the user's list of favorite listings.

5. Reviews and Ratings:
    * Leave a Review: Allow users to leave reviews and ratings for listings they have booked.
    * Get Reviews: An endpoint to retrieve all reviews for a specific listing.

These features will provide a much more complete and user-friendly experience, which is crucial for an MVP.

Would you like me to start implementing any of these features? I can begin with the user profile management, which is a fundamental feature. I will start by creating the necessary files
and then implementing the endpoints.