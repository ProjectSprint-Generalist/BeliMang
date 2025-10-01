-- name: CreateAdmin :exec
INSERT INTO users (
  username, password, email, role
) VALUES (
  $1, $2, $3, 'admin'
);

-- name: CreateUser :exec
INSERT INTO users (
  username, password, email, role
) VALUES (
  $1, $2, $3, 'user'
);

-- name: GetAdminByUsername :one
SELECT * FROM users where username = $1 AND role = 'admin';

-- name: GetUserByUsername :one
SELECT * FROM users where username = $1 AND role = 'user';