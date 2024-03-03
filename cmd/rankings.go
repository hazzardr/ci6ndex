package cmd

import (
	"ci6ndex/domain"
	"ci6ndex/internal"
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

// Refreshes the ranks from google sheets and updates the database
func refreshRanks(cmd *cobra.Command, args []string) {
	ctx := context.Background()
	newRanks, err := internal.GetRankingsFromSheets(config, ctx)
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

	// Get old ranks in case replacement fails
	oldRankings, err := db.Queries.GetRankings(ctx)
	if err != nil {
		slog.Error("Error getting old newRanks from db", "error", err)
		return
	}

	err = db.Queries.DeleteRankings(ctx)
	if err != nil {
		slog.Error("Error deleting old newRanks from db", "error", err)
		return
	}

	dbRanks := make([]domain.CreateRankingsParams, 0, len(newRanks))

	for _, rank := range newRanks {
		p, err := rank.ToRankingDBParam(ctx)
		if err != nil {
			slog.Error("error converting newRanks to db params. will try to reinsert old rankings.",
				"rank", rank, "error", err)
			err = insertRanks(ctx, oldRankings)
			if err != nil {
				slog.Error("error reinserting old rankings. database is empty", "error", err)
				return
			}
			slog.Info("successfully reinserted old rankings", "count", len(oldRankings))
			return
		}
		dbRanks = append(dbRanks, p)
	}

	_, err = db.Queries.CreateRankings(ctx, dbRanks)
	if err != nil {
		slog.Error("error inserting newRanks to db. will try to reinsert old rankings.",
			"newRanks", dbRanks, "error", err)
		err = insertRanks(ctx, oldRankings)
		if err != nil {
			slog.Error("error reinserting old rankings. database is empty", "error", err)
			return
		}
		slog.Info("successfully reinserted old rankings", "count", len(oldRankings))
		return
	}
	slog.Info("successfully inserted new ranks from google sheets", "count", len(dbRanks))
}

// insertRanks is a helper function to reinsert old rankings in case of failure
func insertRanks(ctx context.Context, ranks []domain.Ci6ndexRanking) error {
	toInsert := make([]domain.CreateRankingsParams, 0, len(ranks))
	for _, rank := range ranks {
		p := domain.CreateRankingsParams{
			UserID:   rank.UserID,
			Tier:     rank.Tier,
			LeaderID: rank.LeaderID,
		}
		toInsert = append(toInsert, p)
	}
	_, err := db.Queries.CreateRankings(ctx, toInsert)
	if err != nil {
		return err
	}
	slog.Info("successfully reinserted old rankings", "count", len(toInsert))
	return nil
}

func init() {
	rootCmd.AddCommand(rankingsCmd)
	rankingsCmd.AddCommand(refreshRankingsCmd)
}
