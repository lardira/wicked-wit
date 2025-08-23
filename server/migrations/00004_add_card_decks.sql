-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS template_card (
    id serial PRIMARY KEY,
    placeholders_count int NOT NULL,
    text TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS answer_card (
    id serial PRIMARY KEY,
    text TEXT NOT NULL
);

ALTER TABLE
    round DROP COLUMN presenter_id,
ADD
    template_card_id int NOT NULL,
ADD
    CONSTRAINT template_card_fk FOREIGN KEY (template_card_id) REFERENCES template_card(id);

DROP TABLE user_answer;

CREATE TABLE IF NOT EXISTS user_answer (
    id serial PRIMARY KEY,
    user_id uuid NOT NULL REFERENCES users(id),
    round_id int NOT NULL REFERENCES round(id),
    votes smallint NOT NULL DEFAULT 0,
    created_at timestamp NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS game_used_card (
    id serial PRIMARY KEY,
    game_id uuid NOT NULL REFERENCES game(id),
    answer_card_id int NOT NULL REFERENCES answer_card(id),
    status smallint NOT NULL DEFAULT 0,
    user_id uuid NOT NULL REFERENCES users(id)
);

CREATE TABLE IF NOT EXISTS game_user_played_card (
    user_answer_id int NOT NULL REFERENCES user_answer(id),
    game_used_card_id int NOT NULL REFERENCES game_used_card(id),
    placeholder_index smallint NOT NULL DEFAULT 0,
    CONSTRAINT user_played_card_pk PRIMARY KEY (game_used_card_id, user_answer_id)
);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE game_user_played_card,
game_used_card,
user_answer;

CREATE TABLE IF NOT EXISTS user_answer (
    id serial PRIMARY KEY,
    user_id uuid NOT NULL REFERENCES users(id),
    round_id int NOT NULL REFERENCES round(id),
    answer text NOT NULL,
    position smallint NOT NULL DEFAULT 1,
    created_at timestamp NOT NULL DEFAULT NOW()
);

ALTER TABLE
    round DROP COLUMN template_card_id;

DROP TABLE template_card,
answer_card;

-- +goose StatementEnd