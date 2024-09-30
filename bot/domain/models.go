// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.24.0

package domain

import (
	"database/sql"
)

type Draft struct {
	ID            int64
	DraftStrategy string
	Active        bool
	Players       sql.NullString
}

type DraftStrategy struct {
	Name        string
	Description string
	PoolSize    int64
	Randomize   bool
}

type Leader struct {
	ID                 int64
	CivName            string
	LeaderName         string
	DiscordEmojiString sql.NullString
	Banned             bool
	Ranking            float64
}

type Offered struct {
	UserID  int64
	DraftID int64
	Offered sql.NullString
}
