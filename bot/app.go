package main

import (
	"ci6ndex-bot/ci6ndex"
	"ci6ndex-bot/domain"
	"embed"
	"github.com/charmbracelet/log"
	"os"
	"time"
)

//go:embed db/migrations/*.sql
var embedMigrations embed.FS

type Dependencies struct {
	db     *domain.DatabaseOperations
	config *ci6ndex.AppConfig
	logger *log.Logger
}

func Initialize(c *ci6ndex.AppConfig) (*Dependencies, error) {

	logger := log.NewWithOptions(os.Stderr, log.Options{
		ReportCaller:    true,
		ReportTimestamp: true,
		TimeFormat:      time.Kitchen,
	})
	db := domain.NewDBOperations(embedMigrations, logger)

	return &Dependencies{
		db:     db,
		config: c,
		logger: logger,
	}, nil
}

func (a *Dependencies) Serve() {
	a.db.Close()
}
