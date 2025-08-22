-- +goose Up
-- +goose StatementBegin

ALTER TABLE round 
    ADD game_id uuid NOT NULL,
    ADD CONSTRAINT round_game_fk FOREIGN KEY (game_id) REFERENCES game(id);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

ALTER TABLE round
    DROP COLUMN game_id;
    
-- +goose StatementEnd
