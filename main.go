/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"ci6ndex/bot"
	"ci6ndex/ci6ndex"
	"ci6ndex/cmd"
	"embed"
	"github.com/charmbracelet/log"
	"log/slog"
	"os"
	"time"
)

//go:embed sql/migrations/*.sql
var embedMigrations embed.FS

func main() {
	config, err := loadConfig()
	configureLog()
	if err != nil {
		slog.Error("Failed to load config", slog.Any("err", err))
		os.Exit(1)
	}
	c, err := ci6ndex.New(embedMigrations)
	if err != nil {
		slog.Error("failed to load ci6ndex", slog.Any("err", err))
		os.Exit(1)
	}
	b := bot.New(
		c,
		config.DiscordToken,
		config.GuildIDs,
	)
	err = b.Configure()
	if err != nil {
		slog.Error("Failed to configure bot", slog.Any("err", err))
		os.Exit(1)
	}

	if err := cmd.Exec(b); err != nil {
		slog.Error("Failed to execute command", slog.Any("err", err))
		os.Exit(1)
	}
}

func configureLog() {
	logger := log.NewWithOptions(os.Stderr, log.Options{
		ReportCaller:    true,
		ReportTimestamp: true,
		TimeFormat:      time.Kitchen,
	})
	slog.SetDefault(slog.New(logger))
}
