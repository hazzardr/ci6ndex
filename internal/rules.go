package internal

import (
	"ci6ndex/domain"
	"encoding/json"
	"github.com/pkg/errors"
	"log/slog"
)

var PermaBannedLeaders = []string{
	"ABE",
	"TOMYRIS",
	"GILGAMESH",
	"HAMMURABI",
}

type CivShuffler struct {
	Leaders       []domain.Ci6ndexLeader
	Players       []string
	DraftStrategy domain.Ci6ndexDraftStrategy
	Functions     map[string]*shuffleFunction
	DB            *DatabaseOperations
}

type shuffleFunction struct {
	shuffle  shuffleFunc
	validate validationFunction
}

// validationFunction is a function that takes in a slice of leaders to be assigned,
// and a string representing the user it's being assigned to.
// It returns true if the proposed pool is valid for the user.
type validationFunction func([]domain.Ci6ndexLeader, string, *DatabaseOperations) bool

// validationFunction is a function that takes in a slice of leaders to be assigned,
// and a string representing the user it's being assigned to.
// It returns true if the proposed pool is valid for the user.
type shuffleFunc func([]domain.Ci6ndexLeader, string,
	*DatabaseOperations) ([]domain.Ci6ndexLeader, error)

func NewCivShuffler(leaders []domain.Ci6ndexLeader, players []string, strategy domain.Ci6ndexDraftStrategy) CivShuffler {
	return CivShuffler{
		Leaders:       leaders,
		Players:       players,
		DraftStrategy: strategy,
		Functions: map[string]*shuffleFunction{
			"AllPick": {
				shuffle:  allPick,
				validate: allPickValidate,
			},
		},
	}
}

func (c *CivShuffler) Shuffle() ([]DraftOffering, error) {
	slog.Info("rolling civs for players", "players", c.Players, "strategy", c.DraftStrategy)
	slog.Info("banned leaders", "permaBanned", PermaBannedLeaders)

	eligibleLeaders := make([]domain.Ci6ndexLeader, 0)
	for _, banned := range PermaBannedLeaders {
		for _, leader := range c.Leaders {
			if leader.LeaderName != banned {
				eligibleLeaders = append(eligibleLeaders, leader)
			}
		}
	}

	allRolls := make([]DraftOffering, len(c.Players))

	// All pick
	// No rules in all pick strategies for now
	if c.DraftStrategy.Randomize == false {
		for _, player := range c.Players {
			roll := DraftOffering{
				User:    player,
				Leaders: eligibleLeaders,
			}
			allRolls = append(allRolls, roll)
		}
		return allRolls, nil
	}

	// Random pick
	var rules map[string]interface{}

	if hasRules(c.DraftStrategy) {
		err := json.Unmarshal(c.DraftStrategy.Rules, &rules)
		if err != nil {
			return nil, errors.Wrap(err, "failed to unmarshal draft strategy rules")
		}
	}

	return allRolls, nil
}

func allPickValidate(leaders []domain.Ci6ndexLeader, user string, db *DatabaseOperations) bool {
	return true
}

func allPick(leaders []domain.Ci6ndexLeader, user string, db *DatabaseOperations) ([]domain.Ci6ndexLeader, error) {
	return leaders, nil
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

func hasRules(strategy domain.Ci6ndexDraftStrategy) bool {
	return strategy.Rules != nil && len(strategy.Rules) > 0
}
