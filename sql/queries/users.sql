-- name: CreateUser :one
INSERT INTO users (username, email, hashed_password)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetUserByUsername :one
SELECT * FROM users
WHERE username = $1;

-- name: GetUserFromRefreshToken :one
SELECT * FROM users
WHERE id IN (
    SELECT user_id FROM refresh_tokens
    WHERE token = $1
);