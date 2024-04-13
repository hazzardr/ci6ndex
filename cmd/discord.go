package cmd

import (
	"ci6ndex/internal"
	"fmt"
	"log/slog"

	"github.com/spf13/cobra"
)

var discordCmd = &cobra.Command{
	Use:   "discord",
	Short: "discord interaction tools for ci6ndex",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			err := cmd.Help()
			if err != nil {
				fmt.Println("Error: ", err)
			}
		}
	},
}

var startBot = &cobra.Command{
	Use:   "start",
	Short: "start an instance of the discord bot",
	Long: `
Starts the discord bot, which will listen for commands and respond to them.
`,
	Run: startBotFunc,
}

func startBotFunc(cmd *cobra.Command, args []string) {
	bot, err := internal.NewDiscordBot(db, config)
	if err != nil {
		slog.Error("Error creating discord bot", "error", err)
		return
	}
	err = bot.Start()
	if err != nil {
		slog.Error("Error starting discord bot", "error", err)
		return
	}
}

func init() {
	rootCmd.AddCommand(discordCmd)
	discordCmd.AddCommand(startBot)
}
