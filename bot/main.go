/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"ci6ndex-bot/ci6ndex"
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
	bot.Start()
}
