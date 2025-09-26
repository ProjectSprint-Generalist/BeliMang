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
