package ci6ndex

import (
	"ci6ndex/ci6ndex/generated"
	"context"
	"github.com/pkg/errors"
	"math/rand/v2"
	"slices"
)

type Pool struct {
	Player  generated.Player
	Leaders []generated.Leader
	Rule    Rule
}

func (p *Pool) Evaluate() []generated.Leader {
	return p.Rule.Filter(p.Player, p.Leaders)
}

func NewPool(player generated.Player, leaders []generated.Leader, rule Rule) Pool {
	return Pool{
		Leaders: leaders,
		Player:  player,
		Rule:    rule,
	}
}

// Offering represents a set of leaders offered to a player in a draft.
type Offering struct {
	Player  generated.Player
	Leaders []generated.Leader
	DraftId int64
}

// getPools retrieves all players and leaders from the database,
// and creates a pool for each player / rule combo
func getPools(
	ctx context.Context,
	db *DB,
	playerIds []int64,
	rules []Rule,
) ([]Pool, error) {
	allPlayers, err := db.Queries.GetPlayers(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get players from active draft")
	}

	leaders, err := db.Queries.GetLeaders(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get leaders")
	}
	pools := make([]Pool, 0)
	for _, player := range allPlayers {
		if slices.Contains(playerIds, player.ID) {
			for _, rule := range rules {
				p := NewPool(player, leaders, rule)
				p.Evaluate()
				pools = append(pools, p)
			}
		}
	}
	return pools, nil
}
func (c *Ci6ndex) RollForPlayers(
	guildId uint64,
	playerIds []int64,
	rules []Rule,
) ([]Offering, error) {
	ctx := context.TODO()
	db, err := c.getDB(guildId)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get database for guild %d", guildId)
	}

	pools, err := getPools(ctx, db, playerIds, rules)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get pools for players %v", playerIds)
	}

	// Group pools by player
	playerPools := make(map[int64][]Pool)
	for _, pool := range pools {
		playerPools[pool.Player.ID] = append(playerPools[pool.Player.ID], pool)
	}

	offerings := make([]Offering, 0, len(playerIds))
	draft, err := db.Queries.GetActiveDraft(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get active draft")
	}

	assignedLeaders := make(map[int64]bool)
	poolSize := len(rules)
	for _, playerId := range playerIds {
		playerPoolList, ok := playerPools[playerId]
		if !ok {
			continue // Skip if player has no pools
		}

		player := playerPoolList[0].Player

		// Separate pools by rule type
		var allRulePools []Pool
		var atLeastOnePools []Pool
		for _, pool := range playerPoolList {
			if pool.Rule.Type() == All {
				allRulePools = append(allRulePools, pool)
			} else {
				atLeastOnePools = append(atLeastOnePools, pool)
			}
		}

		// Get leaders that satisfy all "All" type rules
		var validLeaders []generated.Leader
		if len(allRulePools) > 0 {
			validLeaders = allRulePools[0].Evaluate()
			for i := 1; i < len(allRulePools); i++ {
				// Intersect with each subsequent "All" rule pool
				validLeaders = intersectLeaders(validLeaders, allRulePools[i].Evaluate())
			}
		} else if len(playerPoolList) > 0 {
			// If no "All" rules, start with all leaders
			validLeaders = slices.Clone(playerPoolList[0].Leaders)
		}

		// Remove already assigned leaders
		leaderMap := make(map[int64]generated.Leader, len(validLeaders))
		for _, leader := range validLeaders {
			if !assignedLeaders[leader.ID] {
				leaderMap[leader.ID] = leader
			}
		}

		// Convert map back to slice
		validLeaders = make([]generated.Leader, 0, len(leaderMap))
		for _, leader := range leaderMap {
			validLeaders = append(validLeaders, leader)
		}

		if len(validLeaders) < poolSize && len(atLeastOnePools) == 0 {
			return nil, RanOutOfChoicesError{}
		}

		var selectedLeaders []generated.Leader

		if len(atLeastOnePools) > 0 && len(validLeaders) > 0 {
			validLeaderMap := make(map[int64]generated.Leader)
			for _, leader := range validLeaders {
				validLeaderMap[leader.ID] = leader
			}

			atLeastOneSelections := make([]generated.Leader, 0)

			// For each "AtLeastOne" rule, select one valid leader if possible
			for _, pool := range atLeastOnePools {
				possibleLeaders := intersectLeaders(validLeaders, pool.Evaluate())
				if len(possibleLeaders) > 0 {
					// Randomly select one leader that satisfies this rule
					selectedIdx := rand.IntN(len(possibleLeaders))
					selected := possibleLeaders[selectedIdx]
					atLeastOneSelections = append(atLeastOneSelections, selected)

					// Remove from valid leaders map to avoid duplicates
					delete(validLeaderMap, selected.ID)
				}
			}

			// Rebuild validLeaders from the map
			validLeaders = make([]generated.Leader, 0, len(validLeaderMap))
			for _, leader := range validLeaderMap {
				validLeaders = append(validLeaders, leader)
			}

			// If we couldn't satisfy all AtLeastOne rules, return error
			if len(atLeastOneSelections) < len(atLeastOnePools) && len(atLeastOnePools) > 0 {
				return nil, RanOutOfChoicesError{}
			}

			remainingCount := poolSize - len(atLeastOneSelections)
			if remainingCount > 0 && len(validLeaders) > 0 {
				if remainingCount > len(validLeaders) {
					remainingCount = len(validLeaders)
				}

				rand.Shuffle(len(validLeaders), func(i, j int) {
					validLeaders[i], validLeaders[j] = validLeaders[j], validLeaders[i]
				})

				atLeastOneSelections = append(atLeastOneSelections, validLeaders[:remainingCount]...)
			}

			selectedLeaders = atLeastOneSelections
		} else {
			// Without "AtLeastOne" rules, just randomly select poolSize leaders
			if len(validLeaders) > poolSize {
				rand.Shuffle(len(validLeaders), func(i, j int) {
					validLeaders[i], validLeaders[j] = validLeaders[j], validLeaders[i]
				})
				selectedLeaders = validLeaders[:poolSize]
			} else if len(validLeaders) < poolSize {
				// Not enough leaders to satisfy pool size
				return nil, RanOutOfChoicesError{}
			} else {
				selectedLeaders = validLeaders
			}
		}

		for _, leader := range selectedLeaders {
			assignedLeaders[leader.ID] = true
		}

		offering := Offering{
			Player:  player,
			Leaders: selectedLeaders,
			DraftId: draft.ID,
		}
		offerings = append(offerings, offering)
	}

	return offerings, nil
}

func intersectLeaders(a, b []generated.Leader) []generated.Leader {
	result := make([]generated.Leader, 0)
	bMap := make(map[int64]generated.Leader)

	for _, leaderB := range b {
		bMap[leaderB.ID] = leaderB
	}

	for _, leaderA := range a {
		if _, exists := bMap[leaderA.ID]; exists {
			result = append(result, leaderA)
		}
	}

	return result
}

type RanOutOfChoicesError struct{}

func (e RanOutOfChoicesError) Error() string {
	return "no leaders left to pick from"
}
