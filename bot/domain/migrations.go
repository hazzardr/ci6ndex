package domain

import (
	"database/sql"
	"embed"
	"github.com/pressly/goose/v3"
)

func configureMigrationTool(embedMigrations embed.FS) error {
	goose.SetBaseFS(embedMigrations)
	err := goose.SetDialect("sqlite3")
	if err != nil {
		return err
	}
	return nil
}

func (dbo *DatabaseOperations) migrateUp(db *sql.DB) error {
	dbo.logger.Info("running migrations...", "database", db)
	err := goose.Up(db, "db/migrations")
	if err != nil {
		return err
	}
	dbo.logger.Debug("migrations ran successfully", "database", db)
	return nil
}
