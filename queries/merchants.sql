-- name: CreateMerchant :one
INSERT INTO merchants (
  name, merchant_category, url, location
) VALUES (
  $1, $2, $3, $4
) RETURNING id;
