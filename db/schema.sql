CREATE SCHEMA IF NOT EXISTS ci6ndex;

CREATE TABLE ci6ndex.users (
    id BIGSERIAL NOT NULL,
    discord_name TEXT NOT NULL,
    name TEXT NOT NULL,
    PRIMARY KEY (id)
);

CREATE UNIQUE INDEX users_discord_name_uindex ON ci6ndex.users (discord_name);

CREATE TABLE ci6ndex.games (
   id BIGSERIAL NOT NULL,
   user_ids int[] NOT NULL,
   PRIMARY KEY (id)
);

CREATE TABLE ci6ndex.stats (
    id BIGSERIAL NOT NULL,
    stats JSONB NOT NULL,
    user_id int NOT NULL,
    game_id int NOT NULL,
    CONSTRAINT users_fk FOREIGN KEY (user_id) REFERENCES ci6ndex.users (id),
    CONSTRAINT games_fk FOREIGN KEY (game_id) REFERENCES ci6ndex.games (id),
    PRIMARY KEY (id)
);

CREATE TABLE ci6ndex.leaders (
    id BIGSERIAL NOT NULL,
    civ_name TEXT NOT NULL,
    leader_name TEXT NOT NULL,
    PRIMARY KEY (id)
);


CREATE TABLE ci6ndex.rankings (
  id BIGSERIAL NOT NULL,
  user_id INT NOT NULL,
  tier FLOAT NOT NULL,
  leader_id int NOT NULL,
  PRIMARY KEY (id),
  CONSTRAINT users_fk FOREIGN KEY (user_id) REFERENCES ci6ndex.users (id),
  CONSTRAINT leaders_fk FOREIGN KEY (leader_id) REFERENCES ci6ndex.leaders (id)
);
