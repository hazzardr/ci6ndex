-- +goose Up
CREATE TABLE leaders
(
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    civ_name TEXT NOT NULL,
    leader_name TEXT NOT NULL,
    discord_emoji_string TEXT,
    banned BOOLEAN NOT NULL DEFAULT FALSE,
    ranking FLOAT NOT NULL
);

CREATE UNIQUE INDEX leaders_civ_name_leader_name_uindex ON leaders (civ_name, leader_name);

CREATE TABLE draft_strategies
(
    name TEXT UNIQUE NOT NULL,
    description TEXT NOT NULL,
    pool_size INTEGER NOT NULL DEFAULT 3,
    randomize BOOLEAN NOT NULL DEFAULT FALSE,
    PRIMARY KEY (name)
);

CREATE TABLE drafts
(
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    draft_strategy TEXT NOT NULL,
    active BOOLEAN NOT NULL DEFAULT TRUE,
    players TEXT,
    FOREIGN KEY (draft_strategy) REFERENCES draft_strategies (name)
);

CREATE TABLE offered
(
    user_id INTEGER NOT NULL,
    draft_id INTEGER NOT NULL,
    offered TEXT,
    PRIMARY KEY (user_id, draft_id)
);

-- +goose Down
DROP TABLE leaders;
DROP TABLE draft_strategies;
DROP TABLE drafts;
DROP TABLE offered;