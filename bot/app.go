package main

import (
	"ci6ndex-bot/ci6ndex"
	"ci6ndex-bot/domain"
	"github.com/charmbracelet/log"
	"os"
	"time"
)

type Dependencies struct {
	db     *domain.DatabaseOperations
	config *ci6ndex.AppConfig
	logger *log.Logger
}

func Initialize(c *ci6ndex.AppConfig) (*Dependencies, error) {
	db, err := domain.NewDBConnection(c.DatabaseUrl)
	if err != nil {
		panic(err)
	}
	err = db.Health()
	if err != nil {
		panic(err)
	}
	logger := log.NewWithOptions(os.Stderr, log.Options{
		ReportCaller:    true,
		ReportTimestamp: true,
		TimeFormat:      time.Kitchen,
	})

	return &Dependencies{
		db:     db,
		config: c,
		logger: logger,
	}, nil
}

func (a *Dependencies) Serve() {
	a.db.Close()
}
