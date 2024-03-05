-- name: GetUser :one
SELECT * FROM ci6ndex.users WHERE id = $1 LIMIT 1;

-- name: GetUserByName :one
SELECT * FROM ci6ndex.users WHERE name = $1 LIMIT 1;

-- name: GetUserByDiscordName :one
SELECT * FROM ci6ndex.users WHERE discord_name = $1 LIMIT 1;

-- name: GetUsers :many
SELECT * FROM ci6ndex.users;

-- name: CreateUser :one
INSERT INTO ci6ndex.users
(
    name, discord_name
) VALUES (
    $1, $2
)
RETURNING *;

-- name: CreateUsers :copyfrom
INSERT INTO ci6ndex.users
(
    name, discord_name
) VALUES (
    $1, $2
);

-- name: GetRankings :many
SELECT * FROM ci6ndex.rankings;

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

-- name: GetLeaders :many
SELECT * FROM ci6ndex.leaders;

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

-- name: GetDrafts :many
SELECT * FROM ci6ndex.drafts;

-- name: GetActiveDrafts :many
SELECT * FROM ci6ndex.drafts
WHERE active = true;

-- name: SubmitDraftPick :one
INSERT INTO ci6ndex.draft_picks
(
    draft_id, leader_id, user_id, offered
) VALUES (
    $1, $2, $3, $4
)
RETURNING *;

-- name: WipeTables :exec
DO $$
    DECLARE
        r RECORD;
    BEGIN
        FOR r IN (SELECT tablename FROM pg_tables WHERE schemaname = 'ci6ndex')
        LOOP
            IF EXISTS (SELECT 1 FROM pg_tables WHERE schemaname = 'ci6ndex' AND tablename = r.tablename) THEN
                EXECUTE 'TRUNCATE TABLE ci6ndex.' || quote_ident(r.tablename) || ' CASCADE;';
            END IF;
        END LOOP;
    END $$;