-- name: CreateActiveDraft :one
INSERT INTO drafts (
    active
) VALUES (true) RETURNING *;

-- name: RemovePlayersFromDraft :exec
DELETE FROM draft_registry WHERE draft_id = ?;

-- name: AddPlayerToDraft :one
INSERT INTO draft_registry (
    draft_id,
    player_id
) VALUES (
    ?, ?
) RETURNING *;

-- name: AddPlayer :exec
INSERT INTO players (
    id,
    username,
    global_name,
    discord_avatar
) VALUES (
    ?, ?, ?, ?
) ON CONFLICT (id) DO UPDATE SET
    username = EXCLUDED.username,
    global_name = EXCLUDED.global_name,
    discord_avatar = EXCLUDED.discord_avatar;

-- name: AddPool :exec
INSERT INTO pool (
    player_id,
    draft_id,
    leader
) VALUES (
    ?, ?, ?
);

-- name: DeletePoolsForDraftId :exec
DELETE FROM pool WHERE draft_id = ?;

-- name: DeletePoolForPlayer :exec
DELETE FROM pool
       WHERE player_id = ?
        AND draft_id = ?;

-- name: ReturnOffering :exec
DELETE FROM pool
   WHERE player_id = ?
    AND draft_id = ?;

-- name: SubmitRankForPlayer :exec
INSERT INTO ranks (player_id, leader_id, tier)
VALUES (?, ?, ?)
ON CONFLICT (leader_id, player_id)
DO UPDATE SET
    tier = excluded.tier,
    updated_at = CURRENT_TIMESTAMP
;

-- name: UpdateLeaderTier :exec
UPDATE leaders 
SET tier = ?
WHERE id = ?;
