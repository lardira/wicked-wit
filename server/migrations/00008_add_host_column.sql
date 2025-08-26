-- +goose Up
-- +goose StatementBegin
ALTER TABLE
    game
ADD
    COLUMN user_host_id uuid NOT NULL,
ADD
    FOREIGN KEY (user_host_id) REFERENCES users(id);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
ALTER TABLE
    game DROP COLUMN host_user_id;

-- +goose StatementEnd