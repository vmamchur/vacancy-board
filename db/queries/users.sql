-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email, password)
VALUES (
	gen_random_uuid(),
	NOW(),
	NOW(),
	$1,
	$2
)
RETURNING *;

-- name: GetUserByEmail :one
SELECT *
FROM users
WHERE email = $1;

-- name: GetUser :one
SELECT *
FROM users
WHERE id = $1;

