package main

import (
	"ci6ndex-bot/domain"
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type Database struct {
	Queries    *domain.Queries
	Connection *sql.DB
}

func NewDBConnection(dbUrl string) (*Database, error) {
	conn, err := sql.Open("sqlite3", dbUrl)
	if err != nil {
		return nil, err
	}
	return &Database{
		Queries:    domain.New(conn),
		Connection: conn,
	}, nil

}

func (db Database) Health() error {
	err := db.Connection.Ping()
	if err != nil {
		return err
	}
	return nil
}

func (db Database) Close() error {
	return db.Connection.Close()
}
