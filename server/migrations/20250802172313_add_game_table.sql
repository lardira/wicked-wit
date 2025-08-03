-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS users (
    id uuid PRIMARY KEY,
    username VARCHAR(32) NOT NULL UNIQUE,
    password VARCHAR(32) NOT NULL,
    profile_img text NOT NULL,
    created_at timestamp NOT NULL DEFAULT NOW(),
    updated_at timestamp NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS round (
    id SERIAL PRIMARY KEY,
    presenter_id uuid REFERENCES users(id),
    winner_id uuid REFERENCES users(id),
    position smallint NOT NULL,
    created_at timestamp NOT NULL DEFAULT NOW(),
    updated_at timestamp NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS game (
    id uuid PRIMARY KEY,
    title text NOT NULL,
    max_players smallint NOT NULL,
    max_round smallint NOT NULL,
    current_round int NOT NULL REFERENCES round(id),
    created_at timestamp NOT NULL DEFAULT NOW(),
    updated_at timestamp NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS games_users (
    game_id uuid NOT NULL REFERENCES game(id),
    user_id uuid NOT NULL REFERENCES users(id),
    CONSTRAINT game_user_pk PRIMARY KEY (game_id, user_id)
);

CREATE TABLE IF NOT EXISTS user_answer (
    id serial PRIMARY KEY, 
    user_id uuid NOT NULL REFERENCES users(id),
    round_id int NOT NULL REFERENCES round(id),
    answer text NOT NULL,
    position smallint NOT NULL DEFAULT 1,
    created_at timestamp NOT NULL DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE game, round, users, games_users, user_answer CASCADE;
-- +goose StatementEnd
