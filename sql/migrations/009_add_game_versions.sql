-- +goose Up
CREATE TABLE game_versions
(
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT,
    current BOOLEAN NOT NULL DEFAULT FALSE
);

-- +goose Down
DROP TABLE IF EXISTS game_versions;
