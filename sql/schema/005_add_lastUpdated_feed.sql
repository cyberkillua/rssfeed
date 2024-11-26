-- +goose Up
ALTER TABLE feeds ADD COLUMN last_updated TIMESTAMP;

-- +goose Down
ALTER TABLE feeds DROP COLUMN last_updated;