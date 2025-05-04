-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS events (
    id TEXT NOT NULL PRIMARY KEY,
    title TEXT NOT NULL,
    description TEXT NOT NULL,
    postal_code_encrypted TEXT NOT NULL,
    city_encrypted TEXT NOT NULL,
    address1_encrypted TEXT NOT NULL,
    address2_encrypted TEXT,
    placecount INTEGER NOT NULL,
    freeplace INTEGER CHECK(freeplace >= 0),
    begin_at TEXT NOT NULL,
    end_at TEXT NOT NULL,
    price_id_encrypted TEXT NOT NULL UNIQUE,
    day INTEGER NOT NULL CHECK(day > 0 AND day < 31) ,
    month INTEGER NOT NULL CHECK(month > 0 AND month < 13) ,
    year INTEGER NOT NULL,
    UNIQUE(day, month, year)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE events;
-- +goose StatementEnd
