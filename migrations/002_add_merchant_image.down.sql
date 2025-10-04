-- Revert consolidated changes
ALTER TABLE merchants DROP CONSTRAINT IF EXISTS fk_merchants_images_url;

DO $$
BEGIN
  IF EXISTS (
    SELECT 1 FROM information_schema.columns
    WHERE table_name = 'merchants' AND column_name = 'url'
  ) THEN
    ALTER TABLE merchants RENAME COLUMN url TO image_url;
  END IF;
END $$;

ALTER TABLE merchant_items DROP COLUMN IF EXISTS image_url;