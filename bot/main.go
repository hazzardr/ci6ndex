/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"ci6ndex-bot/ci6ndex"
	"ci6ndex-bot/commands"
	"github.com/disgoorg/disgo/handler"
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
	r := handler.New()
	r.Command("/ping", commands.HandlePing)
	r.Command("/roll", commands.HandleRollCivs)

	err = bot.Configure(r)
	if err != nil {
		panic(err)
	}

	//bot.SyncCommands()
	defer func() {
		bot.GracefulShutdown()
	}()
	err = bot.Start()

	if err != nil {
		panic(err)
	}
	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-s
}
