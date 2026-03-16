-- +goose Up
CREATE TABLE phones (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP,
    country_code CHAR(3) NOT NULL,
    number CHAR(20) NOT NULL,
    UNIQUE(country_code, number),
    buyer_id UUID REFERENCES buyers (id) ON DELETE CASCADE,
    talent_id UUID REFERENCES talents (id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE phones;