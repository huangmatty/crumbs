-- +goose Up
CREATE TABLE offers_buyers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    offer_id UUID REFERENCES offers (id) ON DELETE CASCADE,
    buyer_id UUID REFERENCES buyers (id) ON DELETE CASCADE,
    UNIQUE(offer_id, buyer_id)
);

-- +goose Down
DROP TABLE offers_buyers;