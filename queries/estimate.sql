-- name: GetMerchantLocationByID :one
SELECT 
  id::text AS id,
  ST_Y(location::geometry)::float8 AS lat,
  ST_X(location::geometry)::float8 AS long
FROM merchants
WHERE id = ($1)::text::uuid;

-- name: GetMerchantItemPriceByID :one
SELECT 
  id::text AS id,
  price::int4 AS price
FROM merchant_items
WHERE id = ($1)::text::uuid;