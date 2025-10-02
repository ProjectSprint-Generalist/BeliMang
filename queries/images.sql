-- name: CreateImage :one
INSERT INTO images (id, filename, url, size_bytes)
VALUES ($1, $2, $3, $4)
RETURNING id, filename, url, size_bytes, created_at;

-- name: GetImage :one
SELECT id, filename, url, size_bytes, created_at FROM images WHERE id = $1;