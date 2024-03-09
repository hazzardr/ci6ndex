package internal

import (
	"ci6ndex/domain"
	"encoding/json"
	"github.com/pkg/errors"
	"log/slog"
)

var BannedLeaders = []string{
	"ABE",
	"TOMYRIS",
	"GILGAMESH",
	"HAMMURABI",
}

type CivShuffler struct {
	Leaders       []domain.Ci6ndexLeader
	Players       []string
	DraftStrategy domain.Ci6ndexDraftStrategy
	Functions     map[string]shuffleFunction
	DB            *DatabaseOperations
}

// shuffleFunction is a function that takes in a slice of leaders to be assigned,
// and a string representing the user it's being assigned to.
// It returns true if the proposed pool is valid for the user.
type shuffleFunction func([]domain.Ci6ndexLeader, string, *DatabaseOperations) bool

func NewCivShuffler(leaders []domain.Ci6ndexLeader, players []string, strategy domain.Ci6ndexDraftStrategy) CivShuffler {
	return CivShuffler{
		Leaders:       leaders,
		Players:       players,
		DraftStrategy: strategy,
		Functions: map[string]shuffleFunction{
			"allPick":   allPick,
			"allRandom": randomPick,
		},
	}
}

func (c *CivShuffler) Shuffle() ([]DraftOffering, error) {
	allRolls := make([]DraftOffering, 0)
	slog.Info("rolling civs for players", "players", c.Players, "strategy", c.DraftStrategy)

	var rules map[string]interface{}

	err := json.Unmarshal(c.DraftStrategy.Rules, &rules)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal draft strategy rules")
	}

	return allRolls, nil
}

func allPick(leaders []domain.Ci6ndexLeader, user string, db *DatabaseOperations) bool {
	return true
}

func randomPick(leaders []domain.Ci6ndexLeader, user string, db *DatabaseOperations) bool {
	return areElementsUnique(leaders)
}

func areElementsUnique(slice []domain.Ci6ndexLeader) bool {
	elements := make(map[domain.Ci6ndexLeader]bool)
	for _, item := range slice {
		if _, exists := elements[item]; exists {
			// Element already encountered, so elements are not unique
			return false
		}
		// Add element to map
		elements[item] = true
	}
	// All elements are unique
	return true
}
