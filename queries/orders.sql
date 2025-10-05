-- name: CreateCalculatedEstimate :one
INSERT INTO calculated_estimates (
  user_id, total_price, estimated_delivery_time_minutes, estimate_data
) VALUES (
  sqlc.arg(user_id)::uuid, sqlc.arg(total_price), sqlc.arg(estimated_delivery_time_minutes), sqlc.arg(estimate_data)
) RETURNING id;

-- name: GetCalculatedEstimateByID :one
SELECT 
  id,
  user_id,
  total_price,
  estimated_delivery_time_minutes,
  estimate_data,
  created_at
FROM calculated_estimates
WHERE id = sqlc.arg(id)::uuid;

-- name: CreateOrder :one
INSERT INTO orders (
  user_id, calculated_estimate_id
) VALUES (
  sqlc.arg(user_id)::uuid, sqlc.arg(calculated_estimate_id)::uuid
) RETURNING id;

-- name: GetOrdersByUserID :many
SELECT 
  o.id,
  o.created_at,
  ce.estimate_data
FROM orders o
JOIN calculated_estimates ce ON o.calculated_estimate_id = ce.id
WHERE o.user_id = $1::uuid
ORDER BY o.created_at DESC
LIMIT $2 OFFSET $3;

-- name: GetOrdersCountByUserID :one
SELECT COUNT(*)
FROM orders o
WHERE o.user_id = $1::uuid;