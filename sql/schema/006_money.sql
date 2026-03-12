-- +goose Up
CREATE TYPE money AS (
    amount INTEGER,
    currency_alpha CHAR(3)
);

-- +goose Down
DROP TYPE money;