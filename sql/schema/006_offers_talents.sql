-- +goose Up
CREATE TABLE offers_talents (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP,
    offer_id UUID NOT NULL REFERENCES offers (id) ON DELETE CASCADE,
    talent_id UUID NOT NULL REFERENCES talents (id) ON DELETE CASCADE,
    UNIQUE(offer_id, talent_id)
);

-- +goose Down
DROP TABLE offers_talents;