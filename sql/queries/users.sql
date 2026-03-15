-- name: CreateUser :one
INSERT INTO users (username, email, hashed_password)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetUserByID :one
SELECT * FROM users
WHERE id = $1;

-- name: GetUserByUsername :one
SELECT * FROM users
WHERE username = $1;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1;

-- name: GetUserFromRefreshToken :one
SELECT * FROM users
WHERE id IN (
    SELECT user_id FROM refresh_tokens
    WHERE token = $1
);

-- name: UpdateUsername :one
UPDATE users
SET updated_at = NOW(), username = $1
WHERE id = $2
RETURNING *;

-- name: UpdateUserEmail :one
UPDATE users
SET updated_at = NOW(), email = $1
WHERE id = $2
RETURNING *;

-- name: UpdateUserPassword :one
UPDATE users
SET updated_at = NOW(), hashed_password = $1
WHERE id = $2
RETURNING *;