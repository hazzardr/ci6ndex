package cmd

import (
	"ci6ndex/internal"
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"log/slog"
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
		var rules interface{}
		if strategy.Rules == nil {
			rules = "none"
		} else {
			err := json.Unmarshal(strategy.Rules, &rules)
			if err != nil {
				slog.Error("Error decoding rules", "error", err)
				continue
			}
		}

		fmt.Printf("Name: %s, Description: %s, PoolSize: %d, Randomize: %t, Rules: %v\n",
			strategy.Name, strategy.Description, strategy.PoolSize, strategy.Randomize, rules)
	}
}

var (
	rollCivsCmd = &cobra.Command{
		Use:   "roll",
		Short: "rolls civs, optionally for a draft",
		Long: `Rolls civs based on the input draft strategy, draft, players,
and leaders in the database. If you attempt to roll for a draft,
all other parameters will be ignored. If you do not specify a draft, 
you must specify both a draft strategy as well as players to roll for. 
Leaders will be selected from the database, and cannot be provided by the CLI.`,
		Run: rollCivs,
	}
	draft   bool
	players []string
)

func rollCivs(cmd *cobra.Command, args []string) {
	ctx := cmd.Context()
	anyFlags := false
	cmd.Flags().Visit(func(f *pflag.Flag) {
		if f.Changed {
			anyFlags = true
		}
	})

	if !anyFlags {
		err := cmd.Help()
		if err != nil {
			fmt.Println("Error: ", err)
		}
		return
	}

	ls, err := db.Queries.GetLeaders(ctx)
	if err != nil {
		slog.Error("Error fetching leaders", "error", err)
		return
	}
	var shuffler *internal.CivShuffler
	if draft {
		drafts, err := db.Queries.GetActiveDrafts(ctx)
		if err != nil {
			slog.Error("Error fetching active drafts", "error", err)
			return
		}
		if nil == drafts || len(drafts) == 0 {
			slog.Error("Asked to roll for draft, but no active draft found")
			return
		}
		d := drafts[0]
		if len(drafts) > 1 {
			slog.Error("More than one active draft found", "drafts", d)
			return
		}
		ds, err := db.Queries.GetDraftStrategy(ctx, d.DraftStrategy)
		if err != nil {
			slog.Error("Error fetching draft strategy", "draftId", d.ID, "error", err)
			return
		}
		players, err := db.Queries.GetPlayersForDraft(ctx, d.ID)
		if err != nil {
			slog.Error("Error fetching players for draft", "draftId", d.ID, "error", err)
			return
		}
		pNames := make([]string, len(players))
		for i, p := range players {
			pNames[i] = p.DiscordName
		}
		shuffler = internal.NewCivShuffler(ls, pNames, ds, db)
	} else {
		if draftStrategy == "" {
			slog.Error("No draft strategy provided")
			return
		}
		if len(players) == 0 {
			slog.Error("No players provided")
			return
		}
		ds, err := db.Queries.GetDraftStrategy(ctx, draftStrategy)
		if err != nil {
			slog.Error("Error fetching draft strategy", "strategy", draftStrategy, "error", err)
			return
		}
		shuffler = internal.NewCivShuffler(ls, players, ds, db)
	}
	offers, err := shuffler.Shuffle()
	if err != nil {
		slog.Error("Error shuffling civs", "error", err)
		return
	}

	finalRolls := map[string][]string{}
	for _, offer := range offers {
		user := offer.User
		leaders := make([]string, len(offer.Leaders))
		for i, l := range offer.Leaders {
			leaders[i] = l.LeaderName
		}
		finalRolls[user] = leaders
	}
	fmt.Print("Offers are as follows:\n")
	for user := range finalRolls {
		fmt.Printf("user: %s leaders: %s\n", user, finalRolls[user])
	}

}

func init() {
	rootCmd.AddCommand(draftsCmd)
	rootCmd.AddCommand(strategiesCommand)
	rootCmd.AddCommand(rollCivsCmd)
	draftsCmd.AddCommand(startDraftCommand)
	draftsCmd.AddCommand(getDraftsCommand)

	startDraftCommand.Flags().StringVarP(&draftStrategy, "draft-strategy", "s", "",
		"The strategy to use for the draft")
	getDraftsCommand.Flags().Bool("active", false, "Get the active draft")
	rollCivsCmd.Flags().BoolVarP(&draft, "draft", "d", false, "Roll for the active draft")
	rollCivsCmd.Flags().StringVarP(&draftStrategy, "draft-strategy", "s", "RandomPick",
		"Roll using a specific draft strategy")
	rollCivsCmd.Flags().StringSliceVarP(&players, "players", "p", nil, "The players to roll for")
}
