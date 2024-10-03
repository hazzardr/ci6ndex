package main

import (
	"ci6ndex-bot/domain"
	"database/sql"
	"github.com/pkg/errors"

	_ "github.com/mattn/go-sqlite3"
)

type DatabaseOperations struct {
	Queries  *domain.Queries
	readConn *sql.DB
}

func NewDBConnection(dbUrl string) (*DatabaseOperations, error) {
	readConn, err := sql.Open("sqlite3", dbUrl)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to initialize read connection.")
	}
	//writeConn, err := sql.Open("sqlite3", dbUrl)
	//if err != nil {
	//	return nil, errors.Wrap(err, "Failed to initialize write connection.")
	//}
	//
	//// sqlite does not support multiple write connections
	//writeConn.SetMaxOpenConns(1)
	return &DatabaseOperations{
		Queries: domain.New(readConn),
	}, nil

}

func (db *DatabaseOperations) Health() error {
	return db.readConn.Ping()
}

func (db *DatabaseOperations) Close() {
	db.readConn.Close()
}
