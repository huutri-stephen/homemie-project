# Search & Filtering Listings API Design

## Database Schema (PostgreSQL)
```sql
CREATE TYPE property_type_enum AS ENUM('rented_room', 'apartment', 'house', 'villa', 'dormitory', 'office', 'store', 'warehouse', 'land', 'all');
CREATE TYPE listing_type_enum AS ENUM('for_rent', 'for_sale');
CREATE TYPE listing_status_enum AS ENUM('active', 'inactive', 'pending', 'rejected', 'deleted');

CREATE TABLE listings (
    id                BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    owner_id          BIGINT REFERENCES users(id),
    title             VARCHAR(255) NOT NULL,
    description       TEXT,
    property_type     property_type_enum DEFAULT 'rented_room',
    is_shared         BOOLEAN DEFAULT FALSE,
    price             DECIMAL(15,2) NOT NULL,
    area_m2           DECIMAL(10,2),
    address_id        BIGINT REFERENCES addresses(id) NOT NULL,
    contact_phone     VARCHAR(20),
    contact_email     VARCHAR(255),
    contact_name      VARCHAR(100),
    num_bedrooms      INT,
    num_bathrooms     INT,
    num_floors        INT,
    has_balcony       BOOLEAN DEFAULT FALSE,
    has_parking       BOOLEAN DEFAULT FALSE,
    amenities         JSONB,
    pet_allowed       BOOLEAN DEFAULT FALSE,
    allowed_pet_types JSONB,
    listing_type      listing_type_enum DEFAULT 'for_rent',
    deposit_amount    DECIMAL(15,2),
    status            listing_status_enum DEFAULT 'pending',
    is_featured       BOOLEAN DEFAULT FALSE,
    view_count        INT DEFAULT 0,
    created_at        TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at        TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    published_at      TIMESTAMP,
    expires_at        TIMESTAMP
);
```

## Search Requirement
- Input: `keyword` (string).
- Behavior: search in `title` OR `description` using PostgreSQL full-text search or ILIKE.

## Filter Requirement
- property_type: multichoice
- is_shared: boolean
- price: range (min_price, max_price)
- area_m2: range (min_area, max_area)
- num_bathrooms: number
- num_bedrooms: number
- num_floors: number
- has_balcony: boolean
- has_parking: boolean
- amenities: multichoice
- pet_allowed: boolean
- allowed_pet_types: multichoice
- listing_type: one of [for_rent, for_sale]
- status: only 'active'

## API Design
- Endpoint: `GET /listings`
- Query parameters:
  ```
  keyword=...
  property_type=apartment,house
  is_shared=true
  min_price=1000000
  max_price=5000000
  min_area=20
  max_area=50
  num_bedrooms=2
  num_bathrooms=1
  num_floors=3
  has_balcony=true
  has_parking=false
  amenities=wifi,elevator
  pet_allowed=true
  allowed_pet_types=dog,cat
  listing_type=for_rent
  page=1
  limit=20
  ```

## Response Format
```json
{
  "success": true,
  "data": {
    "listings": [
      {
        "id": 101,
        "title": "Cozy Apartment in Hanoi",
        "description": "Fully furnished, near metro",
        "property_type": "apartment",
        "is_shared": false,
        "price": 4500000,
        "area_m2": 35.5,
        "num_bedrooms": 2,
        "num_bathrooms": 1,
        "num_floors": 1,
        "has_balcony": true,
        "has_parking": true,
        "amenities": ["wifi", "elevator"],
        "pet_allowed": true,
        "allowed_pet_types": ["dog", "cat"],
        "listing_type": "for_rent",
        "status": "active",
        "created_at": "2025-09-01T10:00:00Z"
      }
    ],
    "pagination": {
      "page": 1,
      "limit": 20,
      "total": 120,
      "total_pages": 6
    }
  }
}
```

## SQL Query Example
- Add index:
```sql
ALTER TABLE listings 
ADD COLUMN search_vector tsvector 
    GENERATED ALWAYS AS (
        to_tsvector('simple', coalesce(title,'') || ' ' || coalesce(description,''))
    ) STORED;

CREATE INDEX idx_listings_search_vector ON listings USING GIN (search_vector);

-- Full-text search index
CREATE INDEX idx_listings_search_vector ON listings USING GIN (search_vector);

-- JSONB multichoice filter index
CREATE INDEX idx_listings_amenities ON listings USING GIN (amenities);
CREATE INDEX idx_listings_allowed_pet_types ON listings USING GIN (allowed_pet_types);

-- B-tree indexes cho filter chính
CREATE INDEX idx_listings_price ON listings (price);
CREATE INDEX idx_listings_area_m2 ON listings (area_m2);
CREATE INDEX idx_listings_num_bedrooms ON listings (num_bedrooms);
CREATE INDEX idx_listings_property_type ON listings (property_type);
CREATE INDEX idx_listings_listing_type ON listings (listing_type);
CREATE INDEX idx_listings_status ON listings (status);
```

- Sql query
```sql
SELECT *
FROM listings
WHERE status = 'active'
  -- full-text search
  AND (:keyword IS NULL OR search_vector @@ plainto_tsquery('simple', :keyword))
  
  -- property_type multichoice
  AND (:property_types IS NULL OR property_type = ANY(:property_types))
  
  -- boolean filters
  AND (:is_shared IS NULL OR is_shared = :is_shared)
  AND (:has_balcony IS NULL OR has_balcony = :has_balcony)
  AND (:has_parking IS NULL OR has_parking = :has_parking)
  AND (:pet_allowed IS NULL OR pet_allowed = :pet_allowed)
  
  -- range filters
  AND (:min_price IS NULL OR price >= :min_price)
  AND (:max_price IS NULL OR price <= :max_price)
  AND (:min_area IS NULL OR area_m2 >= :min_area)
  AND (:max_area IS NULL OR area_m2 <= :max_area)
  
  -- exact number filters
  AND (:num_bedrooms IS NULL OR num_bedrooms = :num_bedrooms)
  AND (:num_bathrooms IS NULL OR num_bathrooms = :num_bathrooms)
  AND (:num_floors IS NULL OR num_floors = :num_floors)
  
  -- enum filter
  AND (:listing_type IS NULL OR listing_type = :listing_type)
  
  -- JSONB multichoice filters
  AND (:amenities IS NULL OR amenities ?| array[:amenities])
  AND (:allowed_pet_types IS NULL OR allowed_pet_types ?| array[:allowed_pet_types])
  
ORDER BY created_at DESC
LIMIT :limit OFFSET (:page - 1) * :limit;
```
