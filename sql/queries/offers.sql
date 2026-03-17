-- name: CreateOffer :one
INSERT INTO offers (start_date, end_date, buyer_id, user_id)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetOffers :many
SELECT * FROM offers
WHERE deleted_at IS NULL
AND user_id = $1;

-- name: GetOfferByID :one
SELECT * FROM offers
WHERE id = $1;

-- name: GetUserIDForOffer :one
SELECT user_id FROM offers
WHERE id = $1;

-- name: UpdateOfferStartDate :one
UPDATE offers
SET updated_at = NOW(), start_date = $1
WHERE id = $2
RETURNING *;

-- name: UpdateOfferEndDate :one
UPDATE offers
SET updated_at = NOW(), end_date = $1
WHERE id = $2
RETURNING *;

-- name: SoftDeleteOffer :one
UPDATE offers
SET updated_at = NOW(), deleted_at = NOW()
WHERE id = $1
RETURNING *;

-- name: RestoreOffer :one
UPDATE offers
SET updated_at = NOW(), deleted_at = NULL
WHERE id = $1
RETURNING *;