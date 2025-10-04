-- Consolidated migration: finalize image linkage and item image column

-- Drop any prior FKs to avoid conflicts
ALTER TABLE merchants DROP CONSTRAINT IF EXISTS fk_merchants_images;
ALTER TABLE merchants DROP CONSTRAINT IF EXISTS fk_merchants_images_url;

-- Ensure merchants has an 'image_url' column (original design)
ALTER TABLE merchants
  ADD COLUMN IF NOT EXISTS image_url TEXT NOT NULL;

-- Ensure merchant_items has image_url (used by app code)
ALTER TABLE merchant_items
  ADD COLUMN IF NOT EXISTS image_url TEXT NOT NULL DEFAULT '';

-- Intentionally no foreign key to images.url (not unique by design)