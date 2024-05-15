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

-- name: GetRankingsByLeader :many
SELECT * FROM ci6ndex.rankings
WHERE leader_id = $1;

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

-- name: CreateLeaders :copyfrom
INSERT INTO ci6ndex.leaders
(
    leader_name, civ_name
) VALUES (
    $1, $2
);

-- name: UpdateLeaderTier :one
UPDATE ci6ndex.leaders
SET tier = $2
WHERE id = $1
RETURNING *;

-- name: CreateDraftStrategy :one
INSERT INTO ci6ndex.draft_strategies
(
    name, description, pool_size, randomize
) VALUES (
    $1, $2, $3, $4
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

-- name: CreateActiveDraft :one
INSERT INTO ci6ndex.drafts
(
    draft_strategy, active
) VALUES (
    $1, true
)
RETURNING *;

-- name: AddPlayersToActiveDraft :one
UPDATE ci6ndex.drafts
SET players = $1
where active=true
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

-- name: CancelActiveDrafts :many
UPDATE ci6ndex.drafts
SET active = false
WHERE active = true
RETURNING *;

-- name: SubmitDraftPick :one
INSERT INTO ci6ndex.draft_picks
(
    draft_id, leader_id, user_id
) VALUES (
    $1, $2, $3
) ON CONFLICT (draft_id, user_id)
DO UPDATE SET leader_id = $2
RETURNING *;

-- name: RemoveDraftPick :one
DELETE FROM ci6ndex.draft_picks
WHERE draft_id = $1
  AND user_id = $2
RETURNING *;

-- name: GetDraftPicksForDraft :many
SELECT * FROM ci6ndex.draft_picks
WHERE draft_id = $1;

-- name: GetDenormalizedDraftPicksForDraft :many
SELECT u.discord_name, l.leader_name, l.civ_name, l.discord_emoji_string, dp.draft_id, g.start_date
FROM ci6ndex.draft_picks dp
JOIN ci6ndex.games g on dp.draft_id = g.draft_id
JOIN ci6ndex.leaders l ON dp.leader_id = l.id
JOIN ci6ndex.users u ON dp.user_id = u.id
WHERE dp.draft_id = $1;

-- name: WriteOffered :one
INSERT INTO ci6ndex.offered
(
    draft_id, user_id, offered
) VALUES (
    $1, $2, $3
) ON CONFLICT (draft_id, user_id)
DO UPDATE SET offered = $3
RETURNING *;

-- name: ReadOfferedForUser :one
SELECT * FROM ci6ndex.offered
WHERE draft_id = $1
  AND user_id = $2;

-- name: ReadOffer :many
SELECT * FROM ci6ndex.offered
WHERE draft_id = $1;

-- name: DeleteOffered :exec
DELETE FROM ci6ndex.offered
WHERE draft_id = $1
  AND user_id = $2;

-- name: ClearOffered :exec
TRUNCATE TABLE ci6ndex.offered;

-- name: GetDraftPicksForUserFromLastNGames :many
SELECT * FROM ci6ndex.draft_picks
WHERE user_id = (
    SELECT id FROM ci6ndex.users
    WHERE discord_name = $1
)
AND ci6ndex.draft_picks.draft_id IN (
   SELECT ci6ndex.games.draft_id FROM ci6ndex.games
    ORDER BY ci6ndex.games.start_date DESC LIMIT $2
);

-- name: CreateGameFromDraft :one
INSERT INTO ci6ndex.games
(
    draft_id, start_date
) VALUES (
    $1, $2
)
RETURNING *;

-- name: GetGameByDraftID :one
SELECT * FROM ci6ndex.games
WHERE draft_id = $1;

-- name: UpdateGameFromDraftId :exec
UPDATE ci6ndex.games
SET start_date = $2
WHERE draft_id = $1
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