package domain

import (
	"ci6ndex-bot/domain/generated"
	"database/sql"
	"embed"
	"fmt"
	"github.com/charmbracelet/log"
	_ "github.com/mattn/go-sqlite3"
	"github.com/pkg/errors"
	"os"
	"strconv"
)

type DB struct {
	ReadConn  *sql.DB
	WriteConn *sql.DB
	Queries   *generated.Queries
	Writes    *generated.Queries
}

type DatabaseOperations struct {
	logger      *log.Logger
	connections map[uint64]*DB
	path        string
}

func NewDBOperations(embedFileSystem embed.FS, logger *log.Logger) *DatabaseOperations {
	connections := make(map[uint64]*DB)
	err := configureMigrationTool(embedFileSystem)
	if err != nil {
		logger.Error(
			"Failed to configure migrations. Will not be able to create any new databases!",
			"error", err,
		)
	}
	return &DatabaseOperations{
		logger:      logger,
		connections: connections,
		path:        "./data/",
	}
}

func (dbo *DatabaseOperations) openNewConnection(guildId uint64) (*DB, error) {
	dbUrl := "file:" + dbo.path + strconv.FormatUint(guildId, 10) + ".db"
	_, err := os.Stat(strconv.FormatUint(guildId, 10) + ".db")
	if os.IsNotExist(err) {
		dbo.logger.Info("no database exists, creating new one...", "guildId", guildId)
		db, err := sql.Open("sqlite3", dbUrl)
		if err != nil {
			return nil, errors.Wrap(
				err,
				fmt.Sprintf("failed to create new database file for guild %d", guildId),
			)
		}
		err = dbo.migrateUp(db)
		if err != nil {
			return nil, errors.Wrap(
				err,
				fmt.Sprintf("failed to create new database file for guild %d", guildId),
			)
		}
	}

	readConn, err := sql.Open("sqlite3", dbUrl)
	if err != nil {
		return nil, errors.Wrap(err, "failed to initialize database connection.")
	}
	writeConn, err := sql.Open("sqlite3", dbUrl)
	if err != nil {
		return nil, errors.Wrap(err, "failed to initialize database connection.")
	}

	// sqlite does not support multiple write connections
	writeConn.SetMaxOpenConns(1)
	return &DB{
		ReadConn:  readConn,
		WriteConn: writeConn,
		Queries:   generated.New(readConn),
		Writes:    generated.New(writeConn),
	}, nil
}

func (dbo *DatabaseOperations) getDB(guildId uint64) (*DB, error) {
	var db *DB
	db, exists := dbo.connections[guildId]
	if !exists {
		var err error
		db, err = dbo.openNewConnection(guildId)
		if err != nil {
			return nil, err
		}
		dbo.connections[guildId] = db
	}
	return db, nil
}

func (dbo *DatabaseOperations) Health() []error {
	var errs = make([]error, 0)
	for _, db := range dbo.connections {
		if err := db.ReadConn.Ping(); err != nil {
			errs = append(errs, err)
		}
		if err := db.WriteConn.Ping(); err != nil {
			errs = append(errs, err)
		}
	}
	return errs
}

func (dbo *DatabaseOperations) Close() {
	for _, db := range dbo.connections {
		db.ReadConn.Close()
		db.WriteConn.Close()
	}
}

func ResolveOptionalString(s *string) sql.NullString {
	if s == nil {
		return sql.NullString{
			Valid: false,
		}
	}
	return sql.NullString{
		String: *s,
		Valid:  true,
	}
}
