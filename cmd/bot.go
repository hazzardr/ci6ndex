package cmd

import (
	"github.com/spf13/cobra"
)

// botCmd represents the bot command
var botCmd = &cobra.Command{
	Use:   "bot",
	Short: "Start the discord bot",
	Long: `
Starts the discord bot, which will listen for commands and respond to them.
`,
	Run: func(cmd *cobra.Command, args []string) {
	},
}

// startCmd represents the bot command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the discord bot",
	Long: `
Starts the discord bot, which will listen for commands and respond to them.
`,
	Run: func(cmd *cobra.Command, args []string) {
	},
}

func startBot(cmd *cobra.Command, args []string) {

	//bot := &internal.DiscordBot{
	//	db: db,
	//	config: config,
	//}
}

func init() {
	rootCmd.AddCommand(botCmd)
	botCmd.AddCommand(startCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// botCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// botCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
