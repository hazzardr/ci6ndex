package internal

import (
	"ci6ndex/domain"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"regexp"
	"strings"
)

type Ranking struct {
	Player            string
	CombinedLeaderCiv string
	Tier              float64
}

func (r Ranking) ToRankingDBParam(ctx context.Context, db *DatabaseOperations) (domain.CreateRankingsParams, error) {
	user, err := db.Queries.GetUserByName(ctx, r.Player)
	if err != nil {
		return domain.CreateRankingsParams{}, errors.New(fmt.Sprintf("could not find user=%v from google sheets in local database", r.Player))
	}

	re, err := regexp.Compile(`^(.*?) \((.*?)\)$`)
	if err != nil {
		return domain.CreateRankingsParams{}, err
	}
	matches := re.FindStringSubmatch(r.CombinedLeaderCiv)

	var civ string
	var leader string
	if len(matches) == 3 {
		civ = matches[1]
		leader = matches[2]
	} else {
		return domain.CreateRankingsParams{}, errors.New("could not parse civ and leader from google sheets cell")
	}

	l, err := db.Queries.GetLeaderByNameAndCiv(ctx, domain.GetLeaderByNameAndCivParams{
		LeaderName: strings.ToUpper(leader),
		CivName:    strings.ToUpper(civ),
	})

	if err != nil {
		return domain.CreateRankingsParams{}, err
	}

	return domain.CreateRankingsParams{
		UserID:   user.ID,
		Tier:     r.Tier,
		LeaderID: l.ID,
	}, nil
}

// UpdateRankings updates the rankings in the database
func UpdateRankings(ctx context.Context, newRanks []Ranking, db *DatabaseOperations) error {
	// Get old ranks in case replacement fails
	oldRankings, err := db.Queries.GetRankings(ctx)
	if err != nil {
		slog.Error("Error getting old newRanks from db", "error", err)
		return err
	}

	err = db.Queries.DeleteRankings(ctx)
	if err != nil {
		slog.Error("Error deleting old newRanks from db", "error", err)
		return err
	}

	dbRanks := make([]domain.CreateRankingsParams, 0)

	for _, rank := range newRanks {
		p, err := rank.ToRankingDBParam(ctx, db)
		if err != nil {
			slog.Error("error converting newRanks to db params. will try to reinsert old rankings.",
				"rank", rank, "error", err)
			insertErr := insertRanks(ctx, oldRankings, db)
			if insertErr != nil {
				slog.Error("error reinserting old rankings. database is empty", "error", err)
				return insertErr
			}
			slog.Info("successfully reinserted old rankings", "count", len(oldRankings))
			return err
		}
		dbRanks = append(dbRanks, p)
	}
	slog.Info("converted rankings", "count", len(dbRanks))
	_, err = db.Queries.CreateRankings(ctx, dbRanks)
	if err != nil {
		slog.Error("error inserting newRanks to db. will try to reinsert old rankings.", "error", err)
		err = insertRanks(ctx, oldRankings, db)
		if err != nil {
			slog.Error("error reinserting old rankings. database is empty", "error", err)
			return err
		}
		slog.Info("successfully reinserted old rankings", "count", len(oldRankings))
		return err
	}
	return nil
}

// insertRanks is a helper function to reinsert old rankings in case of failure
func insertRanks(ctx context.Context, ranks []domain.Ci6ndexRanking, db *DatabaseOperations) error {
	toInsert := make([]domain.CreateRankingsParams, len(ranks))
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

type Result struct {
	LeaderName string
	Tier       float64
	Err        error
}

func ComputeAverageTier(ctx context.Context, db *DatabaseOperations) error {
	leaders, err := db.Queries.GetLeaders(ctx)
	if err != nil {
		return err
	}
	results := make(chan Result)
	for _, l := range leaders {
		go func(l domain.Ci6ndexLeader) {
			var res Result
			res.LeaderName = l.LeaderName
			slog.Info("computing average tier for", "leader", l.LeaderName)
			ranks, err := db.Queries.GetRankingsByLeader(ctx, l.ID)
			if err != nil {
				res.Err = err
				results <- res
				return
			}
			var average float64
			for _, r := range ranks {
				average += r.Tier
			}
			average = average / float64(len(ranks))
			_, err = db.Queries.UpdateLeaderTier(ctx, domain.UpdateLeaderTierParams{
				ID:   l.ID,
				Tier: average,
			})
			if err != nil {
				res.Err = err
				results <- res
				return
			}
			res.Tier = average
			results <- res
		}(l)
	}

	errs := errors.New("failed to compute average tier for some leaders")
	totalCount := len(leaders)
	for i := 1; i <= totalCount; i++ {
		res := <-results
		if res.Err != nil {
			errs = errors.Join(errs, res.Err)
			continue
		}
		slog.Info("computed average tier", "number", i, "total", totalCount)
	}
	return nil
}
