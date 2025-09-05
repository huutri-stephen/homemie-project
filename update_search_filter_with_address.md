# 📄 Prompt: Update Search & Filter API with Address Support  

## Context  
We already have a `listings` table that supports search and filtering by property type, price, bedrooms, etc.  
Now, we want to **extend the API** to support filtering by **address** (`city`, `ward`, `area`) using the `addresses` and `address_locations` tables.  

## Requirement  

### API Endpoint  
```http
GET /api/v1/listings
```

### Query Parameters  
- `keyword` → optional, search in `title` and `description` (full-text search).  
- `city_id` → optional, filter by city (only one city).  
- `ward_ids` → optional, multiple ward IDs (comma-separated).  
- `area_ids` → optional, multiple area IDs (comma-separated).  
- Other filters remain unchanged:  
  - `price_min`, `price_max`  
  - `area_min`, `area_max`  
  - `property_types[]`  
  - `listing_type`  
  - `num_bedrooms`, `num_bathrooms`, `num_floors`  
  - `has_balcony`, `has_parking`, `amenities[]`, `pet_allowed`, `allowed_pet_types[]`  

---

## Example Request  
```http
GET /api/v1/listings?keyword=studio&city_id=1&ward_ids=2,3&price_min=2000000&price_max=5000000
```

---

## SQL Example (PostgreSQL)  

```sql
WITH filtered_listings AS (
    SELECT l.*
    FROM listings l
    JOIN addresses a ON l.address_id = a.id
    WHERE l.status = 'active'

      -- Search keyword
      AND (
            :keyword IS NULL 
            OR to_tsvector('simple', l.title || ' ' || l.description) @@ plainto_tsquery('simple', :keyword)
          )

      -- City filter (optional)
      AND (
            :city_id IS NULL
            OR a.city_id = :city_id
          )

      -- Ward filter (optional, multi-choice)
      AND (
            :ward_ids IS NULL
            OR a.ward_id = ANY(string_to_array(:ward_ids, ',')::BIGINT[])
          )

      -- Area filter (optional, multi-choice)
      AND (
            :area_ids IS NULL
            OR a.area_id = ANY(string_to_array(:area_ids, ',')::BIGINT[])
          )

      -- Price range
      AND (
            (:price_min IS NULL OR l.price >= :price_min) AND
            (:price_max IS NULL OR l.price <= :price_max)
          )

      -- Bedrooms
      AND (:num_bedrooms IS NULL OR l.num_bedrooms = :num_bedrooms)

      -- Property type
      AND (:property_types IS NULL OR l.property_type = ANY(string_to_array(:property_types, ',')::property_type_enum[]))

      -- Listing type
      AND (:listing_type IS NULL OR l.listing_type = :listing_type)
)
SELECT *
FROM filtered_listings
ORDER BY created_at DESC
LIMIT :limit OFFSET :offset;
```

---

## Index Recommendation  
```sql
-- Address filter indexes
CREATE INDEX idx_addresses_city ON addresses(city_id);
CREATE INDEX idx_addresses_ward ON addresses(ward_id);
CREATE INDEX idx_addresses_area ON addresses(area_id);

-- Join optimization
CREATE INDEX idx_listings_address_id ON listings(address_id);
```

---

## Optional Supporting Endpoint  
For dynamic filtering UI:  
```http
GET /api/v1/locations?city_id=1
```
Returns all wards and areas under a given city.  
