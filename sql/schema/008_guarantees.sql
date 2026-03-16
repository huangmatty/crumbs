-- +goose Up
CREATE TABLE guarantees (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP,
    amount INTEGER NOT NULL,
    currency_alpha CHAR(3) NOT NULL,
    due_date DATE,
    offer_id UUID NOT NULL REFERENCES offers (id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE guarantees;