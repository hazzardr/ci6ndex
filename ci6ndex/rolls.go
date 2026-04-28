package ci6ndex

import (
	"ci6ndex/ci6ndex/generated"
	"context"
	"fmt"
	"math/rand/v2"
)

// Offering represents a set of leaders offered to a player in a draft.
type Offering struct {
	Player  generated.Player
	Leaders []generated.Leader
	DraftId int64
}

// RollForPlayers rolls leaders for a set of players based on the provided rules.
func (c *Ci6ndex) RollForPlayers(
	guildId uint64,
	playerIds []int64,
	rules []Rule,
) ([]Offering, error) {
	ctx := context.TODO()
	db, err := c.getDB(guildId)
	if err != nil {
		return nil, fmt.Errorf("failed to get database for guild %d: %w", guildId, err)
	}

	players, err := db.Queries.GetPlayersFromActiveDraft(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get players: %w", err)
	}

	allLeaders, err := db.Queries.GetEligibleLeaders(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get leaders: %w", err)
	}

	draft, err := db.Queries.GetActiveDraft(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get active draft: %w", err)
	}

	playerMap := make(map[int64]generated.Player, len(players))
	for _, p := range players {
		playerMap[p.ID] = p
	}

	assigned := make(map[int64]struct{})
	poolSize := len(rules)
	offerings := make([]Offering, 0, len(playerIds))

	for _, playerId := range playerIds {
		player, ok := playerMap[playerId]
		if !ok {
			continue
		}

		// Separate rules by type.
		var allRules, atLeastOneRules []Rule
		for _, rule := range rules {
			if rule.Type() == All {
				allRules = append(allRules, rule)
			} else {
				atLeastOneRules = append(atLeastOneRules, rule)
			}
		}

		// Start with leaders that satisfy all "All" rules.
		valid := make([]generated.Leader, len(allLeaders))
		copy(valid, allLeaders)
		for _, rule := range allRules {
			valid = rule.Filter(player, valid)
		}

		// Remove globally assigned leaders.
		valid = filterAssigned(valid, assigned)

		if len(valid) < poolSize && len(atLeastOneRules) == 0 {
			return nil, RanOutOfChoicesError{}
		}

		var selected []generated.Leader

		if len(atLeastOneRules) > 0 && len(valid) > 0 {
			// Satisfy each "AtLeastOne" rule with a distinct leader.
			for _, rule := range atLeastOneRules {
				candidates := rule.Filter(player, valid)
				if len(candidates) == 0 {
					continue
				}
				picked := candidates[rand.IntN(len(candidates))]
				selected = append(selected, picked)
				valid = removeLeader(valid, picked.ID)
			}

			if len(selected) < len(atLeastOneRules) {
				return nil, RanOutOfChoicesError{}
			}

			remaining := poolSize - len(selected)
			if remaining > 0 && len(valid) > 0 {
				if remaining > len(valid) {
					remaining = len(valid)
				}
				selected = append(selected, pickN(valid, remaining)...)
			}
		} else {
			if len(valid) < poolSize {
				return nil, RanOutOfChoicesError{}
			}
			selected = pickN(valid, poolSize)
		}

		for _, l := range selected {
			assigned[l.ID] = struct{}{}
		}

		offerings = append(offerings, Offering{
			Player:  player,
			Leaders: selected,
			DraftId: draft.ID,
		})
	}

	return offerings, nil
}

// filterAssigned returns leaders that have not been assigned yet.
func filterAssigned(leaders []generated.Leader, assigned map[int64]struct{}) []generated.Leader {
	filtered := leaders[:0]
	for _, l := range leaders {
		if _, ok := assigned[l.ID]; !ok {
			filtered = append(filtered, l)
		}
	}
	return filtered
}

// removeLeader returns a new slice with the leader matching id removed.
func removeLeader(leaders []generated.Leader, id int64) []generated.Leader {
	filtered := leaders[:0]
	for _, l := range leaders {
		if l.ID != id {
			filtered = append(filtered, l)
		}
	}
	return filtered
}

// pickN returns n randomly chosen distinct leaders from the slice.
func pickN(leaders []generated.Leader, n int) []generated.Leader {
	if n >= len(leaders) {
		result := make([]generated.Leader, len(leaders))
		copy(result, leaders)
		return result
	}
	shuffled := make([]generated.Leader, len(leaders))
	copy(shuffled, leaders)
	rand.Shuffle(len(shuffled), func(i, j int) {
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	})
	return shuffled[:n]
}

type RanOutOfChoicesError struct{}

func (e RanOutOfChoicesError) Error() string {
	return "no leaders left to pick from"
}
