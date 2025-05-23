package ci6ndex

import (
	"ci6ndex/ci6ndex/generated"
	"context"
	"github.com/pkg/errors"
	"math/rand/v2"
	"slices"
)

// EligibleLeaders is a struct that holds a pool of leaders that are valid for a given player
// The leaders are mapped from a rule to a list of leaders that satisfy the rule
type EligibleLeaders struct {
	leaders []generated.Leader
	rule    Rule
	player  generated.Player
}

type Offering struct {
	Player  generated.Player
	Leaders []generated.Leader
	DraftId int64
}

// NewPool creates a pool of leaders to be offered to a player after filtering based on the rule
func NewPool(player generated.Player, rule Rule, leaders []generated.Leader) EligibleLeaders {
	filteredLeaders := rule.evaluate(player, leaders)
	zeroed := 0
	for _, l := range filteredLeaders {
		if l.ID == 0 {
			zeroed++
		}
	}
	return EligibleLeaders{
		leaders: filteredLeaders,
		rule:    rule,
		player:  player,
	}
}

// TODO: rewrite
func (c *Ci6ndex) RollForPlayers(guildId uint64, poolSize int) ([]Offering, error) {
	ctx := context.TODO()
	db, err := c.getDB(guildId)
	players, err := c.GetPlayersFromActiveDraft(guildId)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get players from active draft")
	}
	d, err := db.Queries.GetActiveDraft(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get active draft")
	}

	// Give everyone a tier 3 or above leader
	// Try to do harder to assign rules first
	pools := make([]EligibleLeaders, 0)
	for _, player := range players {
		pool := NewPool(player, &MinTierRule{minTier: 3.0}, leaders)
		pools = append(pools, pool)
	}
	for i := 1; i < poolSize; i++ {
		for _, player := range players {
			// Get the pool for the player
			pool := NewPool(player, &NoOpRule{}, leaders)
			pools = append(pools, pool)
		}
	}

	offers := make([]Offering, 0)
	alreadyPicked := make([]generated.Leader, 0)

	c.Logger.Info("Rolling for draft", "guildId", guildId)
	for _, p := range pools {
		c.Logger.Debug("Rolling pool", "pool", p)
		offered := p.leaders
		leader, err := tryPickLeader(offered, alreadyPicked)
		if err != nil {
			if errors.As(err, &RanOutOfChoicesError{}) {
				return nil, errors.Wrap(err, "ran out of choices - please retry!")
			}
			return nil, errors.Wrapf(err, "failed to pick leader for pool=%v", p)
		}
		offers = append(offers, Offering{
			Player:  p.player,
			Leaders: []generated.Leader{leader},
			DraftId: d.ID,
		})
		alreadyPicked = append(alreadyPicked, leader)
	}

	offers = flattenOffers(offers, d.ID)
	c.Logger.Info("Finished rolling for draft", "guildId", guildId, "offers", offers)

	// Wipe existing pools
	err = db.Writes.DeletePoolsForDraftId(ctx, d.ID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to delete existing pools. "+
			"Roll succeeded but metadata failed to update")
	}

	for _, o := range offers {
		for _, leader := range o.Leaders {
			err := db.Writes.AddPool(ctx, generated.AddPoolParams{
				PlayerID: o.Player.ID,
				DraftID:  d.ID,
				Leader:   leader.ID,
			})
			if err != nil {
				return nil, errors.Wrap(err, "failed to add pool to database")
			}
		}
	}
	c.Logger.Info("Finished adding pools to database", "guildId", guildId)
	return offers, nil
}

func flattenOffers(offers []Offering, draftId int64) []Offering {
	offerMap := make(map[int64][]generated.Leader)
	userData := make(map[int64]generated.Player)
	for _, o := range offers {
		offerMap[o.Player.ID] = append(offerMap[o.Player.ID], o.Leaders...)
		userData[o.Player.ID] = o.Player
	}
	flattenedOffers := make([]Offering, 0)
	for playerId, leaders := range offerMap {
		flattenedOffers = append(flattenedOffers, Offering{
			Player:  userData[playerId],
			Leaders: leaders,
			DraftId: draftId,
		})
	}
	return flattenedOffers
}

type RanOutOfChoicesError struct{}

func (e RanOutOfChoicesError) Error() string {
	return "no leaders left to pick from"
}

func tryPickLeader(pickFrom []generated.Leader, alreadyPicked []generated.Leader) (generated.Leader, error) {
	if len(pickFrom) == 0 {
		return generated.Leader{}, RanOutOfChoicesError{}
	}
	i := rand.IntN(len(pickFrom))
	randPick := pickFrom[i]
	// this is a hack, no idea why we get a 0 ID leader sometimes
	if randPick.ID == 0 || containsLeader(randPick, alreadyPicked) {
		rest := slices.Delete(pickFrom, i, i+1)
		return tryPickLeader(rest, alreadyPicked)
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
