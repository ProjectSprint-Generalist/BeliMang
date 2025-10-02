-- name: CreateMerchant :exec
INSERT INTO merchants (
  name, merchant_category, image_url, location
) VALUES (
  $1, $2, $3, $4
);
