package pkg

import (
	"ci6ndex-bot/domain"
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DatabaseOperations struct {
	db      *pgxpool.Pool
	Queries *domain.Queries
}

func NewDBConnection(dbUrl string) (*DatabaseOperations, error) {
	conn, err := pgxpool.New(context.Background(), dbUrl)
	if err != nil {
		return nil, err
	}

	err = conn.Ping(context.Background())
	if err != nil {
		return nil, err
	}

	q := domain.New(conn)

	return &DatabaseOperations{db: conn, Queries: q}, nil
}

func (db DatabaseOperations) Health() error {
	err := db.db.Ping(context.Background())
	if err != nil {
		return err
	}
	return nil
}

func (db DatabaseOperations) Close() {
	db.db.Close()
}
