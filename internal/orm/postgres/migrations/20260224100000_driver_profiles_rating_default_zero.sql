-- +goose Up
-- +goose StatementBegin
ALTER TABLE driver_profiles ALTER COLUMN rating SET DEFAULT 0;
UPDATE driver_profiles SET rating = 0 WHERE rating = 5.0;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE driver_profiles ALTER COLUMN rating SET DEFAULT 5.0;
-- +goose StatementEnd
