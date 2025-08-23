package ci6ndex

import (
	"embed"
	"os"
	"time"

	"github.com/charmbracelet/log"
)

type Ci6ndex struct {
	Connections map[uint64]*DB
	Path        string
}

func New(embedMigrations embed.FS) (*Ci6ndex, error) {
	logger := log.NewWithOptions(os.Stderr, log.Options{
		ReportCaller:    true,
		ReportTimestamp: true,
		TimeFormat:      time.Kitchen,
	})

	connections := make(map[uint64]*DB)
	err := ConfigureMigrationTool(embedMigrations)
	if err != nil {
		logger.Error(
			"Failed to configure migrations. Will not be able to create any new databases!",
			"error", err,
		)
	}

	// Ensure data directory exists
	dataPath := "./data/"
	if _, err := os.Stat(dataPath); os.IsNotExist(err) {
		logger.Info("Creating data directory", "path", dataPath)
		if err := os.MkdirAll(dataPath, 0755); err != nil {
			logger.Error("Failed to create data directory", "error", err)
			return nil, err
		}
	}

	return &Ci6ndex{
		Connections: connections,
		Path:        dataPath,
	}, nil
}
