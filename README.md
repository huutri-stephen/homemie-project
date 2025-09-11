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


 I've finished integrating MinIO and implementing presigned URL generation. I added AWS SDK dependencies, configured MinIO, and created the necessary S3 infrastructure,
  services, handlers, and routes. The main application now initializes the S3 client on startup. The feature is ready to use. I'll now notify the user.I've finished
  integrating MinIO and implementing presigned URL generation. I added AWS SDK dependencies, configured MinIO, and created the necessary S3 infrastructure, services,
  handlers, and routes. The main application now initializes the S3 client on startup. The feature is ready to use. I'll now notify the user.I've finished integrating
  MinIO and implementing pre-signed URL generation. The backend now has a /api/v1/media/presigned-url endpoint for file uploads. I'll now notify the user and provide
  usage instructions.I have completed the implementation for integrating MinIO and providing a pre-signed URL generation endpoint.

  Here is a summary of the work done:

   1. Dependencies: Added the aws-sdk-go-v2 library to your project.
   2. Configuration: Updated your .env file with MinIO credentials and created a corresponding configuration struct in config/config.go.
   3. S3 Infrastructure: Created internal/infra/s3.go which handles the S3 client initialization and automatically creates the homemie-media bucket on application start
      if it doesn't exist.
   4. Application Startup: Modified cmd/main.go to initialize the S3 client.
   5. Presigned URL Endpoint:
       * Created a new service (internal/service/media_service.go), handler (internal/handler/media_handler.go), and router (internal/router/media_routes.go).
       * This exposes a new protected endpoint: POST /api/v1/media/presigned-url.

  How to Use the New Endpoint

  To generate a pre-signed URL for a file upload, you need to make a POST request to /api/v1/media/presigned-url with a valid JWT token and the following JSON body:

   1 {
   2   "objectName": "user-avatars/some-unique-filename.jpg",
   3   "bucketName": "homemie-media"
   4 }

  The response will be a JSON object containing the pre-signed URL:

   1 {
   2   "url": "http://localhost:9000/homemie-media/user-avatars/some-unique-filename.jpg?X-Amz-Algorithm=..."
   3 }

  Your frontend can then use this URL to upload the file directly to MinIO using an HTTP PUT request with the file content as the request body.

  This setup allows you to develop locally using MinIO. When you are ready to move to AWS S3, you will only need to update the S3_* environment variables in your .env
  file to point to your S3 bucket and credentials. No code changes will be required.