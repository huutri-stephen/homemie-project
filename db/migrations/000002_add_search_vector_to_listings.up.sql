-- Drop index
DROP INDEX IF EXISTS idx_listings_search_vector;

-- Drop trigger
DROP TRIGGER IF EXISTS tsvectorupdate ON listings;

-- Drop trigger function
DROP FUNCTION IF EXISTS listings_search_trigger;

-- Drop column
ALTER TABLE listings DROP COLUMN IF EXISTS search_vector;
