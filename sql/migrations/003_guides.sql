-- +goose Up
CREATE TABLE documents
(    
    id INTEGER PRIMARY KEY, 
    leader_id INTEGER NOT NULL,
    doc_name TEXT NOT NULL,
    link TEXT NOT NULL,
    FOREIGN KEY (leader_id) REFERENCES leaders (id)
);

INSERT INTO documents (leader_id, doc_name, link)
VALUES
-- America
    (1, 'Zigzagzigal', 'https://steamcommunity.com/sharedfiles/filedetails/?id=2466666825'),
    (2, 'Zigzagzigal', 'https://steamcommunity.com/sharedfiles/filedetails/?id=2466666825'),
    (3, 'Zigzagzigal', 'https://steamcommunity.com/sharedfiles/filedetails/?id=2466666825'),
-- Arabia
    (4, 'Zigzagzigal', 'https://steamcommunity.com/sharedfiles/filedetails/?id=2444171808'),
    (5, 'Zigzagzigal', 'https://steamcommunity.com/sharedfiles/filedetails/?id=2444171808'),
-- Australia
    (6, 'Zigzagzigal', 'https://steamcommunity.com/sharedfiles/filedetails/?id=2542999669'),
-- Aztec
    (7, 'Zigzagzigal', 'https://steamcommunity.com/sharedfiles/filedetails/?id=1929535744'),
-- Babylon
    (8, 'Zigzagzigal', 'https://steamcommunity.com/sharedfiles/filedetails/?id=2387413001'),
-- Brazil
    (9, 'Zigzagzigal', 'https://steamcommunity.com/sharedfiles/filedetails/?id=1994150511'),
-- Byzantium
    (10, 'Zigzagzigal', 'https://steamcommunity.com/sharedfiles/filedetails/?id=2250797073'),
    (11, 'Zigzagzigal', 'https://steamcommunity.com/sharedfiles/filedetails/?id=2250797073'),
-- Canada
    (12, 'Zigzagzigal', 'https://steamcommunity.com/sharedfiles/filedetails/?id=1826638208'),
-- China
    (13, 'Zigzagzigal', 'https://steamcommunity.com/sharedfiles/filedetails/?id=1894800189'),
    (14, 'Zigzagzigal', 'https://steamcommunity.com/sharedfiles/filedetails/?id=1894800189'),
    (15, 'Zigzagzigal', 'https://steamcommunity.com/sharedfiles/filedetails/?id=1894800189'),
    (16, 'Zigzagzigal', 'https://steamcommunity.com/sharedfiles/filedetails/?id=1894800189'),
    (17, 'Zigzagzigal', 'https://steamcommunity.com/sharedfiles/filedetails/?id=1894800189'),
-- Cree
    (18, 'Zigzagzigal', 'https://steamcommunity.com/sharedfiles/filedetails/?id=1911704556'),
-- Dutch
    (19, 'Zigzagzigal', 'https://steamcommunity.com/sharedfiles/filedetails/?id=2474216503'),
-- Egypt
    (20, 'Zigzagzigal', 'https://steamcommunity.com/sharedfiles/filedetails/?id=2129421830'),
    (21, 'Zigzagzigal', 'https://steamcommunity.com/sharedfiles/filedetails/?id=2129421830'),
    (22, 'Zigzagzigal', 'https://steamcommunity.com/sharedfiles/filedetails/?id=2129421830'),
-- England
    (23, 'Zigzagzigal', 'https://steamcommunity.com/sharedfiles/filedetails/?id=1905994186'),
    (24, 'Zigzagzigal', 'https://steamcommunity.com/sharedfiles/filedetails/?id=1905994186'),
    (25, 'Zigzagzigal', 'https://steamcommunity.com/sharedfiles/filedetails/?id=1905994186'),
    (26, 'Zigzagzigal', 'https://steamcommunity.com/sharedfiles/filedetails/?id=1905994186'),
-- Ethiopia
    (27, 'Zigzagzigal', 'https://steamcommunity.com/sharedfiles/filedetails/?id=2540764190'),
-- France
    (28, 'Zigzagzigal', 'https://steamcommunity.com/sharedfiles/filedetails/?id=1787356472'),
    (29, 'Zigzagzigal', 'https://steamcommunity.com/sharedfiles/filedetails/?id=1787356472'),
    (30, 'Zigzagzigal', 'https://steamcommunity.com/sharedfiles/filedetails/?id=1787356472'),
-- Gaul
    (31, 'Zigzagzigal', 'https://steamcommunity.com/sharedfiles/filedetails/?id=2238038366'),
-- Georgia
    (32, 'Zigzagzigal', 'https://steamcommunity.com/sharedfiles/filedetails/?id=2095198631'),
-- Germany
    (33, 'Zigzagzigal', 'https://steamcommunity.com/sharedfiles/filedetails/?id=2536738340'),
    (34, 'Zigzagzigal', 'https://steamcommunity.com/sharedfiles/filedetails/?id=2536738340'),
-- Gran Colombia
    (35, 'Zigzagzigal', 'https://steamcommunity.com/sharedfiles/filedetails/?id=2106914837'),
-- Greece
    (36, 'Zigzagzigal', 'https://steamcommunity.com/sharedfiles/filedetails/?id=1980899859'),
    (37, 'Zigzagzigal', 'https://steamcommunity.com/sharedfiles/filedetails/?id=1980899859'),
-- Hungary
    (38, 'Zigzagzigal', 'https://steamcommunity.com/sharedfiles/filedetails/?id=1792612933'),
