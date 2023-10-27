-- name: GetUserByName :one
SELECT * FROM ci6ndex.users WHERE name = $1 LIMIT 1;

-- name: GetUserByDiscordName :one
SELECT * FROM ci6ndex.users WHERE discord_name = $1 LIMIT 1;

-- name: CreateUser :one
INSERT INTO ci6ndex.users
(
    name, discord_name
) VALUES (
    $1, $2
)
RETURNING *;

-- name: DeleteRankings :exec
DELETE FROM ci6ndex.rankings
RETURNING *;

-- name: CreateRanking :one
INSERT INTO ci6ndex.rankings
(
    user_id, tier, leader_id
) VALUES (
    $1, $2, $3
)
RETURNING *;

-- name: CreateRankings :copyfrom
INSERT INTO ci6ndex.rankings
(
    user_id, tier, leader_id
) VALUES (
    $1, $2, $3
);

-- name: GetLeader :one
SELECT * FROM ci6ndex.leaders
WHERE leader_name = $1
AND civ_name = $2
LIMIT 1;