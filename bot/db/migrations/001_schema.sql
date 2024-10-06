-- +goose Up
CREATE TABLE leaders
(
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    civ_name TEXT NOT NULL,
    leader_name TEXT NOT NULL,
    discord_emoji_string TEXT,
    banned BOOLEAN NOT NULL DEFAULT FALSE,
    tier FLOAT NOT NULL
);

CREATE UNIQUE INDEX leaders_civ_name_leader_name_uindex ON leaders (civ_name, leader_name);

CREATE TABLE drafts
(
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    active BOOLEAN NOT NULL DEFAULT TRUE,
    players TEXT -- String array, comma separated.
);

CREATE TABLE offered
(
    player TEXT NOT NULL,
    draft_id INTEGER NOT NULL,
    offered TEXT,
    PRIMARY KEY (player, draft_id)
);

-- +goose Down
DROP TABLE leaders;
DROP TABLE drafts;
DROP TABLE offered;