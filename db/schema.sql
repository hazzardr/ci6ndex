CREATE SCHEMA IF NOT EXISTS ci6ndex;

CREATE TABLE ci6ndex.users (
    id BIGSERIAL NOT NULL,
    discord_name TEXT NOT NULL,
    name TEXT NOT NULL,
    PRIMARY KEY (id)
);

CREATE UNIQUE INDEX users_discord_name_uindex ON ci6ndex.users (discord_name);

CREATE TABLE ci6ndex.leaders
(
    id BIGSERIAL NOT NULL,
    civ_name TEXT NOT NULL,
    leader_name TEXT NOT NULL,
    PRIMARY KEY (id)
);

CREATE UNIQUE INDEX leaders_civ_name_leader_name_uindex ON ci6ndex.leaders (civ_name, leader_name);

CREATE TABLE ci6ndex.rankings
(
    id BIGSERIAL NOT NULL,
    user_id BIGINT NOT NULL,
    tier FLOAT NOT NULL,
    leader_id BIGINT NOT NULL,
    PRIMARY KEY (id),
    CONSTRAINT users_fk FOREIGN KEY (user_id) REFERENCES ci6ndex.users (id),
    CONSTRAINT leaders_fk FOREIGN KEY (leader_id) REFERENCES ci6ndex.leaders (id)
);

CREATE TABLE ci6ndex.draft_strategies
(
  name TEXT UNIQUE NOT NULL,
  description TEXT NOT NULL,
  rules JSONB,
  PRIMARY KEY (name)
);

COMMENT ON TABLE ci6ndex.draft_strategies IS 'The strategies that can be used to draft a civ';
COMMENT ON COLUMN ci6ndex.draft_strategies.rules IS 'Specific rules that this draft has to follow.';

CREATE TABLE ci6ndex.drafts
(
    id BIGSERIAL NOT NULL,
    draft_strategy TEXT NOT NULL,
    PRIMARY KEY (id),
    CONSTRAINT draft_strategies_fk FOREIGN KEY (draft_strategy) REFERENCES ci6ndex.draft_strategies (name)
);

CREATE TABLE ci6ndex.draft_picks
(
    id BIGSERIAL NOT NULL,
    draft_id BIGINT NOT NULL,
    user_id BIGINT NOT NULL,
    leader_id BIGINT NOT NULL,
    offered JSONB,
    PRIMARY KEY (id),
    CONSTRAINT drafts_fk FOREIGN KEY (draft_id) REFERENCES ci6ndex.drafts (id),
    CONSTRAINT users_fk FOREIGN KEY (user_id) REFERENCES ci6ndex.users (id),
    CONSTRAINT leaders_fk FOREIGN KEY (leader_id) REFERENCES ci6ndex.leaders (id)
);

COMMENT ON COLUMN ci6ndex.draft_picks.offered IS 'The civs that were offered to the user.';
CREATE UNIQUE INDEX draft_picks_draft_id_user_id_uindex ON ci6ndex.draft_picks (draft_id, user_id);

CREATE TABLE ci6ndex.games
(
    id BIGSERIAL NOT NULL,
    draft_id BIGINT NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE,
    game_stats JSONB,
    PRIMARY KEY (id),
    CONSTRAINT drafts_fk FOREIGN KEY (draft_id) REFERENCES ci6ndex.drafts (id)
);

CREATE TABLE ci6ndex.stats
(
    id BIGSERIAL NOT NULL,
    stats JSONB NOT NULL,
    user_id BIGINT NOT NULL,
    game_id BIGINT NOT NULL,
    CONSTRAINT users_fk FOREIGN KEY (user_id) REFERENCES ci6ndex.users (id),
    CONSTRAINT games_fk FOREIGN KEY (game_id) REFERENCES ci6ndex.games (id),
    PRIMARY KEY (id)
);
