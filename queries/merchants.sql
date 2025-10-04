-- name: CreateMerchant :one
INSERT INTO merchants (
  name, merchant_category, image_url, location
) VALUES (
  $1, $2, $3, ST_SetSRID(ST_MakePoint($4, $5), 4326)
) RETURNING id;
