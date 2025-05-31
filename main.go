/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"ci6ndex/cmd"
	"embed"
)

//go:embed sql/migrations/*.sql
var embedMigrations embed.FS

func main() {
	cmd.Exec(embedMigrations)
}
