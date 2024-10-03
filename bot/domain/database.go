package domain

import (
	"ci6ndex-bot/domain/generated"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"github.com/pkg/errors"
)

type DatabaseOperations struct {
	Queries  *generated.Queries
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
		Queries: generated.New(readConn),
	}, nil

}

func (db *DatabaseOperations) Health() error {
	return db.readConn.Ping()
}

func (db *DatabaseOperations) Close() {
	db.readConn.Close()
}
