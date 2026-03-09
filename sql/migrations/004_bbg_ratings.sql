-- +goose Up
UPDATE leaders SET banned = FALSE;

INSERT INTO leaders (civ_name, leader_name, discord_emoji_string, banned, tier)
VALUES
    ('GAUL', 'VERCINGETORIX','<:Ambiorix_Civ6:1229388711087702046>',false,3),
    ('MACEDON', 'OLYMPIAS','<:Alexander_Civ6:1229381348251406417>',false,3),
    ('MAYA', 'TE KINICH II','<:Lady_Six_Sky_Civ6:1229393212645572618>',false,3),
    ('PHOENICIAN', 'AHIRAM','<:Dido_Civ6:1229388795565309984>',false,3),
    ('SWAHILI', 'AL-HASAN IBN SULAIMAN','<:Suleiman_Civ6:1229599564277874720>',false,2.75),
    ('TEOTIHUACAN', 'SPEARTHROWER OWL','<:Montezuma_Civ6:1229393522717622313>', false, 1.92),
    ('THULE', 'KIVIUQ','<:Gandhi_Civ6:1229388944014049300>',false, 6),
    ('TIBET', 'TRISONG DETSEN','<:Gandhi_Civ6:1229388944014049300>',false, 6);


ALTER TABLE ranks ADD COLUMN bbg BOOLEAN NOT NULL DEFAULT FALSE;

DROP INDEX ranks_leader_id_player_id_uindex;

CREATE UNIQUE INDEX ranks_leader_id_player_id_bbg_uindex ON ranks(leader_id, player_id, bbg);

