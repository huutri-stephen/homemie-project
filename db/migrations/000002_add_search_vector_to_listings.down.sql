-- Add search_vector column
ALTER TABLE listings ADD COLUMN search_vector tsvector;

-- Initialize existing data
UPDATE listings 
SET search_vector = to_tsvector('simple', coalesce(title,'') || ' ' || coalesce(description,''));

-- Create trigger function to auto-update search_vector
CREATE FUNCTION listings_search_trigger() RETURNS trigger AS $$
BEGIN
  NEW.search_vector :=
    to_tsvector('simple', coalesce(NEW.title,'') || ' ' || coalesce(NEW.description,''));
  RETURN NEW;
END
$$ LANGUAGE plpgsql;

-- Attach trigger to listings
CREATE TRIGGER tsvectorupdate
BEFORE INSERT OR UPDATE ON listings
FOR EACH ROW EXECUTE FUNCTION listings_search_trigger();

-- Create GIN index for fast search
CREATE INDEX idx_listings_search_vector ON listings USING gin(search_vector);