-- Inca
    (39, 'Zigzagzigal', 'https://steamcommunity.com/sharedfiles/filedetails/?id=1672524904'),
-- India
    (40, 'Zigzagzigal', 'https://steamcommunity.com/sharedfiles/filedetails/?id=1984819134'),
    (41, 'Zigzagzigal', 'https://steamcommunity.com/sharedfiles/filedetails/?id=1984819134'),
-- Indonesia
    (42, 'Zigzagzigal', 'https://steamcommunity.com/sharedfiles/filedetails/?id=2529715786'),
-- Japan
    (43, 'Zigzagzigal', 'https://steamcommunity.com/sharedfiles/filedetails/?id=2423465889'),
    (44, 'Zigzagzigal', 'https://steamcommunity.com/sharedfiles/filedetails/?id=2423465889'),
-- Khmer
    (45, 'Zigzagzigal', 'https://steamcommunity.com/sharedfiles/filedetails/?id=1911702737'),
-- Kongo
    (46, 'Zigzagzigal', 'https://steamcommunity.com/sharedfiles/filedetails/?id=2533397288'),
    (47, 'Zigzagzigal', 'https://steamcommunity.com/sharedfiles/filedetails/?id=2533397288'),
-- Korea
    (48, 'Zigzagzigal', 'https://steamcommunity.com/sharedfiles/filedetails/?id=1928634264'),
    (49, 'Zigzagzigal', 'https://steamcommunity.com/sharedfiles/filedetails/?id=1928634264'),
-- Macedon
    (50, 'Zigzagzigal', 'https://steamcommunity.com/sharedfiles/filedetails/?id=1902563180'),
-- Mali
    (51, 'Zigzagzigal', 'https://steamcommunity.com/sharedfiles/filedetails/?id=1918595516'),
    (52, 'Zigzagzigal', 'https://steamcommunity.com/sharedfiles/filedetails/?id=1918595516'),
-- Maori
    (53, 'Zigzagzigal', 'https://steamcommunity.com/sharedfiles/filedetails/?id=1800180355'),
-- Mapuche
    (54, 'Zigzagzigal', 'https://steamcommunity.com/sharedfiles/filedetails/?id=2477806398'),
-- Maya
    (55, 'Zigzagzigal', 'https://steamcommunity.com/sharedfiles/filedetails/?id=2104285340'),
-- Mongolia
    (56, 'Zigzagzigal', 'https://steamcommunity.com/sharedfiles/filedetails/?id=2446715153'),
    (57, 'Zigzagzigal', 'https://steamcommunity.com/sharedfiles/filedetails/?id=2446715153'),
-- Norway
    (58, 'Zigzagzigal', 'https://steamcommunity.com/sharedfiles/filedetails/?id=1935960977'),
    (59, 'Zigzagzigal', 'https://steamcommunity.com/sharedfiles/filedetails/?id=1935960977'),
-- Nubia
    (60, 'Zigzagzigal', 'https://steamcommunity.com/sharedfiles/filedetails/?id=1905728008'),
-- Ottoman
    (61, 'Zigzagzigal', 'https://steamcommunity.com/sharedfiles/filedetails/?id=1840773267'),
    (62, 'Zigzagzigal', 'https://steamcommunity.com/sharedfiles/filedetails/?id=1840773267'),
-- Persia
    (63, 'Zigzagzigal', 'https://steamcommunity.com/sharedfiles/filedetails/?id=2096074206'),
    (64, 'Zigzagzigal', 'https://steamcommunity.com/sharedfiles/filedetails/?id=2096074206'),
-- Phoencia
    (65, 'Zigzagzigal', 'https://steamcommunity.com/sharedfiles/filedetails/?id=1896495196'),
-- Poland
    (66, 'Zigzagzigal', 'https://steamcommunity.com/sharedfiles/filedetails/?id=1931279541'),
-- Portugal
    (67, 'Zigzagzigal', 'https://steamcommunity.com/sharedfiles/filedetails/?id=2436226637'),
-- Rome
    (68, 'Zigzagzigal', 'https://steamcommunity.com/sharedfiles/filedetails/?id=2244514508'),
    (69, 'Zigzagzigal', 'https://steamcommunity.com/sharedfiles/filedetails/?id=2244514508'),
-- Russia
    (70, 'Zigzagzigal', 'https://steamcommunity.com/sharedfiles/filedetails/?id=1923759261'),
-- Scotland
    (71, 'Zigzagzigal', 'https://steamcommunity.com/sharedfiles/filedetails/?id=2527343587'),
-- Scythia
    (72, 'Zigzagzigal', 'https://steamcommunity.com/sharedfiles/filedetails/?id=2094299592'),
-- Spain
    (73, 'Zigzagzigal', 'https://steamcommunity.com/sharedfiles/filedetails/?id=2527840104'),
-- Sumeria
    (74, 'Zigzagzigal', 'https://steamcommunity.com/sharedfiles/filedetails/?id=1905022232'),
-- Sweden
    (75, 'Zigzagzigal', 'https://steamcommunity.com/sharedfiles/filedetails/?id=1891321127'),
-- Vietnam
    (76, 'Zigzagzigal', 'https://steamcommunity.com/sharedfiles/filedetails/?id=2589592137'),
-- Zulu
    (77, 'Zigzagzigal', 'https://steamcommunity.com/sharedfiles/filedetails/?id=1983949096')
;

-- +goose Down
DROP TABLE IF EXISTS documents;
