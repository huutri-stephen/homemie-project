ALTER TABLE listings
    ALTER COLUMN amenities TYPE JSONB USING amenities::JSONB,
    ALTER COLUMN allowed_pet_types TYPE JSONB USING allowed_pet_types::JSONB;
