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


Example curl Requests

  Below are the curl commands to test the new booking management endpoints. Please replace {AUTH_TOKEN}, {BOOKING_ID}, and other placeholders with actual values from your
  system.

  1. Approve/Reject Booking (Owner)

  This endpoint is for the property owner to accept or reject a booking request.

  Approve a Booking:

   1 curl -X PUT http://localhost:8080/api/v1/bookings/{BOOKING_ID}/respond \
   2 -H "Authorization: Bearer {OWNER_AUTH_TOKEN}" \
   3 -H "Content-Type: application/json" \
   4 -d '{
   5     "status": "accepted",
   6     "response_message_from_owner": "Your booking is confirmed. I look forward to welcoming you."
   7 }'

  Reject a Booking:

   1 curl -X PUT http://localhost:8080/api/v1/bookings/{BOOKING_ID}/respond \
   2 -H "Authorization: Bearer {OWNER_AUTH_TOKEN}" \
   3 -H "Content-Type: application/json" \
   4 -d '{
   5     "status": "rejected",
   6     "response_message_from_owner": "Sorry, the property is not available at the requested time."
   7 }'

  2. Cancel Booking (Renter or Owner)

  This endpoint can be used by either the renter or the owner to cancel a booking.

  Cancel as Renter:

   1 curl -X PUT http://localhost:8080/api/v1/bookings/{BOOKING_ID}/cancel \
   2 -H "Authorization: Bearer {RENTER_AUTH_TOKEN}"

  Cancel as Owner:

   1 curl -X PUT http://localhost:8080/api/v1/bookings/{BOOKING_ID}/cancel \
   2 -H "Authorization: Bearer {OWNER_AUTH_TOKEN}"
