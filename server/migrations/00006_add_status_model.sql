-- +goose Up
-- +goose StatementBegin
ALTER TABLE game 
ADD "status" smallint NOT NULL DEFAULT 0;

ALTER TABLE round 
ADD "status" smallint NOT NULL DEFAULT 0;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE game
DROP COLUMN "status";

ALTER TABLE round
DROP COLUMN "status";
-- +goose StatementEnd
