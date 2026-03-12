-- name: CreateTalent :one
INSERT INTO talents (name, user_id)
VALUES ($1, $2)
RETURNING *;

-- name: GetTalents :many
SELECT * FROM talents
WHERE deleted_at IS NULL
AND user_id = $1
ORDER BY name;

-- name: GetTalentByID :one
SELECT * FROM talents
WHERE id = $1;

-- name: SoftDeleteTalent :one
UPDATE talents
SET deleted_at = NOW()
WHERE id = $1
RETURNING *;