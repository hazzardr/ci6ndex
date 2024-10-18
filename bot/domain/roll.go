package domain

import (
	"ci6ndex-bot/domain/generated"
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

func (dbo *DatabaseOperations) RollForPlayers(guildId uint64, poolSize int) ([]Offering, error) {
	ctx := context.TODO()
	db, err := dbo.getDB(guildId)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get database connection")
	}
	players, err := dbo.GetPlayersFromActiveDraft(guildId)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get players from active draft")
	}
	leaders, err := db.Queries.GetEligibleLeaders(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get eligible leaders")
	}
	d, err := db.Queries.GetActiveDraft(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get active draft")
	}

	for _, l := range leaders {
		dbo.logger.Debug("Leader", "leader", l)
		if l.ID == 0 {
			dbo.logger.Error("Leader ID is 0", "leader", l)
			return nil, errors.New("leader ID is 0")
		}
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
	for _, p := range pools {
		for _, l := range p.leaders {
			if l.ID == 0 {
				dbo.logger.Error("Leader ID is 0", "leader", l)
				return nil, errors.New("leader ID is 0")
			}
		}
	}

	offers := make([]Offering, 0)
	alreadyPicked := make([]generated.Leader, 0)

	dbo.logger.Info("Rolling for draft", "guildId", guildId)
	for _, p := range pools {
		dbo.logger.Debug("Rolling pool", "pool", p)
		leader, err := tryPickLeader(p.leaders, alreadyPicked)
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
	dbo.logger.Info("Finished rolling for draft", "guildId", guildId, "offers", offers)

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
	dbo.logger.Info("Finished adding pools to database", "guildId", guildId)
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
