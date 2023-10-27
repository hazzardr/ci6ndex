package cmd

import (
	"ci6ndex/internal"
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
		internal.Start(mode)
	},
}

func init() {
	startCmd.PersistentFlags().StringVarP(&mode, "mode", "m", "bot", "The mode the app should start in.")
	rootCmd.AddCommand(startCmd)
}
