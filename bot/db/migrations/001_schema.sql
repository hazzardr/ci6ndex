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
    active BOOLEAN NOT NULL DEFAULT TRUE
);

CREATE TABLE draft_registry
(
    player_id INTEGER NOT NULL,
    draft_id INTEGER NOT NULL,
    primary key (player_id, draft_id),
    FOREIGN KEY (player_id) REFERENCES players (id),
    FOREIGN KEY (draft_id) REFERENCES drafts (id)
);

CREATE TABLE pool
(
    player_id INTEGER NOT NULL,
    draft_id INTEGER NOT NULL,
    leader INTEGER NOT NULL,
    FOREIGN KEY (player_id) REFERENCES players (id),
    FOREIGN KEY (draft_id) REFERENCES drafts (id),
    FOREIGN KEY (leader) REFERENCES leaders (id)
);

CREATE TABLE picks
(
    player TEXT NOT NULL,
    draft_id INTEGER NOT NULL,
    pick INTEGER NOT NULL,
    offered TEXT, -- comma-delimited list of leader ids
    PRIMARY KEY (player, draft_id),
    FOREIGN KEY (pick) REFERENCES leaders (id),
    FOREIGN KEY (draft_id) REFERENCES drafts (id)
);

CREATE TABLE players
(
    id INTEGER PRIMARY KEY,
    username TEXT NOT NULL,
    global_name TEXT,
    discord_avatar TEXT
);

-- +goose Down
DROP TABLE leaders;
DROP TABLE drafts;
DROP TABLE pool;
DROP TABLE picks;
DROP TABLE players;
DROP TABLE draft_registry;