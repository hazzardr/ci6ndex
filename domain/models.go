// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.24.0

package domain

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type Ci6ndexDraft struct {
	ID            int64
	DraftStrategy string
	Active        bool
}

type Ci6ndexDraftPick struct {
	ID       int64
	DraftID  int64
	UserID   int64
	LeaderID pgtype.Int8
	// The civs that were offered to the user.
	Offered []int64
}

// The strategies that can be used to draft a civ
type Ci6ndexDraftStrategy struct {
	Name        string
	Description string
	PoolSize    int32
	Randomize   bool
	// Specific rules that this draft has to follow.
	Rules []byte
}

type Ci6ndexGame struct {
	ID        int64
	DraftID   int64
	StartDate pgtype.Date
	EndDate   pgtype.Date
	GameStats []byte
}

type Ci6ndexLeader struct {
	ID         int64
	CivName    string
	LeaderName string
	IconUrl    pgtype.Text
	Banned     bool
}

type Ci6ndexOffered struct {
	UserID  int64
	DraftID int64
	Offered []byte
}

type Ci6ndexRanking struct {
	ID       int64
	UserID   int64
	Tier     float64
	LeaderID int64
}

type Ci6ndexStat struct {
	ID     int64
	Stats  []byte
	UserID int64
	GameID int64
}

type Ci6ndexUser struct {
	ID          int64
	DiscordName string
	Name        string
}
