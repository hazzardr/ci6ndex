package cmd

import (
	"ci6ndex/internal/discord"
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
Starts the discord bot, which will listen for slashCommands and respond to them.
`,
	Run: startBotFunc,
}

func startBotFunc(cmd *cobra.Command, args []string) {
	bot, err := discord.NewBot(db, config)
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

// -------- BEGIN SLASH COMMANDS --------
var guildId string

var slashCommands = &cobra.Command{
	Use:   "commands",
	Short: "manage slash (/) commands for the discord bot",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			err := cmd.Help()
			if err != nil {
				fmt.Println("Error: ", err)
			}
		}
	},
}

var refreshCmd = &cobra.Command{
	Use:   "refresh",
	Short: "refresh all slash commands from the discord bot.",
	Long: `
This command will refresh all slash commands registered for the bot. 
By default, this will apply globally, but you can specify a guild ID to 
apply to a specific guild. This is required if you change the *commands*,
but not the handlers.
`,
	Run: refreshCommandsFunc,
}

func refreshCommandsFunc(cmd *cobra.Command, args []string) {
	bot, err := discord.NewBot(db, config)
	if err != nil {
		slog.Error("error creating discord bot", "guildId", guildId, "error", err)
		return
	}
	err = bot.RemoveSlashCommands(guildId)
	if err != nil {
		slog.Error("error removing slash (/) commands", "guildId", guildId, "error", err)
		return
	}
	_, err = bot.RegisterSlashCommands(guildId)
	if err != nil {
		slog.Error("error registering slash (/) commands", "error", err)
		return
	}
	slog.Info("successfully registered commands")
}

func init() {
	rootCmd.AddCommand(discordCmd)
	discordCmd.AddCommand(startBot)
	discordCmd.AddCommand(slashCommands)

	slashCommands.AddCommand(refreshCmd)

	refreshCmd.Flags().StringVarP(&guildId, "guild-id", "g", "",
		"Which guild to apply commands for. Defaults to global.")

}
