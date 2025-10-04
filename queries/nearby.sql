-- name: GetNearbyMerchants :many
SELECT
  m.id,
  m.name,
  m.merchant_category,
  COALESCE(m.image_url, '') AS image_url,
  ST_Y(m.location::geometry) AS lat,
  ST_X(m.location::geometry) AS long,
  m.created_at,
  ST_DistanceSphere(m.location, ST_SetSRID(ST_MakePoint(sqlc.arg(long), sqlc.arg(lat)), 4326)) AS distance
FROM merchants m
WHERE
  (sqlc.narg(merchant_id)::text IS NULL OR m.id::text = sqlc.narg(merchant_id))
  AND (sqlc.narg(merchant_category)::text IS NULL OR m.merchant_category::text = sqlc.narg(merchant_category))
  AND (
    sqlc.narg(name)::text IS NULL
    OR LOWER(m.name) LIKE LOWER('%' || sqlc.narg(name) || '%')
    OR EXISTS (
      SELECT 1 FROM merchant_items mi
      WHERE mi.merchant_id = m.id
        AND LOWER(mi.name) LIKE LOWER('%' || sqlc.narg(name) || '%')
    )
  )
ORDER BY distance ASC, m.id ASC
LIMIT sqlc.arg(row_limit)::int OFFSET sqlc.arg(row_offset)::int;

-- name: CountNearbyMerchants :one
SELECT COUNT(*)
FROM merchants m
WHERE
  (sqlc.narg(merchant_id)::text IS NULL OR m.id::text = sqlc.narg(merchant_id))
  AND (sqlc.narg(merchant_category)::text IS NULL OR m.merchant_category::text = sqlc.narg(merchant_category))
  AND (
    sqlc.narg(name)::text IS NULL
    OR LOWER(m.name) LIKE LOWER('%' || sqlc.narg(name) || '%')
    OR EXISTS (
      SELECT 1 FROM merchant_items mi
      WHERE mi.merchant_id = m.id
        AND LOWER(mi.name) LIKE LOWER('%' || sqlc.narg(name) || '%')
    )
  );


