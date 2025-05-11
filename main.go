/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"ci6ndex/bot"
	"ci6ndex/ci6ndex"
	"embed"
	"github.com/charmbracelet/log"
	"github.com/disgoorg/disgo/handler"
	"os"
	"os/signal"
	"syscall"
	"time"
)

//go:embed sql/migrations/*.sql
var embedMigrations embed.FS

func Initialize() (*ci6ndex.Ci6ndex, error) {
	logger := log.NewWithOptions(os.Stderr, log.Options{
		ReportCaller:    true,
		ReportTimestamp: true,
		TimeFormat:      time.Kitchen,
	})

	connections := make(map[uint64]*ci6ndex.DB)
	err := ci6ndex.ConfigureMigrationTool(embedMigrations)
	if err != nil {
		logger.Error(
			"Failed to configure migrations. Will not be able to create any new databases!",
			"error", err,
		)
	}

	return &ci6ndex.Ci6ndex{
		Logger:      logger,
		Connections: connections,
		Path:        "./data/",
	}, nil
}

func main() {
	config, err := bot.LoadConfig()
	if err != nil {
		panic(err)
	}
	c, err := Initialize()
	if err != nil {
		panic(err)
	}

	b := bot.New(*config, c, *c.Logger)
	commandHandler := handler.New()
	commandHandler.Command("/ping", bot.HandlePing)
	commandHandler.Command("/roll", bot.HandleRollCivs(b))
	commandHandler.SelectMenuComponent("/select-player", bot.HandlePlayerSelect(b))
	commandHandler.SelectMenuComponent("/select-reroll-player", bot.HandlePlayerSelectReRoll(b))
	commandHandler.ButtonComponent("/confirm-roll", bot.HandleConfirmRoll(b))
	err = bot.Configure(b, commandHandler)
	if err != nil {
		panic(err)
	}

	defer bot.GracefulShutdown(b)
	err = bot.Start(b)
	if err != nil {
		panic(err)
	}

	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-s
}
