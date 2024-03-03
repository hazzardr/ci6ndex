package cmd

import (
	"ci6ndex/internal"
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "ci6ndex",
	Short: "Civ 6 Index is a tool for drafting Civilization 6 leaders.",
	Long: `Civ 6 Index (CiVIndex) is primarily used to define drafts for Civilization games. 

It contains static data such as leaders + their civs, users who participate, rankings given by the users to the leaders.
This tool manages the above through the CLI, and the rankings themselves are managed through a Google Sheet.

This tool is also used to define the different draft strategies available, and the status of the draft itself.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

var db *internal.DatabaseOperations
var config *internal.AppConfig

func init() {
	db = internal.Initialize()
	config = internal.GetConfig()
}
