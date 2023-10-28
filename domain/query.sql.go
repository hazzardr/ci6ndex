// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.23.0
// source: query.sql

package domain

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const createDraft = `-- name: CreateDraft :one
INSERT INTO ci6ndex.drafts
(
    draft_strategy
) VALUES (
    $1
)
RETURNING id, draft_strategy, active
`

func (q *Queries) CreateDraft(ctx context.Context, draftStrategy string) (Ci6ndexDraft, error) {
	row := q.db.QueryRow(ctx, createDraft, draftStrategy)
	var i Ci6ndexDraft
	err := row.Scan(&i.ID, &i.DraftStrategy, &i.Active)
	return i, err
}

const createDraftStrategy = `-- name: CreateDraftStrategy :one
INSERT INTO ci6ndex.draft_strategies
(
    name, description
) VALUES (
    $1, $2
)
RETURNING name, description, rules
`

type CreateDraftStrategyParams struct {
	Name        string
	Description string
}

func (q *Queries) CreateDraftStrategy(ctx context.Context, arg CreateDraftStrategyParams) (Ci6ndexDraftStrategy, error) {
	row := q.db.QueryRow(ctx, createDraftStrategy, arg.Name, arg.Description)
	var i Ci6ndexDraftStrategy
	err := row.Scan(&i.Name, &i.Description, &i.Rules)
	return i, err
}

const createRanking = `-- name: CreateRanking :one
INSERT INTO ci6ndex.rankings
(
    user_id, tier, leader_id
) VALUES (
    $1, $2, $3
)
RETURNING id, user_id, tier, leader_id
`

type CreateRankingParams struct {
	UserID   int64
	Tier     float64
	LeaderID int64
}

func (q *Queries) CreateRanking(ctx context.Context, arg CreateRankingParams) (Ci6ndexRanking, error) {
	row := q.db.QueryRow(ctx, createRanking, arg.UserID, arg.Tier, arg.LeaderID)
	var i Ci6ndexRanking
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Tier,
		&i.LeaderID,
	)
	return i, err
}

type CreateRankingsParams struct {
	UserID   int64
	Tier     float64
	LeaderID int64
}

const createUser = `-- name: CreateUser :one
INSERT INTO ci6ndex.users
(
    name, discord_name
) VALUES (
    $1, $2
)
RETURNING id, discord_name, name
`

