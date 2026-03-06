-- name: CreateTalent :one
INSERT INTO talents (name, user_id)
VALUES ($1, $2)
RETURNING *;