-- Consolidated migration: finalize image linkage and item image column

-- Drop any prior FKs to avoid conflicts
ALTER TABLE merchants DROP CONSTRAINT IF EXISTS fk_merchants_images;
ALTER TABLE merchants DROP CONSTRAINT IF EXISTS fk_merchants_images_url;

-- Ensure merchants has a 'url' column (rename from image_url if present, else add)
DO $$
BEGIN
  IF EXISTS (
    SELECT 1 FROM information_schema.columns
    WHERE table_name = 'merchants' AND column_name = 'image_url'
  ) THEN
    ALTER TABLE merchants RENAME COLUMN image_url TO url;
  ELSIF NOT EXISTS (
    SELECT 1 FROM information_schema.columns
    WHERE table_name = 'merchants' AND column_name = 'url'
  ) THEN
    ALTER TABLE merchants ADD COLUMN url TEXT NOT NULL;
  END IF;
END $$;

-- Ensure merchant_items has image_url (used by app code)
ALTER TABLE merchant_items
  ADD COLUMN IF NOT EXISTS image_url TEXT NOT NULL DEFAULT '';

-- Create FK from merchants.url to images.url
ALTER TABLE merchants
  ADD CONSTRAINT fk_merchants_images_url FOREIGN KEY (url) REFERENCES images(url);