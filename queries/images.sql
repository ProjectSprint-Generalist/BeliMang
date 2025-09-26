-- name: InsertImage :exec
INSERT INTO images (
  object_key
) VALUES (
  $1
);

-- name: GetImage :one
SELECT * FROM images WHERE object_key = $1;