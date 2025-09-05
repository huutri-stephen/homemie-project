# API Update: Search Listings with Address and Radius Filter

## 1. Description
This update enhances the **Search Listings API** by adding filtering capabilities:
- **Filter by Address**:
  - `city_id` (optional, choose 1 city).
  - `ward_ids` and `area_ids` (multiple selection, only available when `city_id` is provided).
- **Filter by Radius**:
  - Allow users to search for listings within a given distance (km) from a specific location (e.g., home, office).

## 2. New Query Parameters
```
GET /api/v1/listings/search
```

| Parameter    | Type      | Required | Description |
|--------------|-----------|----------|-------------|
| `city_id`    | bigint    | No       | ID of the city. |
| `ward_ids`   | array     | No       | List of ward IDs (must belong to city). |
| `area_ids`   | array     | No       | List of area IDs (must belong to city). |
| `radius_km`  | float     | No       | Radius in kilometers. |
| `lat`        | float     | Required if `radius_km` is provided | Latitude of reference point. |
| `lon`        | float     | Required if `radius_km` is provided | Longitude of reference point. |

## 3. SQL Implementation
```sql
SELECT l.id, l.title, l.price, a.latitude, a.longitude,
       6371 * acos(
           cos(radians(:lat)) * cos(radians(a.latitude)) * cos(radians(a.longitude) - radians(:lon))
           + sin(radians(:lat)) * sin(radians(a.latitude))
       ) AS distance_km
FROM listings l
JOIN addresses a ON l.address_id = a.id
WHERE l.status = 'active'

-- Filter by city
AND (:city_id IS NULL OR a.city_id = :city_id)

-- Filter by ward
AND (:ward_ids IS NULL OR a.ward_id = ANY(:ward_ids))

-- Filter by area
AND (:area_ids IS NULL OR a.area_id = ANY(:area_ids))

-- Filter by radius
AND (
    :radius_km IS NULL OR (
        6371 * acos(
            cos(radians(:lat)) * cos(radians(a.latitude)) * cos(radians(a.longitude) - radians(:lon))
            + sin(radians(:lat)) * sin(radians(a.latitude))
        )
    ) <= :radius_km
)
ORDER BY distance_km ASC;
```

## 4. Golang Example (Haversine formula)
```go
func Haversine(lat1, lon1, lat2, lon2 float64) float64 {
    const EarthRadius = 6371.0 // km
    dLat := (lat2 - lat1) * math.Pi / 180.0
    dLon := (lon2 - lon1) * math.Pi / 180.0

    a := math.Sin(dLat/2)*math.Sin(dLat/2) +
        math.Cos(lat1*math.Pi/180.0)*math.Cos(lat2*math.Pi/180.0)*
            math.Sin(dLon/2)*math.Sin(dLon/2)
    c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

    return EarthRadius * c
}
```

## 5. Example Request
```
GET /api/v1/listings/search?city_id=1&ward_ids=2,3&radius_km=5&lat=10.762622&lon=106.660172
```

## 6. Example Response
```json
{
  "success": true,
  "data": [
    {
      "id": 101,
      "title": "Căn hộ 2PN quận 1",
      "price": 15000000,
      "latitude": 10.762500,
      "longitude": 106.660300,
      "distance_km": 0.3
    },
    {
      "id": 102,
      "title": "Phòng trọ Quận 3",
      "price": 5000000,
      "latitude": 10.771000,
      "longitude": 106.670000,
      "distance_km": 1.2
    }
  ]
}
```
