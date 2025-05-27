-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS products.products (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    description TEXT NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS offers (
    id TEXT PRIMARY KEY,
    product_id TEXT NOT NULL,
    name TEXT NOT NULL UNIQUE,
    description TEXT NOT NULL UNIQUE,
    picture_encrypted TEXT NOT NULL UNIQUE,
    duration INTEGER NOT NULL,
    price INTEGER NOT NULL,
    price_id_encrypted TEXT NOT NULL UNIQUE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE products.products;
DROP TABLE products.offers;
-- +goose StatementEnd
