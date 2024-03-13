package internal

import (
	"ci6ndex/domain"
	"context"
	"github.com/pkg/errors"
	"log/slog"
	"math/rand/v2"
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

// shuffleFunc is a function that takes in a slice of leaders to be assigned,
// and a string representing the user it's being assigned to.
// It returns the output offering based on rules defined in the function.
type shuffleFunc func([]domain.Ci6ndexLeader, string, int,
	*DatabaseOperations) ([]domain.Ci6ndexLeader, error)

func NewCivShuffler(leaders []domain.Ci6ndexLeader, players []string, strategy domain.Ci6ndexDraftStrategy, db *DatabaseOperations) CivShuffler {
	return CivShuffler{
		Leaders:       leaders,
		Players:       players,
		DraftStrategy: strategy,
		DB:            db,
		Functions: map[string]*shuffleFunction{
			"AllPick": {
				shuffle:  allPick,
				validate: allPickValidate,
			},
			"RandomPick": {
				shuffle:  randomPick,
				validate: randomPickValidate,
			},
		},
	}
}

func (c *CivShuffler) Shuffle() ([]DraftOffering, error) {
	slog.Info("rolling civs for players", "players", c.Players, "strategy", c.DraftStrategy)
	slog.Info("banned leaders", "permaBanned", PermaBannedLeaders)

	fullPool := make([]domain.Ci6ndexLeader, 0)
	for _, leader := range c.Leaders {
		if !contains(PermaBannedLeaders, leader.LeaderName) {
			fullPool = append(fullPool, leader)
		}

	}

	var allRolls []DraftOffering

	totalNumTries := 50
	attempt := 0
	for attempt < totalNumTries && len(allRolls) < len(c.Players) {
		slog.Info("attempting to shuffle leaders", "attempt", attempt, "strategy", c.DraftStrategy.Name, "players", c.Players)
		eligibleLeaders := fullPool
		allRolls = make([]DraftOffering, len(c.Players))

		playerIndex := 0
		exhaustedRolls := false
		for playerIndex < len(c.Players) && !exhaustedRolls {
			player := c.Players[playerIndex]
			attemptPerPlayer := 0
			numTriesPerPlayer := 10
			valid := false
			for attemptPerPlayer < numTriesPerPlayer && !valid {
				slog.Info("rolling civs for player", "player", player, "strategy", c.DraftStrategy.Name)
				shuffle := c.Functions[c.DraftStrategy.Name].shuffle
				validate := c.Functions[c.DraftStrategy.Name].validate

				offered, err := shuffle(eligibleLeaders, player, int(c.DraftStrategy.PoolSize), c.DB)
				if err != nil {
					slog.Error("failed to shuffle leaders", "error", err, "player", player, "strategy", c.DraftStrategy.Name)
					return nil, errors.Wrap(err, "failed to shuffle leaders")
				}
				if validate(offered, player, c.DB) {
					roll := DraftOffering{
						User:    player,
						Leaders: offered,
					}
					allRolls = append(allRolls, roll)
					if c.DraftStrategy.Name != "AllPick" {
						eligibleLeaders = RemoveOffered(eligibleLeaders, offered)
					}
					valid = true
				} else {
					slog.Debug("invalid roll, retrying", "player", player, "strategy", c.DraftStrategy.Name, "offered", offered)
					attemptPerPlayer++
				}
			}
			if valid {
				slog.Debug("valid roll for player", "player", player, "strategy", c.DraftStrategy.Name, "rolls", allRolls)
			} else {
				slog.Warn("failed to roll valid offering for player", "player", player, "strategy", c.DraftStrategy.Name)
				exhaustedRolls = true
			}
			playerIndex++
		}
		attempt++
	}
	if len(allRolls) < len(c.Players) {
		slog.Warn("failed to roll valid offerings for all players", "strategy", c.DraftStrategy.Name, "players", c.Players, "totalTries", totalNumTries)
		return nil, errors.New("failed to roll valid offerings for all players")
	}

	return allRolls, nil
}

func allPickValidate(leaders []domain.Ci6ndexLeader, user string, db *DatabaseOperations) bool {
	return true
}

func allPick(leaders []domain.Ci6ndexLeader, user string, poolSize int, db *DatabaseOperations) ([]domain.Ci6ndexLeader, error) {
	return leaders, nil
}

func randomPick(leaders []domain.Ci6ndexLeader, user string, poolSize int, db *DatabaseOperations) ([]domain.Ci6ndexLeader, error) {
	offering := make([]domain.Ci6ndexLeader, poolSize)
	localLeaders := leaders
	for i := 0; i < poolSize; i++ {
		randIndex := rand.N(len(localLeaders))
		offering[i] = localLeaders[randIndex]
		localLeaders = removeIndex(localLeaders, randIndex)
	}
	return leaders, nil
}

func randomPickValidate(leaders []domain.Ci6ndexLeader, user string, db *DatabaseOperations) bool {
	if !areElementsUnique(leaders) {
		return false
	}
	noRecentPick := true
	for _, leader := range leaders {
		noRecentPick = noRecentPick && hasNoRecentPick(leader, user, db)
	}
	return noRecentPick
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

func hasNoRecentPick(leader domain.Ci6ndexLeader, user string, db *DatabaseOperations) bool {
	params := domain.GetDraftPicksForUserFromLastNGamesParams{
		DiscordName: user,
		Limit:       3,
	}

	picks, err := db.Queries.GetDraftPicksForUserFromLastNGames(context.Background(), params)

	if err != nil {
		slog.Error("failed to query draft picks for user while validating recent picks",
			"params", params, "error", err)
		return false
	}
	for _, pick := range picks {
		if pick.LeaderID.Int64 == leader.ID {
			return false
		}
	}
	return true
}

func RemoveOffered(leaders []domain.Ci6ndexLeader, offered []domain.Ci6ndexLeader) []domain.Ci6ndexLeader {
	for _, off := range offered {
		for i, l := range leaders {
			if l.ID == off.ID {
				leaders = removeIndex(leaders, i)
				break
			}
		}
	}
	return leaders

}

func removeIndex(s []domain.Ci6ndexLeader, index int) []domain.Ci6ndexLeader {
	ret := make([]domain.Ci6ndexLeader, 0)
	ret = append(ret, s[:index]...)
	return append(ret, s[index+1:]...)
}

func contains(leaders []string, leaderName string) bool {
	for _, leader := range leaders {
		if leader == leaderName {
			return true
		}
	}
	return false
}
