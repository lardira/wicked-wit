-- +goose Up
-- +goose StatementBegin
ALTER TABLE game 
    ALTER COLUMN current_round DROP NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE game 
    ALTER COLUMN current_round SET NOT NULL;

-- +goose StatementEnd
