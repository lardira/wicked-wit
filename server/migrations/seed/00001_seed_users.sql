-- +goose Up
-- +goose StatementBegin

INSERT INTO users (id, username, password, profile_img) VALUES (
    'c5eedc3c-0e51-4cb8-bfdd-a64babc67725',
    'admin',
    'admin',
    'https://www.meme-arsenal.com/memes/4c8f80315713d6ede660cc124e3f67d5.jpg'
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DELETE FROM users u WHERE u.id = 'c5eedc3c-0e51-4cb8-bfdd-a64babc67725';

-- +goose StatementEnd
