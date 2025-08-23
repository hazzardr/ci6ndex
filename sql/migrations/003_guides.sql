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
    (1, "Zigzagzigal", 'https://steamcommunity.com/sharedfiles/filedetails/?id=2466666825')
    (2, "Zigzagzigal", 'https://steamcommunity.com/sharedfiles/filedetails/?id=2466666825')
    (3, "Zigzagzigal", 'https://steamcommunity.com/sharedfiles/filedetails/?id=2466666825')
    (4, "Zigzagzigal", 'https://steamcommunity.com/sharedfiles/filedetails/?id=2444171808')
    (5, "Zigzagzigal", 'https://steamcommunity.com/sharedfiles/filedetails/?id=2444171808')
    (6, "Zigzagzigal", 'https://steamcommunity.com/sharedfiles/filedetails/?id=2542999669')
    (7, "Zigzagzigal", 'https://steamcommunity.com/sharedfiles/filedetails/?id=1929535744')
    (8, "Zigzagzigal", 'https://steamcommunity.com/sharedfiles/filedetails/?id=2387413001')
;
-- +goose Down
DROP TABLE IF EXISTS documents;
