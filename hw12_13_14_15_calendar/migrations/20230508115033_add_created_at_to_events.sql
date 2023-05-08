-- +goose Up
-- +goose StatementBegin
ALTER TABLE events ADD COLUMN created_at TIMESTAMP NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE events DROP COLUMN created_at;
-- +goose StatementEnd
