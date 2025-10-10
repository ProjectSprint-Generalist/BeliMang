-- name: CreateUser :exec
INSERT INTO users (
  username, password, email, role
) VALUES (
  $1, $2, $3, $4
);

-- name: GetUserByUsername :one
SELECT * FROM users where username = $1 AND role = $2;