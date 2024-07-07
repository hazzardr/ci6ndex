package cmd

import (
	"ci6ndex-bot/pkg"
	"fmt"
	"github.com/spf13/viper"
	"log/slog"
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

// Execute adds all child slashCommands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

var db *pkg.DatabaseOperations
var config *pkg.AppConfig

func init() {
	err := initializeConfig()
	if err != nil {
		slog.Error("Error initializing config", "error", err)
		return
	}
	slog.Info("initializing db...")
	db, err = pkg.NewDBConnection(config.DatabaseUrl)
	slog.Info("done!")
	if err != nil {
		slog.Error("Error initializing db", "error", err)
		return
	}
}

func initializeConfig() error {
	viper.SetConfigFile(".env")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("failed to load configuration, error=%w", err))
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		panic(fmt.Errorf("failed to load configuration, error=%w", err))
	}
	return nil
}
