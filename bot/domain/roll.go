package domain

import (
	"ci6ndex-bot/domain/generated"
	"context"
	"github.com/pkg/errors"
	"math/rand/v2"
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
	for _, p := range pools {
		i := rand.IntN(len(p.leaders)) - 1
		randPick := p.leaders[i]
		if offers[p.player.ID] == nil {
			offers[p.player.ID] = make([]generated.Leader, poolSize)
		}
		offers[p.player.ID] = append(offers[p.player.ID], randPick)
		alreadyPicked = append(alreadyPicked, randPick)

	}
}

func tryPickLeader(p OfferPool, alreadyPicked []generated.Leader) (generated.Leader, error) {
	i := rand.IntN(len(p.leaders)) - 1
	randPick := p.leaders[i]
	if !containsLeader(randPick, alreadyPicked) {
		return randPick, nil
	}
}

func containsLeader(pick generated.Leader, alreadyPicked []generated.Leader) bool {
	for _, p := range alreadyPicked {
		if p.ID == pick.ID {
			return true
		}
	}
	return false
}
