package cmd

import (
	"ci6ndex-bot/pkg"
	"context"
	"fmt"
	"github.com/spf13/cobra"
	"log/slog"
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
	Short: "refreshes used by the tool",
	Long: "We use a Google Sheet to manage the rankings. " +
		"This command will refresh the rankings from the Google Sheet and update the database.",
	Run: refreshRanks,
}

// Refreshes the ranks from Google sheets and updates the database
func refreshRanks(cmd *cobra.Command, args []string) {
	ctx := context.Background()
	newRanks, err := pkg.GetRankingsFromSheets(config, ctx)
	if err != nil {
		slog.Error("Error getting newRanks from google sheets", "error", err)
		return
	}
	if newRanks == nil {
		slog.Error("Error getting newRanks from google sheets", "error", "newRanks is nil")
	}
	if len(newRanks) == 0 {
		slog.Error("Error getting newRanks from google sheets", "error", "newRanks is empty")
	}
	slog.Info("Successfully got newRanks from google sheets", "count", len(newRanks))

	err = pkg.UpdateRankings(ctx, newRanks, db)
	if err != nil {
		slog.Error("Error updating pkg rankings", "error", err)
		return
	}
	slog.Info("Successfully updated pkg rankings", "count", len(newRanks))
	err = pkg.ComputeAverageTier(ctx, db)
	if err != nil {
		slog.Error("Error computing average tier", "error", err)
		return
	}
	slog.Info("Successfully computed average tier for each leader")
}

func init() {
	rootCmd.AddCommand(rankingsCmd)
	rankingsCmd.AddCommand(refreshRankingsCmd)
}
