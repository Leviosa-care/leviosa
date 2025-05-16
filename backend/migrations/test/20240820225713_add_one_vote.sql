-- +goose Up
-- +goose StatementBegin
INSERT INTO votes ( user_id, days, month, year)
    VALUES ( '123e4567-e89b-12d3-a456-426614174000', '23|12|6', 4, 2025);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM votes
WHERE
  rowid = (SELECT
  MAX(rowid)
FROM
  votes);
-- +goose StatementEnd
