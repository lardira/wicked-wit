-- +goose Up
-- +goose StatementBegin
ALTER TABLE
    game DROP COLUMN current_round;

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
ALTER TABLE
    game
ADD
    COLUMN current_round int NOT NULL
ADD
    CONSTRAINT current_round_fk FOREIGN KEY (current_round) REFERENCES round(id);

-- +goose StatementEnd