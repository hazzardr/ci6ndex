-- name: GetLeaders :one
SELECT * FROM leaders;

-- name: GetActiveDraft :one
SELECT * FROM drafts WHERE active = true;