package cmd

import (
	"ci6ndex/internal"
	"fmt"
	"log/slog"

	"github.com/spf13/cobra"
)

// draftsCmd represents the drafts command
var draftsCmd = &cobra.Command{
	Use:   "drafts",
	Short: "command to interact with drafts",
	Long: `Allows for general CRUD operations, as well as things like:
1. Starting drafts
2. Offering picks
3. Submitting picks
4. Ending drafts
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

var (
	draftStrategy string
)

var startDraftCommand = &cobra.Command{
	Use:   "start",
	Short: "starts a new draft",
	Long: "Starts a new draft, so long as there is not one in progress." +
		" This will create a new draft in the database and return the draft ID.",
	Run: startDraft,
}

func startDraft(cmd *cobra.Command, args []string) {
	ctx := cmd.Context()
	cdr := internal.CreateDraftRequest{DraftStrategy: draftStrategy}
	draft, err := internal.CreateDraft(ctx, cdr, db)
	if err != nil {
		slog.Error("Error creating draft", "error", err)
		return
	}
	slog.Info("draft created successfully",
		"draft_id", draft.ID, "draft_strategy", draft.DraftStrategy)
}

var getDraftsCommand = &cobra.Command{
	Use:   "get",
	Short: "command to get drafts",
	Long: `Allows for getting drafts from the database.
You can either get all drafts, or you can get the active draft via the --active flag.
`,
	Run: getDrafts,
}

func getDrafts(cmd *cobra.Command, args []string) {
	ctx := cmd.Context()
	active, err := cmd.Flags().GetBool("active")
	if err != nil {
		slog.Error("unable to parse active flag", "error", err)
		return
	}
	drafts, err := internal.GetDrafts(ctx, db, active)
	if err != nil {
		slog.Error("Error fetching drafts", "error", err)
		return
	}
	if drafts == nil {
		fmt.Println("Nil drafts returned")
		return
	}

	for _, draft := range drafts {
		fmt.Println(draft)
	}
}

var strategiesCommand = &cobra.Command{
	Use:   "draft-strategies",
	Short: "list all draft strategies",
	Long:  "Lists all draft strategies available in the database.",
	Run:   listDraftStrategies,
}

func listDraftStrategies(cmd *cobra.Command, args []string) {
	ctx := cmd.Context()
	strategies, err := internal.GetDraftStrategies(ctx, db)
	if err != nil {
		slog.Error("Error fetching draft strategies", "error", err)
		return
	}
	for _, strategy := range strategies {
		fmt.Println(strategy)
	}
}

func init() {
	rootCmd.AddCommand(draftsCmd)
	rootCmd.AddCommand(strategiesCommand)
	draftsCmd.AddCommand(startDraftCommand)
	draftsCmd.AddCommand(getDraftsCommand)

	startDraftCommand.Flags().StringVarP(&draftStrategy, "draft-strategy", "s", "",
		"The strategy to use for the draft")
	getDraftsCommand.Flags().Bool("active", false, "Get the active draft")
}