type CreateUserParams struct {
	Name        string
	DiscordName string
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (Ci6ndexUser, error) {
	row := q.db.QueryRow(ctx, createUser, arg.Name, arg.DiscordName)
	var i Ci6ndexUser
	err := row.Scan(&i.ID, &i.DiscordName, &i.Name)
	return i, err
}

const deleteRankings = `-- name: DeleteRankings :exec
DELETE FROM ci6ndex.rankings
RETURNING id, user_id, tier, leader_id
`

func (q *Queries) DeleteRankings(ctx context.Context) error {
	_, err := q.db.Exec(ctx, deleteRankings)
	return err
}

const getActiveDrafts = `-- name: GetActiveDrafts :many
SELECT id, draft_strategy, active FROM ci6ndex.drafts
WHERE active = true
`

func (q *Queries) GetActiveDrafts(ctx context.Context) ([]Ci6ndexDraft, error) {
	rows, err := q.db.Query(ctx, getActiveDrafts)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Ci6ndexDraft
	for rows.Next() {
		var i Ci6ndexDraft
		if err := rows.Scan(&i.ID, &i.DraftStrategy, &i.Active); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getDraft = `-- name: GetDraft :one
SELECT id, draft_strategy, active FROM ci6ndex.drafts
WHERE id = $1
LIMIT 1
`

func (q *Queries) GetDraft(ctx context.Context, id int64) (Ci6ndexDraft, error) {
	row := q.db.QueryRow(ctx, getDraft, id)
	var i Ci6ndexDraft
	err := row.Scan(&i.ID, &i.DraftStrategy, &i.Active)
	return i, err
}

const getDraftStrategies = `-- name: GetDraftStrategies :many
SELECT name, description, rules FROM ci6ndex.draft_strategies
`

func (q *Queries) GetDraftStrategies(ctx context.Context) ([]Ci6ndexDraftStrategy, error) {
	rows, err := q.db.Query(ctx, getDraftStrategies)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Ci6ndexDraftStrategy
	for rows.Next() {
		var i Ci6ndexDraftStrategy
		if err := rows.Scan(&i.Name, &i.Description, &i.Rules); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getDraftStrategy = `-- name: GetDraftStrategy :one
SELECT name, description, rules FROM ci6ndex.draft_strategies
WHERE name = $1
LIMIT 1
`

func (q *Queries) GetDraftStrategy(ctx context.Context, name string) (Ci6ndexDraftStrategy, error) {
	row := q.db.QueryRow(ctx, getDraftStrategy, name)
	var i Ci6ndexDraftStrategy
	err := row.Scan(&i.Name, &i.Description, &i.Rules)
	return i, err
}

const getLeader = `-- name: GetLeader :one
SELECT id, civ_name, leader_name FROM ci6ndex.leaders
         WHERE id = $1
         LIMIT 1
`

func (q *Queries) GetLeader(ctx context.Context, id int64) (Ci6ndexLeader, error) {
	row := q.db.QueryRow(ctx, getLeader, id)
	var i Ci6ndexLeader
	err := row.Scan(&i.ID, &i.CivName, &i.LeaderName)
	return i, err
}

const getLeaderByNameAndCiv = `-- name: GetLeaderByNameAndCiv :one
SELECT id, civ_name, leader_name FROM ci6ndex.leaders
WHERE leader_name = $1
AND civ_name = $2
LIMIT 1
`

type GetLeaderByNameAndCivParams struct {
	LeaderName string
	CivName    string
}

func (q *Queries) GetLeaderByNameAndCiv(ctx context.Context, arg GetLeaderByNameAndCivParams) (Ci6ndexLeader, error) {
	row := q.db.QueryRow(ctx, getLeaderByNameAndCiv, arg.LeaderName, arg.CivName)
	var i Ci6ndexLeader
	err := row.Scan(&i.ID, &i.CivName, &i.LeaderName)
	return i, err
}

const getUser = `-- name: GetUser :one
SELECT id, discord_name, name FROM ci6ndex.users WHERE id = $1 LIMIT 1
`

func (q *Queries) GetUser(ctx context.Context, id int64) (Ci6ndexUser, error) {
	row := q.db.QueryRow(ctx, getUser, id)
	var i Ci6ndexUser
	err := row.Scan(&i.ID, &i.DiscordName, &i.Name)
	return i, err
}

const getUserByDiscordName = `-- name: GetUserByDiscordName :one
SELECT id, discord_name, name FROM ci6ndex.users WHERE discord_name = $1 LIMIT 1
`

func (q *Queries) GetUserByDiscordName(ctx context.Context, discordName string) (Ci6ndexUser, error) {
	row := q.db.QueryRow(ctx, getUserByDiscordName, discordName)
	var i Ci6ndexUser
	err := row.Scan(&i.ID, &i.DiscordName, &i.Name)
	return i, err
}

const getUserByName = `-- name: GetUserByName :one
SELECT id, discord_name, name FROM ci6ndex.users WHERE name = $1 LIMIT 1
`

func (q *Queries) GetUserByName(ctx context.Context, name string) (Ci6ndexUser, error) {
	row := q.db.QueryRow(ctx, getUserByName, name)
	var i Ci6ndexUser
	err := row.Scan(&i.ID, &i.DiscordName, &i.Name)
	return i, err
}

const submitDraftPick = `-- name: SubmitDraftPick :one
INSERT INTO ci6ndex.draft_picks
(
    draft_id, leader_id, user_id, offered
) VALUES (
    $1, $2, $3, $4
)
RETURNING id, draft_id, user_id, leader_id, offered
`

type SubmitDraftPickParams struct {
	DraftID  int64
	LeaderID pgtype.Int8
	UserID   int64
	Offered  []byte
}

func (q *Queries) SubmitDraftPick(ctx context.Context, arg SubmitDraftPickParams) (Ci6ndexDraftPick, error) {
	row := q.db.QueryRow(ctx, submitDraftPick,
		arg.DraftID,
		arg.LeaderID,
		arg.UserID,
		arg.Offered,
	)
	var i Ci6ndexDraftPick
	err := row.Scan(
		&i.ID,
		&i.DraftID,
		&i.UserID,
		&i.LeaderID,
		&i.Offered,
	)
	return i, err
}
