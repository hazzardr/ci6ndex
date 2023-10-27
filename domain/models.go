// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.23.0

package domain

import ()

type Ci6ndexGame struct {
	ID      int64
	UserIds []int32
}

type Ci6ndexLeader struct {
	ID         int64
	CivName    string
	LeaderName string
}

type Ci6ndexRanking struct {
	ID       int64
	UserID   int32
	Tier     int32
	LeaderID int32
}

type Ci6ndexStat struct {
	ID     int64
	Stats  []byte
	UserID int32
	GameID int32
}

type Ci6ndexUser struct {
	ID          int64
	DiscordName string
	Name        string
}