/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package main

import "log/slog"

func main() {
	db, err := NewDBConnection("file:bot.db")
	if err != nil {
		panic(err)
	}
	err = db.Health()
	if err != nil {
		panic(err)
	}
	_, err = loadConfig()
	if err != nil {
		panic(err)
	}
	slog.Info("Initialization successful")

	db.Close()
}
