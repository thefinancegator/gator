-- +goose Up
ALTER TABLE feeds
ADD COLUMN last_fetched_At TIMESTAMP;

-- +goose Down
ALTER TABLE feeds
DROP COLUMN last_fetched_At;