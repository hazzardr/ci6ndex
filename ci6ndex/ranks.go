package ci6ndex

import (
	"ci6ndex/ci6ndex/generated"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"golang.org/x/sync/errgroup"
	"golang.org/x/sync/semaphore"
)

func (c *Ci6ndex) SubmitRankForPlayer(guildID uint64, rank string, playerID int64, leaderID int64) error {
	db, err := c.getDB(guildID)
	if err != nil {
		return err
	}
	tier, err := GetTierByName(rank)
	if err != nil {
		return err
	}

	err = db.Writes.SubmitRankForPlayer(context.Background(), generated.SubmitRankForPlayerParams{
		PlayerID: playerID,
		LeaderID: leaderID,
		Tier:     tier.Value(),
	})
	if err != nil {
		return fmt.Errorf("failed to submit ranking of %v for playerID %d: %w", tier, playerID, err)
	}
	return nil
}

func (c *Ci6ndex) CalculateTierForLeader(guildID uint64, leaderID int64) error {
	slog.Info("calculating tier", "guild", guildID, "leader", leaderID)
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	db, err := c.getDB(guildID)
	if err != nil {
		return err
	}

	ranks, err := db.Queries.GetAllRanksForLeader(ctx, leaderID)
	if err != nil {
		return err
	}
	if len(ranks) == 0 {
		slog.Info("no rankings found", "guildID", guildID, "leader", leaderID)
		return nil
	}

	sum := 0.0
	for _, r := range ranks {
		sum += r.Tier
	}
	averageRank := sum / float64(len(ranks))
	slog.Info("updating tier", "guild", guildID, "leader", leaderID, "averageRank", averageRank, "allRanks", ranks)
	err = db.Writes.UpdateLeaderTier(ctx, generated.UpdateLeaderTierParams{
		Tier: averageRank,
		ID:   leaderID,
	})
	return err
}

// CalculateTiers will compute average tiers based on user rankings and update leaders table with the result
func (c *Ci6ndex) CalculateTiers(guildID uint64) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	db, err := c.getDB(guildID)
	if err != nil {
		return err
	}

	ranks, err := db.Queries.GetAllRanks(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil
		}
		return err
	}

	// max db conns
	var sem = semaphore.NewWeighted(20)
	eg, egCtx := errgroup.WithContext(ctx)

	ranksByLeads := make(map[int64][]generated.Rank)
	for _, r := range ranks {
		ranksByLeads[r.LeaderID] = append(ranksByLeads[r.LeaderID], r)
	}

	for leaderID, rs := range ranksByLeads {
		sum := 0.0
		for _, r := range rs {
			sum += r.Tier
		}
		averageRank := sum / float64(len(rs))

		eg.Go(func() error {
			err := sem.Acquire(egCtx, 1)
			if err != nil {
				return err
			}
			defer sem.Release(1)
			err = db.Writes.UpdateLeaderTier(egCtx, generated.UpdateLeaderTierParams{
				Tier: averageRank,
				ID:   leaderID,
			})
			return err
		})
	}

	err = eg.Wait()
	return err
}
