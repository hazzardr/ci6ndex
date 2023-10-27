-- name: GetUser :one
SELECT * FROM ci6ndex.users WHERE id = $1 LIMIT 1;

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
         WHERE id = $1
         LIMIT 1;

-- name: GetLeaderByNameAndCiv :one
SELECT * FROM ci6ndex.leaders
WHERE leader_name = $1
AND civ_name = $2
LIMIT 1;

-- name: CreateDraftStrategy :one
INSERT INTO ci6ndex.draft_strategies
(
    name, description
) VALUES (
    $1, $2
)
RETURNING *;

-- name: GetDraftStrategy :one
SELECT * FROM ci6ndex.draft_strategies
WHERE name = $1
LIMIT 1;

-- name: GetDraftStrategies :many
SELECT * FROM ci6ndex.draft_strategies;

-- name: CreateDraft :one
INSERT INTO ci6ndex.drafts
(
    draft_strategy
) VALUES (
    $1
)
RETURNING *;

-- name: GetDraft :one
SELECT * FROM ci6ndex.drafts
WHERE id = $1
LIMIT 1;

-- name: SubmitDraftPick :one
INSERT INTO ci6ndex.draft_picks
(
    draft_id, leader_id, user_id, offered
) VALUES (
    $1, $2, $3, $4
)
RETURNING *;