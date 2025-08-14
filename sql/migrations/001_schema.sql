-- +goose Up
CREATE TABLE leaders
(
    id INTEGER PRIMARY KEY, -- Use database-native auto-increment
    civ_name TEXT NOT NULL,
    leader_name TEXT NOT NULL,
    discord_emoji_string TEXT,
    banned BOOLEAN NOT NULL DEFAULT FALSE,
    tier FLOAT NOT NULL
);

CREATE UNIQUE INDEX leaders_civ_name_leader_name_uindex ON leaders (civ_name, leader_name);

CREATE TABLE drafts
(
    id INTEGER PRIMARY KEY,
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
    PRIMARY KEY (player_id, draft_id, leader), -- Add composite primary key
    FOREIGN KEY (player_id) REFERENCES players (id),
    FOREIGN KEY (draft_id) REFERENCES drafts (id),
    FOREIGN KEY (leader) REFERENCES leaders (id)
);

CREATE TABLE picks
(
    player TEXT NOT NULL,
    draft_id INTEGER NOT NULL,
    pick INTEGER NOT NULL,
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

CREATE TABLE ranks
(
    id INTEGER PRIMARY KEY,
    leader_id INTEGER NOT NULL,
    player_id INTEGER NOT NULL,
    tier FLOAT NOT NULL,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (leader_id) references leaders (id),
    FOREIGN KEY (player_id) references players (id)
);

CREATE TRIGGER update_ranks_timestamp
    AFTER UPDATE ON ranks
    FOR EACH ROW
    UPDATE ranks SET updated_at = CURRENT_TIMESTAMP WHERE id = old.id;

-- Add useful indexes
CREATE INDEX idx_pool_draft_id ON pool (draft_id);
CREATE INDEX idx_picks_draft_id ON picks (draft_id);
CREATE INDEX idx_draft_registry_draft_id ON draft_registry (draft_id);

-- +goose Down
DROP TRIGGER IF EXISTS update_ranks_timestamp;
DROP TABLE IF EXISTS ranks;
DROP TABLE IF EXISTS picks;
DROP TABLE IF EXISTS pool;
DROP TABLE IF EXISTS draft_registry;
DROP TABLE IF EXISTS players;
DROP TABLE IF EXISTS drafts;
DROP TABLE IF EXISTS leaders;
DROP INDEX IF EXISTS leaders_civ_name_leader_name_uindex;
DROP INDEX IF EXISTS idx_pool_draft_id;
DROP INDEX IF EXISTS idx_picks_draft_id;
DROP INDEX IF EXISTS idx_draft_registry_draft_id;
