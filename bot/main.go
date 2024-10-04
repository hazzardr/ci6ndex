/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"ci6ndex-bot/ci6ndex"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	config, err := ci6ndex.LoadConfig()
	if err != nil {
		panic(err)
	}
	deps, err := Initialize(config)
	if err != nil {
		panic(err)
	}

	bot := ci6ndex.New(*config, deps.db, *deps.logger)
	if err := bot.Configure(); err != nil {
		panic(err)
	}

	defer func() {
		bot.GracefulShutdown()
	}()
	bot.Start()

	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-s
}
