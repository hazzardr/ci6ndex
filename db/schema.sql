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

CREATE TABLE ci6ndex.rankings
(
    id BIGSERIAL NOT NULL,
    user_id INT NOT NULL,
    tier FLOAT NOT NULL,
    leader_id int NOT NULL,
    PRIMARY KEY (id),
    CONSTRAINT users_fk FOREIGN KEY (user_id) REFERENCES ci6ndex.users (id),
    CONSTRAINT leaders_fk FOREIGN KEY (leader_id) REFERENCES ci6ndex.leaders (id)
);

CREATE TABLE ci6ndex.draft_strategies
(
  name TEXT NOT NULL,
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
    picks JSONB[],
    PRIMARY KEY (id),
    CONSTRAINT draft_strategies_fk FOREIGN KEY (draft_strategy) REFERENCES ci6ndex.draft_strategies (name)
);
COMMENT ON COLUMN ci6ndex.drafts.picks IS 'The picks in the draft. Should include at least leader and user';

CREATE TABLE ci6ndex.games
(
    id BIGSERIAL NOT NULL,
    draft_id int ,
    start_date DATE NOT NULL,
    end_date DATE,
    PRIMARY KEY (id),
    CONSTRAINT drafts_fk FOREIGN KEY (draft_id) REFERENCES ci6ndex.drafts (id)
);

CREATE TABLE ci6ndex.stats
(
    id BIGSERIAL NOT NULL,
    stats JSONB NOT NULL,
    user_id int NOT NULL,
    game_id int NOT NULL,
    CONSTRAINT users_fk FOREIGN KEY (user_id) REFERENCES ci6ndex.users (id),
    CONSTRAINT games_fk FOREIGN KEY (game_id) REFERENCES ci6ndex.games (id),
    PRIMARY KEY (id)
);
