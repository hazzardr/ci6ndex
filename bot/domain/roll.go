package domain

import (
	"ci6ndex-bot/domain/generated"
	"context"
	"github.com/pkg/errors"
	"math/rand/v2"
	"slices"
)

// OfferPool is a struct that holds a pool of leaders that are valid for a given player
// The leaders are mapped from a rule to a list of leaders that satisfy the rule
type OfferPool struct {
	leaders []generated.Leader
	rule    Rule
	player  generated.Player
}

// NewPool creates a pool of leaders to be offered to a player after filtering based on the rule
func NewPool(player generated.Player, rule Rule, leaders []generated.Leader) *OfferPool {
	filteredLeaders := rule.evaluate(player, leaders)
	return &OfferPool{
		leaders: filteredLeaders,
		rule:    rule,
		player:  player,
	}
}

func (dbo *DatabaseOperations) RollForPlayers(guildId uint64, poolSize int) error {
	ctx := context.TODO()
	db, err := dbo.getDB(guildId)
	if err != nil {
		return errors.Wrap(err, "failed to get database connection")
	}
	players, err := dbo.GetPlayersFromActiveDraft(guildId)
	if err != nil {
		return errors.Wrap(err, "failed to get players from active draft")
	}
	leaders, err := db.Queries.GetEligibleLeaders(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to get eligible leaders")
	}

	// Give everyone a tier 3 or above leader
	// Try to do harder to assign rules first
	pools := make([]OfferPool, poolSize)
	for _, player := range players {
		pool := NewPool(player, &MinTierRule{minTier: 3.0}, leaders)
		pools = append(pools, *pool)
	}
	for i := 1; i < poolSize; i++ {
		for _, player := range players {
			// Get the pool for the player
			pool := NewPool(player, &NoOpRule{}, leaders)
			pools = append(pools, *pool)
		}
	}

	offers := make(map[int64][]generated.Leader)
	alreadyPicked := make([]generated.Leader, 0)

	dbo.logger.Info("Rolling for draft", "guildId", guildId)
	for _, p := range pools {
		dbo.logger.Debug("Rolling pool", "pool", p)
		leader, err := tryPickLeader(p.leaders, alreadyPicked)
		if err != nil {
			if errors.As(err, &RanOutOfChoicesError{}) {
				return errors.Wrap(err, "ran out of choices - please retry!")
			}
			return errors.Wrapf(err, "failed to pick leader for pool=%v", p)
		}
		offers[p.player.ID] = append(offers[p.player.ID], leader)
		alreadyPicked = append(alreadyPicked, leader)
	}
	dbo.logger.Info("Finished rolling for draft", "guildId", guildId, "offers", offers)

	d, err := db.Queries.GetActiveDraft(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to get active draft")
	}

	for playerId, leaders := range offers {
		for _, leader := range leaders {
			err := db.Writes.AddPool(ctx, generated.AddPoolParams{
				PlayerID: playerId,
				DraftID:  d.ID,
				Leader:   leader.ID,
			})
			if err != nil {
				return errors.Wrap(err, "failed to add pool to database")
			}
		}
	}
	dbo.logger.Info("Finished adding pools to database", "guildId", guildId)
	return nil
}

type RanOutOfChoicesError struct{}

func (e RanOutOfChoicesError) Error() string {
	return "no leaders left to pick from"
}

func tryPickLeader(pickFrom []generated.Leader, alreadyPicked []generated.Leader) (generated.Leader, error) {
	if len(pickFrom) == 0 {
		return generated.Leader{}, RanOutOfChoicesError{}
	}
	i := rand.IntN(len(pickFrom)) - 1
	randPick := pickFrom[i]
	if containsLeader(randPick, alreadyPicked) {
		pickFrom = slices.Delete(pickFrom, i, i+1)
		return tryPickLeader(pickFrom, alreadyPicked)
	}
	return randPick, nil
}

func containsLeader(pick generated.Leader, alreadyPicked []generated.Leader) bool {
	for _, p := range alreadyPicked {
		if p.ID == pick.ID {
			return true
		}
	}
	return false
}
