-- name: GetLeaders :many
SELECT * FROM leaders;

-- name: GetEligibleLeaders :many
SELECT * FROM leaders WHERE banned = false;

-- name: GetActiveDraft :one
SELECT * FROM drafts WHERE active = true;

-- name: GetOffersByDraftId :many
SELECT * FROM pool WHERE draft_id = ?;

-- name: GetPlayersFromDraft :many
SELECT p.*
FROM draft_registry dr
JOIN players p ON dr.player_id = p.id
WHERE dr.draft_id = ?;

-- name: GetPlayersFromActiveDraft :many
SELECT p.*
FROM draft_registry dr
JOIN players p ON dr.player_id = p.id
JOIN drafts d ON dr.draft_id = d.id
WHERE d.active = true;

-- name: GetPlayers :many
SELECT * FROM players;

-- name: GetPlayer :one
SELECT *
FROM players
WHERE id = ?;

-- name: GetLeadersByLimitAndOffset :many
SELECT *
FROM leaders l
ORDER BY l.leader_name
LIMIT ? OFFSET ?;
