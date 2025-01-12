-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, username, hashed_password)
VALUES(gen_random_uuid(), NOW(), NOW(), $1, $2)
RETURNING *;

-- name: GetUserByID :one
SELECT *
FROM users
WHERE id = $1;

-- name: GetUserByUsername :one
SELECT * FROM users
WHERE username =$1;