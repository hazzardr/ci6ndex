-- +goose Up
ALTER TABLE leaders ADD COLUMN unranked BOOLEAN NOT NULL DEFAULT FALSE;

-- Mark the 4 new BBG Expanded leaders as unranked until tiers are established.
INSERT INTO leaders (civ_name, leader_name, friendly_name, discord_emoji_string, banned, tier, unranked)
VALUES
    ('TAÍNO', 'ANACAONA', 'Anacaona', NULL, false, 0, true),
    ('POLAND', 'STANISLAW II', 'Stanislaw II', NULL, false, 0, true),
    ('AUSTRIA', 'MARIA THERESA', 'Maria Theresa', NULL, false, 0, true),
    ('GOTHS', 'THEODORIC', 'Theodoric', NULL, false, 0, true);

INSERT INTO documents (leader_id, doc_name, link)
VALUES
    (86, 'BBG', 'https://civ6bbg.github.io/en_US/leaders_7.5.html#Ta%C3%ADno%20Anacaona'),
    (87, 'BBG', 'https://civ6bbg.github.io/en_US/leaders_7.5.html#Poland%20Stanislaw%20II'),
    (88, 'BBG', 'https://civ6bbg.github.io/en_US/leaders_7.5.html#Austria%20Maria%20Theresa'),
    (89, 'BBG', 'https://civ6bbg.github.io/en_US/leaders_7.5.html#Goths%20Theodoric');

-- +goose Down
DELETE FROM documents WHERE leader_id IN (86, 87, 88, 89) AND doc_name = 'BBG';
DELETE FROM leaders WHERE id IN (86, 87, 88, 89);
ALTER TABLE leaders DROP COLUMN unranked;
