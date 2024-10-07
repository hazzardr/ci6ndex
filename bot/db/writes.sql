-- name: CreateDraft :one
INSERT INTO drafts (
    active,
    players
) VALUES (?, ?) RETURNING *;