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
  tier INT NOT NULL,
  leader_id int NOT NULL,
  PRIMARY KEY (id),
  CONSTRAINT users_fk FOREIGN KEY (user_id) REFERENCES ci6ndex.users (id),
  CONSTRAINT leaders_fk FOREIGN KEY (leader_id) REFERENCES ci6ndex.leaders (id)
);


INSERT INTO ci6ndex.leaders (civ_name, leader_name)
VALUES
    ('America', 'Abe'),
    ('America', 'Abe'),
    ('America', 'Bullmoose Teddy'),
    ('America', 'Rough Rider Teddy'),
    ('America', 'Teddy'),
    ('Arabia', 'Saladin Sultan'),
    ('Arabia', 'Saladin Vizir'),
    ('Australia', 'John Curtin'),
    ('Aztec', 'Montezuma'),
    ('Babylon', 'Hammurabi'),
    ('Brazil', 'Pedro II'),
    ('Byzantium', 'Basil II'),
    ('Byzantium', 'Theodora'),
    ('Canada', 'Wilfrid Laurier'),
    ('China', 'Kublai Khan'),
    ('China', 'Qin Shi Huang Unifier'),
    ('China', 'Qin Shi Huang'),
    ('China', 'Wu Zeitan'),
    ('China', 'Yongle'),
    ('Cree', 'Poundmaker'),
    ('Dutch', 'Wilhelmina'),
    ('Egypt', 'Cleopatra'),
    ('Egypt', 'Ptolemeic Cleo'),
    ('Egypt', 'Ramseys'),
    ('England', 'Eleanor of Aquitaine'),
    ('England', 'Elizabeth'),
    ('England', 'Steamy Vicky'),
    ('England', 'Victoria'),
    ('Ethiopia', 'Menelik II'),
    ('France', 'Catherine de Medici'),
    ('France', 'Eleanor Aquitaine'),
    ('France', 'Magnificence Catherine'),
    ('Gaul', 'Ambiorix'),
    ('Georgia', 'Tamar'),
    ('Germany', 'Frederick Barbarossa'),
    ('Germany', 'Ludwig'),
    ('Gran Columbia', 'Simon Bolivar'),
    ('Greece', 'Gorgo'),
    ('Greece', 'Pericles'),
    ('Hungary', 'Matthias Corvinus'),
    ('Inca', 'Pachacuti'),
    ('India', 'Chandragupta'),
    ('India', 'Ghandi'),
    ('Indonesia', 'Gitarja'),
    ('Japan', 'Hojo Tokimune'),
    ('Japan', 'Tokugawa'),
    ('Khmer', 'Jayavarman VII'),
    ('Kongo', 'Mvemba a Nzinga'),
    ('Kongo', 'Nzinga Mbande'),
    ('Korea', 'Sejong'),
    ('Korea', 'Seondeok'),
    ('Macedon', 'Alexander'),
    ('Mali', 'Mansa Musa'),
    ('Mali', 'Sundiata Keita'),
    ('Maori', 'Kupe'),
    ('Mapuche', 'Lautara'),
    ('Maya', 'Lady Six Sky'),
    ('Mongolia', 'Genghis Khan'),
    ('Mongolia', 'Kublai Khan'),
    ('Norway', 'Harald Hardrada'),
    ('Norway', 'Varangia Harald'),
    ('Nubia', 'Amanitore'),
    ('Ottoman', 'Suleiman Muhteşem'),
    ('Ottoman', 'Suleiman'),
    ('Persia', 'Cyrus'),
    ('Persia', 'Nader Shah'),
    ('Phoenician', 'Dido'),
    ('Polish', 'Jadwiga'),
    ('Portugal', 'João III'),
    ('Rome', 'Julius Caeser'),
    ('Rome', 'Trajan'),
    ('Russia', 'Peter'),
    ('Scotland', 'Robert the Bruce'),
    ('Scythia', 'Tomyris'),
    ('Spain', 'Phillip II'),
    ('Sumeria', 'Gilgamesh'),
    ('Sweden', 'Kristina'),
    ('Vietnam', 'Bà Triệu'),
    ('Zulu', 'Shaka')
;
