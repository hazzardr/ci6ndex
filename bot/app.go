package main

import (
	"ci6ndex-bot/domain"
	"github.com/charmbracelet/log"
	"os"
	"time"
)

type App struct {
	db     *domain.DatabaseOperations
	config *AppConfig
	logger *log.Logger
}

func InitNewApp(c *AppConfig) (*App, error) {
	db, err := domain.NewDBConnection(c.DatabaseUrl)
	if err != nil {
		panic(err)
	}
	logger := log.NewWithOptions(os.Stderr, log.Options{
		ReportCaller:    true,
		ReportTimestamp: true,
		TimeFormat:      time.Kitchen,
	})

	return &App{
		db:     db,
		config: c,
		logger: logger,
	}, nil
}

func (a *App) Serve() {
	a.db.Close()
}
