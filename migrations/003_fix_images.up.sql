-- Fix invalid FK on merchants.image_url and add missing merchant_items.image_url

-- Drop invalid foreign key constraint if it exists (it referenced a non-existent column)
ALTER TABLE merchants DROP CONSTRAINT IF EXISTS fk_merchants_images;

-- Rename merchants.image_url to merchants.url for consistency
DO $$
BEGIN
  IF EXISTS (
    SELECT 1 FROM information_schema.columns
    WHERE table_name = 'merchants' AND column_name = 'image_url'
  ) THEN
    ALTER TABLE merchants RENAME COLUMN image_url TO url;
  END IF;
END $$;

-- Ensure merchant_items has image_url column used by application code
ALTER TABLE merchant_items
  ADD COLUMN IF NOT EXISTS image_url TEXT NOT NULL DEFAULT '';

-- Create FK to images.url using merchants.url
DO $$
BEGIN
  IF EXISTS (
    SELECT 1 FROM information_schema.columns
    WHERE table_name = 'merchants' AND column_name = 'url'
  ) THEN
    ALTER TABLE merchants
      ADD CONSTRAINT fk_merchants_images_url FOREIGN KEY (url) REFERENCES images(url);
  END IF;
END $$;


