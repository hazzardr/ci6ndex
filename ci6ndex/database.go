package ci6ndex

import (
	"ci6ndex/ci6ndex/generated"
	"database/sql"
	"embed"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"github.com/pkg/errors"
	"github.com/pressly/goose/v3"
	"log/slog"
	"os"
	"strconv"
)

type DB struct {
	readConn  *sql.DB
	writeConn *sql.DB
	Queries   *generated.Queries
	Writes    *generated.Queries
}

func (c *Ci6ndex) openNewConnection(path string, guildId uint64) (*DB, error) {
	dbUrl := "file:" + path + strconv.FormatUint(guildId, 10) + ".db"

	_, err := os.Stat(c.Path + strconv.FormatUint(guildId, 10) + ".db")

	if os.IsNotExist(err) {
		slog.Info("no database exists, creating new one...", "guildId", guildId)
		db, err := sql.Open("sqlite3", dbUrl)
		if err != nil {
			return nil, errors.Wrap(
				err,
				fmt.Sprintf("failed to create new database file for guild %d", guildId),
			)
		}

		err = c.migrateUp(db)
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

	// sqlite does not support multiple write Connections
	writeConn.SetMaxOpenConns(1)
	return &DB{
		readConn:  readConn,
		writeConn: writeConn,
		Queries:   generated.New(readConn),
		Writes:    generated.New(writeConn),
	}, nil
}

func (c *Ci6ndex) getDB(guildId uint64) (*DB, error) {
	var db *DB
	db, exists := c.Connections[guildId]
	if !exists {
		var err error
		db, err = c.openNewConnection(c.Path, guildId)
		if err != nil {
			return nil, err
		}
		c.Connections[guildId] = db
	}
	return db, nil
}

func (c *Ci6ndex) Health() []error {
	var errs = make([]error, 0)
	for _, db := range c.Connections {
		if err := db.readConn.Ping(); err != nil {
			errs = append(errs, err)
		}
		if err := db.writeConn.Ping(); err != nil {
			errs = append(errs, err)
		}
	}
	return errs
}

func (c *Ci6ndex) Close() {
	for _, db := range c.Connections {
		db.readConn.Close()
		db.writeConn.Close()
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

func ConfigureMigrationTool(embedMigrations embed.FS) error {
	goose.SetBaseFS(embedMigrations)
	err := goose.SetDialect("sqlite3")
	if err != nil {
		return err
	}
	return nil
}

func (c *Ci6ndex) migrateUp(db *sql.DB) error {
	err := goose.Up(db, "sql/migrations")
	if err != nil {
		return err
	}
	return nil
}
