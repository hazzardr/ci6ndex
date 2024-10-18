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