/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"log/slog"
)

func main() {
	_, err := InitNewApp()
	if err != nil {
		panic(err)
	}
	slog.Info("initialization successful")
}
