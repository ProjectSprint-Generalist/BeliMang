-- name: CreateMerchant :one
INSERT INTO merchants (
  name, merchant_category, image_url, location
) VALUES (
  $1, $2, $3, ST_SetSRID(ST_MakePoint($4, $5), 4326)
) RETURNING id;

-- name: GetMerchants :many
SELECT
  m.id,
  m.name,
  m.merchant_category,
  COALESCE(m.image_url, '') as image_url,
  ST_Y(m.location::geometry) as lat,
  ST_X(m.location::geometry) as long,
  m.created_at
FROM merchants m
WHERE
  (sqlc.narg(merchant_id)::text IS NULL OR m.id::text = sqlc.narg(merchant_id))
  AND (sqlc.narg(merchant_category)::text IS NULL OR m.merchant_category::text = sqlc.narg(merchant_category))
  AND (
    sqlc.narg(name)::text IS NULL
    OR LOWER(m.name) LIKE LOWER('%' || sqlc.narg(name) || '%')
  )
ORDER BY
  CASE WHEN sqlc.arg(created_at) = 'asc' THEN m.created_at END ASC,
  CASE WHEN sqlc.arg(created_at) = 'desc' THEN m.created_at END DESC,
  m.id ASC
LIMIT sqlc.arg(limit_val)::int OFFSET sqlc.arg(offset_val)::int;

-- name: CountMerchants :one
SELECT COUNT(*)
FROM merchants m
WHERE
  (sqlc.narg(merchant_id)::text IS NULL OR m.id::text = sqlc.narg(merchant_id))
  AND (sqlc.narg(merchant_category)::text IS NULL OR m.merchant_category::text = sqlc.narg(merchant_category))
  AND (
    sqlc.narg(name)::text IS NULL
    OR LOWER(m.name) LIKE LOWER('%' || sqlc.narg(name) || '%')
  );

-- name: CreateMerchantItem :one
INSERT INTO merchant_items (
  merchant_id, name, product_category, price, image_url
) VALUES (
  $1, $2, $3, $4, $5
) RETURNING id;

-- name: GetMerchantByID :one
SELECT EXISTS(SELECT 1 FROM merchants WHERE id = $1);

-- name: GetMerchantItems :many
SELECT
  mi.id,
  mi.name,
  mi.product_category,
  mi.price,
  COALESCE(mi.image_url, '') as image_url,
  mi.created_at
FROM merchant_items mi
WHERE mi.merchant_id = sqlc.arg(merchant_id)
  AND (sqlc.narg(item_id)::text IS NULL OR mi.id::text = sqlc.narg(item_id))
  AND (sqlc.narg(product_category)::text IS NULL OR mi.product_category::text = sqlc.narg(product_category))
  AND (
    sqlc.narg(name)::text IS NULL
    OR LOWER(mi.name) LIKE LOWER('%' || sqlc.narg(name) || '%')
  )
ORDER BY
  CASE WHEN sqlc.arg(created_at) = 'asc' THEN mi.created_at END ASC,
  CASE WHEN sqlc.arg(created_at) = 'desc' THEN mi.created_at END DESC,
  mi.id ASC
LIMIT sqlc.arg(limit_val)::int OFFSET sqlc.arg(offset_val)::int;

-- name: CountMerchantItems :one
SELECT COUNT(*)
FROM merchant_items mi
WHERE mi.merchant_id = sqlc.arg(merchant_id)
  AND (sqlc.narg(item_id)::text IS NULL OR mi.id::text = sqlc.narg(item_id))
  AND (sqlc.narg(product_category)::text IS NULL OR mi.product_category::text = sqlc.narg(product_category))
  AND (
    sqlc.narg(name)::text IS NULL
    OR LOWER(mi.name) LIKE LOWER('%' || sqlc.narg(name) || '%')
  );

-- name: GetMerchantDetailsByID :one
SELECT
  m.id,
  m.name,
  m.merchant_category,
  COALESCE(m.image_url, '') as image_url,
  ST_Y(m.location::geometry) as lat,
  ST_X(m.location::geometry) as long,
  m.created_at
FROM merchants m
WHERE m.id = sqlc.arg(id)::uuid;

-- name: GetMerchantItemByID :one
SELECT
  mi.id,
  mi.name,
  mi.product_category,
  mi.price,
  COALESCE(mi.image_url, '') as image_url,
  mi.created_at
FROM merchant_items mi
WHERE mi.id = sqlc.arg(id)::uuid;
