-- +goose Up
CREATE TYPE phone AS (
    country_code CHAR(3),
    number CHAR(20)
);

-- +goose Down
DROP TYPE phone;