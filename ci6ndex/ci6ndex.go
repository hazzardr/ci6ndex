package ci6ndex

import (
	"embed"
	"github.com/charmbracelet/log"
	"os"
	"time"
)

type Ci6ndex struct {
	Logger      *log.Logger
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

	return &Ci6ndex{
		Logger:      logger,
		Connections: connections,
		Path:        "./data/",
	}, nil
}
