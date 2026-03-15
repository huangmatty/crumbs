-- name: CreateBuyer :one
INSERT INTO buyers (name, email, user_id)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetBuyers :many
SELECT * FROM buyers
WHERE deleted_at IS NULL
AND user_id = $1;

-- name: GetBuyerByID :one
SELECT * FROM buyers
WHERE id = $1;

-- name: GetUserIDForBuyer :one
SELECT user_id FROM buyers
WHERE id = $1;

-- name: UpdateBuyerName :one
UPDATE buyers
SET updated_at = NOW(), name = $1
WHERE id = $2
RETURNING *;

-- name: UpdateBuyerEmail :one
UPDATE buyers
SET updated_at = NOW(), email = $1
WHERE id = $2
RETURNING *;

-- name: SoftDeleteBuyer :one
UPDATE buyers
SET updated_at = NOW(), deleted_at = NOW()
WHERE id = $1
RETURNING *;

-- name: RestoreBuyer :one
UPDATE buyers
SET updated_at = NOW(), deleted_at = NULL
WHERE id = $1
RETURNING *;