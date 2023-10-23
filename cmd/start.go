package cmd

import (
	"ci6ndex/api"
	"github.com/spf13/cobra"
)

var mode string

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Starts the application",
	Long: `
Starts the application.
`,
	Run: func(cmd *cobra.Command, args []string) {
		api.Start(mode)
	},
}

func init() {
	startCmd.PersistentFlags().StringVarP(&mode, "mode", "m", "bot", "The mode the app should start in.")
	rootCmd.AddCommand(startCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// startCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// startCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
