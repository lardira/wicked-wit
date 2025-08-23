-- +goose Up
-- +goose StatementBegin
ALTER TABLE
    users
ALTER COLUMN
    profile_img DROP NOT NULL;

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
ALTER TABLE
    users
ALTER COLUMN
    profile_img
SET
    NOT NULL;

-- +goose StatementEnd