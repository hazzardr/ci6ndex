-- +goose Up
CREATE TABLE rerolls
(
    player_id INTEGER NOT NULL
);

-- +goose Down
DROP TABLE rerolls;