package internal

import (
	"ci6ndex/domain"
	"log/slog"
	"math/rand/v2"
)

type CivShuffler struct {
	Leaders       []domain.Ci6ndexLeader
	Players       []string
	DraftStrategy domain.Ci6ndexDraftStrategy
	Functions     map[string]shuffleFunction
}

// shuffleFunction
type shuffleFunction func([]domain.Ci6ndexLeader, int) []domain.Ci6ndexLeader

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

	return allRolls, nil
}

func allPick(leaders []domain.Ci6ndexLeader, numPicks int) []domain.Ci6ndexLeader {
	return leaders
}

func randomPick(leaders []domain.Ci6ndexLeader, numPicks int) []domain.Ci6ndexLeader {
	offers := make([]domain.Ci6ndexLeader, numPicks)

	for i := 0; i < numPicks; i++ {
		r := rand.N(len(leaders))
		randomLeader := leaders[r]
		if !elementExists(offers, randomLeader) {
			offers = append(offers, randomLeader)
		}
	}
	return offers
}

func elementExists(slice []domain.Ci6ndexLeader, leader domain.Ci6ndexLeader) bool {
	for _, item := range slice {
		if item == leader {
			return true
		}
	}
	return false
}
