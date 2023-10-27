DELETE FROM ci6ndex.draft_picks cascade;
DELETE FROM ci6ndex.rankings cascade;
DELETE FROM ci6ndex.leaders cascade;

INSERT INTO ci6ndex.leaders (civ_name, leader_name)
VALUES
    ('AMERICA', 'ABE'),
    ('AMERICA', 'BULLMOOSE TEDDY'),
    ('AMERICA', 'ROUGH RIDER TEDDY'),
    ('AMERICA', 'TEDDY'),
    ('ARABIA', 'SALADIN SULTAN'),
    ('ARABIA', 'SALADIN VIZIR'),
    ('AUSTRALIA', 'JOHN CURTIN'),
    ('AZTEC', 'MONTEZUMA'),
    ('BABYLON', 'HAMMURABI'),
    ('BRAZIL', 'PEDRO II'),
    ('BYZANTIUM', 'BASIL II'),
    ('BYZANTIUM', 'THEODORA'),
    ('CANADA', 'WILFRID LAURIER'),
    ('CHINA', 'KUBLAI KHAN'),
    ('CHINA', 'QIN SHI HUANG UNIFIER'),
    ('CHINA', 'QIN SHI HUANG'),
    ('CHINA', 'WU ZEITAN'),
    ('CHINA', 'YONGLE'),
    ('CREE', 'POUNDMAKER'),
    ('DUTCH', 'WILHELMINA'),
    ('EGYPT', 'CLEOPATRA'),
    ('EGYPT', 'PTOLEMEIC CLEO'),
    ('EGYPT', 'RAMSEYS'),
    ('ENGLAND', 'ELEANOR OF AQUITAINE'),
    ('ENGLAND', 'ELIZABETH'),
    ('ENGLAND', 'STEAMY VICKY'),
    ('ENGLAND', 'VICTORIA'),
    ('ETHIOPIA', 'MENELIK II'),
    ('FRANCE', 'CATHERINE DE MEDICI'),
    ('FRANCE', 'ELEANOR AQUITAINE'),
    ('FRANCE', 'MAGNIFICENCE CATHERINE'),
    ('GAUL', 'AMBIORIX'),
    ('GEORGIA', 'TAMAR'),
    ('GERMANY', 'FREDERICK BARBAROSSA'),
    ('GERMANY', 'LUDWIG'),
    ('GRAN COLUMBIA', 'SIMON BOLIVAR'),
    ('GREECE', 'GORGO'),
    ('GREECE', 'PERICLES'),
    ('HUNGARY', 'MATTHIAS CORVINUS'),
    ('INCA', 'PACHACUTI'),
    ('INDIA', 'CHANDRAGUPTA'),
    ('INDIA', 'GHANDI'),
    ('INDONESIA', 'GITARJA'),
    ('JAPAN', 'HOJO TOKIMUNE'),
    ('JAPAN', 'TOKUGAWA'),
    ('KHMER', 'JAYAVARMAN VII'),
    ('KONGO', 'MVEMBA A NZINGA'),
    ('KONGO', 'NZINGA MBANDE'),
    ('KOREA', 'SEJONG'),
    ('KOREA', 'SEONDEOK'),
    ('MACEDON', 'ALEXANDER'),
    ('MALI', 'MANSA MUSA'),
    ('MALI', 'SUNDIATA KEITA'),
    ('MAORI', 'KUPE'),
    ('MAPUCHE', 'LAUTARO'),
    ('MAYA', 'LADY SIX SKY'),
    ('MONGOLIA', 'GENGHIS KHAN'),
    ('MONGOLIA', 'KUBLAI KHAN'),
    ('NORWAY', 'HARALD HARDRADA'),
    ('NORWAY', 'VARANGIA HARALD'),
    ('NUBIA', 'AMANITORE'),
    ('OTTOMAN', 'SULEIMAN MUHTEŞEM'),
    ('OTTOMAN', 'SULEIMAN'),
    ('PERSIA', 'CYRUS'),
    ('PERSIA', 'NADER SHAH'),
    ('PHOENICIAN', 'DIDO'),
    ('POLISH', 'JADWIGA'),
    ('PORTUGAL', 'JOÃO III'),
    ('ROME', 'JULIUS CAESER'),
    ('ROME', 'TRAJAN'),
    ('RUSSIA', 'PETER'),
    ('SCOTLAND', 'ROBERT THE BRUCE'),
    ('SCYTHIA', 'TOMYRIS'),
    ('SPAIN', 'PHILLIP II'),
    ('SUMERIA', 'GILGAMESH'),
    ('SWEDEN', 'KRISTINA'),
    ('VIETNAM', 'BÀ TRIỆU'),
    ('ZULU', 'SHAKA')
;

INSERT INTO ci6ndex.draft_strategies (name, description, rules)
VALUES
(
 'AllPick',
 'Everyone can freely pick a leader, no restrictions.',
 null
),
(
 'RandomPick',
 'Everyone gets a randomized leader.',
 '{ "randomize": true }'
),
(
 'RandomPickPool3',
 'Everyone gets assigned a pool of 3 random leaders to pick from.',
 '{ "randomize": true, "pool_size": 3 }'
),
--     Potentially in the future add attr. like monger, culture etc?
(
 'RandomPickPool3Standard',
 'Everyone gets assigned a pool of 3 random leaders to pick from. Includes the discord \"Standard\" rules.',
 '{ "randomize": true, "pool_size": 3, "tiers_offered_min": {"1": 1}}'
)
;