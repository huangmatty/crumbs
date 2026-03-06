-- name: CreateTalent :one
INSERT INTO talents (name, user_id)
VALUES ($1, $2)
RETURNING *;

-- name: GetTalents :many
SELECT * FROM talents
ORDER BY name;

-- name: GetTalentByID :one
SELECT * FROM talents
WHERE id = $1;