INSERT INTO documents (leader_id, doc_name, link)
VALUES
    (1, 'BBG', 'https://civ6bbg.github.io/en_US/leaders_7.2.html#AMERICA%20ABE'),
    (2, 'BBG', 'https://civ6bbg.github.io/en_US/leaders_7.2.html#AMERICA%20BULLMOOSE%20TEDDY'),
    (3, 'BBG', 'https://civ6bbg.github.io/en_US/leaders_7.2.html#AMERICA%20ROUGH%20RIDER%20TEDDY'),
    (4, 'BBG', 'https://civ6bbg.github.io/en_US/leaders_7.2.html#ARABIA%20SALADIN%20SULTAN'),
    (5, 'BBG', 'https://civ6bbg.github.io/en_US/leaders_7.2.html#ARABIA%20SALADIN%20VIZIR'),
    (6, 'BBG', 'https://civ6bbg.github.io/en_US/leaders_7.2.html#AUSTRALIA%20JOHN%20CURTIN'),
    (7, 'BBG', 'https://civ6bbg.github.io/en_US/leaders_7.2.html#AZTEC%20MONTEZUMA'),
    (8, 'BBG', 'https://civ6bbg.github.io/en_US/leaders_7.2.html#BABYLON%20HAMMURABI'),
    (9, 'BBG', 'https://civ6bbg.github.io/en_US/leaders_7.2.html#BRAZIL%20PEDRO%20II'),
    (10, 'BBG', 'https://civ6bbg.github.io/en_US/leaders_7.2.html#BYZANTIUM%20BASIL%20II'),
    (11, 'BBG', 'https://civ6bbg.github.io/en_US/leaders_7.2.html#BYZANTIUM%20THEODORA'),
    (12, 'BBG', 'https://civ6bbg.github.io/en_US/leaders_7.2.html#CANADA%20WILFRID%20LAURIER'),
    (13, 'BBG', 'https://civ6bbg.github.io/en_US/leaders_7.2.html#CHINA%20KUBLAI%20KHAN'),
    (14, 'BBG', 'https://civ6bbg.github.io/en_US/leaders_7.2.html#CHINA%20QIN%20SHI%20HUANG%20UNIFIER'),
    (15, 'BBG', 'https://civ6bbg.github.io/en_US/leaders_7.2.html#CHINA%20QIN%20SHI%20HUANG'),
    (16, 'BBG', 'https://civ6bbg.github.io/en_US/leaders_7.2.html#CHINA%20WU%20ZEITAN'),
    (17, 'BBG', 'https://civ6bbg.github.io/en_US/leaders_7.2.html#CHINA%20YONGLE'),
    (18, 'BBG', 'https://civ6bbg.github.io/en_US/leaders_7.2.html#CREE%20POUNDMAKER'),
    (19, 'BBG', 'https://civ6bbg.github.io/en_US/leaders_7.2.html#DUTCH%20WILHELMINA'),
    (20, 'BBG', 'https://civ6bbg.github.io/en_US/leaders_7.2.html#EGYPT%20CLEOPATRA'),
    (21, 'BBG', 'https://civ6bbg.github.io/en_US/leaders_7.2.html#EGYPT%20PTOLEMEIC%20CLEO'),
    (22, 'BBG', 'https://civ6bbg.github.io/en_US/leaders_7.2.html#EGYPT%20RAMSEYS'),
    (23, 'BBG', 'https://civ6bbg.github.io/en_US/leaders_7.2.html#ENGLAND%20ELEANOR%20OF%20AQUITAINE'),
    (24, 'BBG', 'https://civ6bbg.github.io/en_US/leaders_7.2.html#ENGLAND%20ELIZABETH'),
    (25, 'BBG', 'https://civ6bbg.github.io/en_US/leaders_7.2.html#ENGLAND%20STEAMY%20VICKY'),
    (26, 'BBG', 'https://civ6bbg.github.io/en_US/leaders_7.2.html#ENGLAND%20VICTORIA'),
    (27, 'BBG', 'https://civ6bbg.github.io/en_US/leaders_7.2.html#ETHIOPIA%20MENELIK%20II'),
    (28, 'BBG', 'https://civ6bbg.github.io/en_US/leaders_7.2.html#FRANCE%20CATHERINE%20DE%20MEDICI'),
    (29, 'BBG', 'https://civ6bbg.github.io/en_US/leaders_7.2.html#FRANCE%20ELEANOR%20AQUITAINE'),
    (30, 'BBG', 'https://civ6bbg.github.io/en_US/leaders_7.2.html#FRANCE%20MAGNIFICENCE%20CATHERINE'),
    (31, 'BBG', 'https://civ6bbg.github.io/en_US/leaders_7.2.html#GAUL%20AMBIORIX'),
    (32, 'BBG', 'https://civ6bbg.github.io/en_US/leaders_7.2.html#GEORGIA%20TAMAR'),
    (33, 'BBG', 'https://civ6bbg.github.io/en_US/leaders_7.2.html#GERMANY%20FREDERICK%20BARBAROSSA'),
    (34, 'BBG', 'https://civ6bbg.github.io/en_US/leaders_7.2.html#GERMANY%20LUDWIG'),
    (35, 'BBG', 'https://civ6bbg.github.io/en_US/leaders_7.2.html#GRAN%20COLUMBIA%20SIMON%20BOLIVAR'),
    (36, 'BBG', 'https://civ6bbg.github.io/en_US/leaders_7.2.html#GREECE%20GORGO'),
    (37, 'BBG', 'https://civ6bbg.github.io/en_US/leaders_7.2.html#GREECE%20PERICLES'),
    (38, 'BBG', 'https://civ6bbg.github.io/en_US/leaders_7.2.html#HUNGARY%20MATTHIAS%20CORVINUS'),
    (39, 'BBG', 'https://civ6bbg.github.io/en_US/leaders_7.2.html#INCA%20PACHACUTI'),
    (40, 'BBG', 'https://civ6bbg.github.io/en_US/leaders_7.2.html#INDIA%20CHANDRAGUPTA'),
    (41, 'BBG', 'https://civ6bbg.github.io/en_US/leaders_7.2.html#INDIA%20GANDHI'),
    (42, 'BBG', 'https://civ6bbg.github.io/en_US/leaders_7.2.html#INDONESIA%20GITARJA'),
    (43, 'BBG', 'https://civ6bbg.github.io/en_US/leaders_7.2.html#JAPAN%20HOJO%20TOKIMUNE'),
    (44, 'BBG', 'https://civ6bbg.github.io/en_US/leaders_7.2.html#JAPAN%20TOKUGAWA'),
    (45, 'BBG', 'https://civ6bbg.github.io/en_US/leaders_7.2.html#KHMER%20JAYAVARMAN%20VII'),
    (46, 'BBG', 'https://civ6bbg.github.io/en_US/leaders_7.2.html#KONGO%20MVEMBA%20A%20NZINGA'),
    (47, 'BBG', 'https://civ6bbg.github.io/en_US/leaders_7.2.html#KONGO%20NZINGA%20MBANDE'),
    (48, 'BBG', 'https://civ6bbg.github.io/en_US/leaders_7.2.html#KOREA%20SEJONG'),
    (49, 'BBG', 'https://civ6bbg.github.io/en_US/leaders_7.2.html#KOREA%20SEONDEOK'),
    (50, 'BBG', 'https://civ6bbg.github.io/en_US/leaders_7.2.html#MACEDON%20ALEXANDER'),
    (51, 'BBG', 'https://civ6bbg.github.io/en_US/leaders_7.2.html#MALI%20MANSA%20MUSA'),
    (52, 'BBG', 'https://civ6bbg.github.io/en_US/leaders_7.2.html#MALI%20SUNDIATA%20KEITA'),
    (53, 'BBG', 'https://civ6bbg.github.io/en_US/leaders_7.2.html#MAORI%20KUPE'),
    (54, 'BBG', 'https://civ6bbg.github.io/en_US/leaders_7.2.html#MAPUCHE%20LAUTARO'),
    (55, 'BBG', 'https://civ6bbg.github.io/en_US/leaders_7.2.html#MAYA%20LADY%20SIX%20SKY'),
    (56, 'BBG', 'https://civ6bbg.github.io/en_US/leaders_7.2.html#MONGOLIA%20GENGHIS%20KHAN'),
    (57, 'BBG', 'https://civ6bbg.github.io/en_US/leaders_7.2.html#MONGOLIA%20KUBLAI%20KHAN'),
    (58, 'BBG', 'https://civ6bbg.github.io/en_US/leaders_7.2.html#NORWAY%20HARALD%20HARDRADA'),
    (59, 'BBG', 'https://civ6bbg.github.io/en_US/leaders_7.2.html#NORWAY%20VARANGIA%20HARALD'),
    (60, 'BBG', 'https://civ6bbg.github.io/en_US/leaders_7.2.html#NUBIA%20AMANITORE'),
    (61, 'BBG', 'https://civ6bbg.github.io/en_US/leaders_7.2.html#OTTOMAN%20SULEIMAN%20MUHTEŞEM'),
    (62, 'BBG', 'https://civ6bbg.github.io/en_US/leaders_7.2.html#OTTOMAN%20SULEIMAN'),
    (63, 'BBG', 'https://civ6bbg.github.io/en_US/leaders_7.2.html#PERSIA%20CYRUS'),
    (64, 'BBG', 'https://civ6bbg.github.io/en_US/leaders_7.2.html#PERSIA%20NADER%20SHAH'),
    (65, 'BBG', 'https://civ6bbg.github.io/en_US/leaders_7.2.html#PHOENICIAN%20DIDO'),
    (66, 'BBG', 'https://civ6bbg.github.io/en_US/leaders_7.2.html#POLISH%20JADWIGA'),
    (67, 'BBG', 'https://civ6bbg.github.io/en_US/leaders_7.2.html#PORTUGAL%20JOÃO%20III'),
    (68, 'BBG', 'https://civ6bbg.github.io/en_US/leaders_7.2.html#ROME%20JULIUS%20CAESER'),
    (69, 'BBG', 'https://civ6bbg.github.io/en_US/leaders_7.2.html#ROME%20TRAJAN'),
    (70, 'BBG', 'https://civ6bbg.github.io/en_US/leaders_7.2.html#RUSSIA%20PETER'),
    (71, 'BBG', 'https://civ6bbg.github.io/en_US/leaders_7.2.html#SCOTLAND%20ROBERT%20THE%20BRUCE'),
    (72, 'BBG', 'https://civ6bbg.github.io/en_US/leaders_7.2.html#SCYTHIA%20TOMYRIS'),
    (73, 'BBG', 'https://civ6bbg.github.io/en_US/leaders_7.2.html#SPAIN%20PHILLIP%20II'),
    (74, 'BBG', 'https://civ6bbg.github.io/en_US/leaders_7.2.html#SUMERIA%20GILGAMESH'),
    (75, 'BBG', 'https://civ6bbg.github.io/en_US/leaders_7.2.html#SWEDEN%20KRISTINA'),
    (76, 'BBG', 'https://civ6bbg.github.io/en_US/leaders_7.2.html#VIETNAM%20BÀ%20TRIỆU'),
    (77, 'BBG', 'https://civ6bbg.github.io/en_US/leaders_7.2.html#ZULU%20SHAKA');

-- +goose Down
DELETE FROM documents WHERE doc_name = 'BBG';

DROP INDEX ranks_leader_id_player_id_bbg_uindex;

CREATE UNIQUE INDEX ranks_leader_id_player_id_uindex ON ranks(leader_id, player_id);

ALTER TABLE ranks DROP COLUMN bbg;
