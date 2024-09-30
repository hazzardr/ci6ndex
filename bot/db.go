package main

import (
	"ci6ndex-bot/domain"
	"database/sql"
	"github.com/pkg/errors"

	_ "github.com/mattn/go-sqlite3"
)

type DatabaseOperations struct {
	Queries         *domain.Queries
	ReadConnection  *sql.DB
	WriteConnection *sql.DB
}

func NewDBConnection(dbUrl string) (*DatabaseOperations, error) {
	readConn, err := sql.Open("sqlite3", dbUrl)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to initialize read connection.")
	}
	writeConn, err := sql.Open("sqlite3", dbUrl)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to initialize write connection.")
	}

	// sqlite does not support multiple write connections
	writeConn.SetMaxOpenConns(1)
	return &DatabaseOperations{
		Queries:         domain.New(readConn),
		ReadConnection:  readConn,
		WriteConnection: writeConn,
	}, nil

}

func (db DatabaseOperations) Health() error {
	err := db.ReadConnection.Ping()
	if err != nil {
		return err
	}
	err = db.WriteConnection.Ping()
	if err != nil {
		return err
	}
	return nil
}

func (db DatabaseOperations) Close() error {
	err := db.ReadConnection.Close()
	if err != nil {
		return err
	}
	err = db.WriteConnection.Close()
	if err != nil {
		return err
	}
	return nil
}
