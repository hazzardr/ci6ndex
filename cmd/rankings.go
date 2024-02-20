package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var rankingsCmd = &cobra.Command{
	Use:   "rankings",
	Short: "ci6ndex rankings tools",
	Long: `command used to manage rankings. Rankings are managed through a Google Sheet.
`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			err := cmd.Help()
			if err != nil {
				fmt.Println("Error: ", err)
			}
		}
	},
}

var refreshRankingsCmd = &cobra.Command{
	Use:   "refresh",
	Short: "refresh rankings",
	Long: `command used to refresh rankings by calling out to the google sheet defined in config.
`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("rankings called")
	},
}

func init() {
	rootCmd.AddCommand(rankingsCmd)
	rankingsCmd.AddCommand(refreshRankingsCmd)
}
