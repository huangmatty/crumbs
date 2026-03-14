-- name: CreateTalent :one
INSERT INTO talents (name, email, user_id)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetTalents :many
SELECT * FROM talents
WHERE deleted_at IS NULL
AND user_id = $1
ORDER BY name;

-- name: GetTalentByID :one
SELECT * FROM talents
WHERE id = $1;

-- name: UpdateTalentName :one
UPDATE talents
SET updated_at = NOW(), name = $1
WHERE id = $2
RETURNING *;

-- name: UpdateTalentEmail :one
UPDATE talents
SET updated_at = NOW(), email = $1
WHERE id = $2
RETURNING *;

-- name: SoftDeleteTalent :one
UPDATE talents
SET updated_at = NOW(), deleted_at = NOW()
WHERE id = $1
RETURNING *;