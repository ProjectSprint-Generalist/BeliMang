ALTER TABLE merchants ADD COLUMN image_url TEXT NOT NULL;
ALTER TABLE merchants ADD CONSTRAINT fk_merchants_images FOREIGN KEY (image_url) REFERENCES images(object_key);