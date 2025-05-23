// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.24.0
// source: writes.sql

package generated

import (
	"context"
	"database/sql"
)

const addPlayer = `-- name: AddPlayer :exec
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
    discord_avatar = EXCLUDED.discord_avatar
`

type AddPlayerParams struct {
	ID            int64
	Username      string
	GlobalName    sql.NullString
	DiscordAvatar sql.NullString
}

func (q *Queries) AddPlayer(ctx context.Context, arg AddPlayerParams) error {
	_, err := q.db.ExecContext(ctx, addPlayer,
		arg.ID,
		arg.Username,
		arg.GlobalName,
		arg.DiscordAvatar,
	)
	return err
}

const addPlayerToDraft = `-- name: AddPlayerToDraft :one
INSERT INTO draft_registry (
    draft_id,
    player_id
) VALUES (
    ?, ?
) RETURNING player_id, draft_id
`

type AddPlayerToDraftParams struct {
	DraftID  int64
	PlayerID int64
}

func (q *Queries) AddPlayerToDraft(ctx context.Context, arg AddPlayerToDraftParams) (DraftRegistry, error) {
	row := q.db.QueryRowContext(ctx, addPlayerToDraft, arg.DraftID, arg.PlayerID)
	var i DraftRegistry
	err := row.Scan(&i.PlayerID, &i.DraftID)
	return i, err
}

const addPool = `-- name: AddPool :exec
INSERT INTO pool (
    player_id,
    draft_id,
    leader
) VALUES (
    ?, ?, ?
)
`

type AddPoolParams struct {
	PlayerID int64
	DraftID  int64
	Leader   int64
}

func (q *Queries) AddPool(ctx context.Context, arg AddPoolParams) error {
	_, err := q.db.ExecContext(ctx, addPool, arg.PlayerID, arg.DraftID, arg.Leader)
	return err
}

const createActiveDraft = `-- name: CreateActiveDraft :one
INSERT INTO drafts (
    active
) VALUES (true) RETURNING id, active
`

func (q *Queries) CreateActiveDraft(ctx context.Context) (Draft, error) {
	row := q.db.QueryRowContext(ctx, createActiveDraft)
	var i Draft
	err := row.Scan(&i.ID, &i.Active)
	return i, err
}

const deletePoolForPlayer = `-- name: DeletePoolForPlayer :exec
DELETE FROM pool
       WHERE player_id = ?
        AND draft_id = ?
`

type DeletePoolForPlayerParams struct {
	PlayerID int64
	DraftID  int64
}

func (q *Queries) DeletePoolForPlayer(ctx context.Context, arg DeletePoolForPlayerParams) error {
	_, err := q.db.ExecContext(ctx, deletePoolForPlayer, arg.PlayerID, arg.DraftID)
	return err
}

const deletePoolsForDraftId = `-- name: DeletePoolsForDraftId :exec
DELETE FROM pool WHERE draft_id = ?
`

func (q *Queries) DeletePoolsForDraftId(ctx context.Context, draftID int64) error {
	_, err := q.db.ExecContext(ctx, deletePoolsForDraftId, draftID)
	return err
}

const removePlayersFromDraft = `-- name: RemovePlayersFromDraft :exec
DELETE FROM draft_registry WHERE draft_id = ?
`

func (q *Queries) RemovePlayersFromDraft(ctx context.Context, draftID int64) error {
	_, err := q.db.ExecContext(ctx, removePlayersFromDraft, draftID)
	return err
}

const returnOffering = `-- name: ReturnOffering :exec
DELETE FROM pool
   WHERE player_id = ?
    AND draft_id = ?
`

type ReturnOfferingParams struct {
	PlayerID int64
	DraftID  int64
}

func (q *Queries) ReturnOffering(ctx context.Context, arg ReturnOfferingParams) error {
	_, err := q.db.ExecContext(ctx, returnOffering, arg.PlayerID, arg.DraftID)
	return err
}
