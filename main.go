/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"ci6ndex/bot"
	"ci6ndex/ci6ndex"
	"embed"
	"github.com/disgoorg/disgo/handler"
	"os"
	"os/signal"
	"syscall"
)

//go:embed sql/migrations/*.sql
var embedMigrations embed.FS

func main() {
	config, err := bot.LoadConfig()
	if err != nil {
		panic(err)
	}
	c, err := ci6ndex.New(embedMigrations)
	if err != nil {
		panic(err)
	}

	b := bot.New(*config, c)
	err = bot.Configure(b, handler.New())
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
