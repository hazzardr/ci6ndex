package pkg

import (
	"ci6ndex-bot/domain"
	"context"
	"encoding/json"
	"github.com/pkg/errors"
	"log/slog"
	"math/rand/v2"
)

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
type validationFunction func([]domain.Ci6ndexLeader, string, domain.Ci6ndexDraftStrategy,
	*DatabaseOperations) bool

// shuffleFunc is a function that takes in a slice of leaders to be assigned,
// and a string representing the user it's being assigned to.
// It returns the output offering based on rules defined in the function.
type shuffleFunc func([]domain.Ci6ndexLeader, string, domain.Ci6ndexDraftStrategy,
	*DatabaseOperations) ([]domain.Ci6ndexLeader, error)

func NewCivShuffler(leaders []domain.Ci6ndexLeader, players []string,
	strategy domain.Ci6ndexDraftStrategy,
	db *DatabaseOperations) *CivShuffler {

	return &CivShuffler{
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
			"RandomPickPool3": {
				shuffle:  randomPick,
				validate: randomPickValidate,
			},
			"RandomPickNoRepeats": {
				shuffle:  randomPick,
				validate: randomPickValidate,
			},
		},
	}
}

func (c *CivShuffler) Shuffle() ([]DraftOffering, error) {

	fullPool := make([]domain.Ci6ndexLeader, 0)
	for _, leader := range c.Leaders {
		if !leader.Banned {
			fullPool = append(fullPool, leader)
		}
	}
	var allRolls []DraftOffering

	totalNumTries := 2000
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
			numTriesPerPlayer := 1000
			valid := false
			for attemptPerPlayer < numTriesPerPlayer && !valid {
				slog.Debug("rolling civs for player", "player", player,
					"strategy", c.DraftStrategy.Name, "attempt", attemptPerPlayer+1)
				shuffle := c.Functions[c.DraftStrategy.Name].shuffle
				validate := c.Functions[c.DraftStrategy.Name].validate

				offered, err := shuffle(eligibleLeaders, player, c.DraftStrategy, c.DB)
				if err != nil {
					slog.Error("failed to shuffle leaders", "error", err, "player", player, "strategy", c.DraftStrategy.Name)
					return nil, errors.Wrap(err, "failed to shuffle leaders")
				}
				if validate(offered, player, c.DraftStrategy, c.DB) {
					roll := DraftOffering{
						User:    player,
						Leaders: offered,
					}
					allRolls[playerIndex] = roll
					if c.DraftStrategy.Name != "AllPick" {
						eligibleLeaders = RemoveOffered(eligibleLeaders, offered)
					}
					valid = true
				} else {
					// TODO: don't include for player? would have to know the offending leader
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
			if len(eligibleLeaders) < int(c.DraftStrategy.PoolSize) {
				slog.Warn("exhausted leader pool", "poolSize", c.DraftStrategy.PoolSize, "eligibleLeaders", len(eligibleLeaders))
				exhaustedRolls = true
			}
			playerIndex++
		}
		attempt++
	}
	//if countNonNil(allRolls) < len(c.Players) {
	//	slog.Warn("failed to roll valid offerings for all players", "strategy", c.DraftStrategy.Name, "players", c.Players, "totalTries", totalNumTries)
	//	return nil, errors.New("failed to roll valid offerings for all players")
	//}

	return allRolls, nil
}

func allPickValidate(leaders []domain.Ci6ndexLeader, user string,
	strat domain.Ci6ndexDraftStrategy, db *DatabaseOperations) bool {
	return true
}

func allPick(leaders []domain.Ci6ndexLeader, user string, strat domain.Ci6ndexDraftStrategy,
	db *DatabaseOperations) ([]domain.Ci6ndexLeader, error) {
	return leaders, nil
}

func randomPick(leaders []domain.Ci6ndexLeader, user string, strat domain.Ci6ndexDraftStrategy,
	db *DatabaseOperations) ([]domain.Ci6ndexLeader, error) {
	offering := make([]domain.Ci6ndexLeader, strat.PoolSize)
	localLeaders := leaders

	for i := 0; i < int(strat.PoolSize); i++ {
		randIndex := rand.N(len(localLeaders))
		offering[i] = localLeaders[randIndex]
		localLeaders = removeIndex(localLeaders, randIndex)
	}
	return offering, nil
}

func randomPickValidate(leaders []domain.Ci6ndexLeader, user string,
	strat domain.Ci6ndexDraftStrategy, db *DatabaseOperations) bool {

	if !areElementsUnique(leaders) {
		return false
	}
	valid := true
	if hasRules(strat) {
		var rules map[string]interface{}
		err := json.Unmarshal(strat.Rules, &rules)
		if err != nil {
			slog.Error("failed to unmarshal rules", "error", err, "strat", strat.Name)
			return false
		}
		numGames, checkNumRepeats := rules["noRepeats"]
		if checkNumRepeats {
			numGames, ok := numGames.(float64) // default serializing here no idea why
			if !ok {
				slog.Error("failed to convert noRepeats to int", "numGames", numGames)
				return false
			}
			numGamesInt := int32(numGames)
			noRecentPick := hasNoRecentPick(leaders, user, numGamesInt, db)
			valid = valid && noRecentPick
		}
		minTier, hasMinTier := rules["minTier"]
		if hasMinTier {
			minTier, ok := minTier.(float64)
			if !ok {
				slog.Error("failed to convert minTier to float", "minTier", minTier)
				return false
			}
			hasAboveTier := false
			for _, leader := range leaders {
				if leader.Tier <= minTier {
					hasAboveTier = true
					break
				}
			}
			if !hasAboveTier {
				slog.Warn("no leaders above minTier", "minTier", minTier, "user", user, "leaders", leaders)
			}
			valid = valid && hasAboveTier
		}
	}
	return valid
}

func areElementsUnique(leaders []domain.Ci6ndexLeader) bool {
	elements := make(map[domain.Ci6ndexLeader]bool)
	for _, leader := range leaders {
		if _, exists := elements[leader]; exists {
			slog.Warn("duplicate leader found in pool", "leader", leader)
			// Element already encountered, so elements are not unique
			return false
		}
		// Add element to map
		elements[leader] = true
	}
	// All elements are unique
	return true
}

func hasRules(strategy domain.Ci6ndexDraftStrategy) bool {
	return strategy.Rules != nil && len(strategy.Rules) > 0
}

func hasNoRecentPick(leaders []domain.Ci6ndexLeader, user string,
	numGames int32, db *DatabaseOperations) bool {
	params := domain.GetDraftPicksForUserFromLastNGamesParams{
		DiscordName: user,
		Limit:       numGames,
	}

	picks, err := db.Queries.GetDraftPicksForUserFromLastNGames(context.Background(), params)

	if err != nil {
		slog.Error("failed to query draft picks for user while validating recent picks",
			"params", params, "error", err)
		return false
	}
	for _, pick := range picks {
		for _, leader := range leaders {
			if pick.LeaderID.Int64 == leader.ID {
				slog.Warn("leader has been picked recently", "leader", leader, "user", user)
				return false
			}
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

func countNonNil(slice []DraftOffering) int {
	count := 0
	for _, element := range slice {
		if element.User != "" && element.Leaders != nil {
			count++
		}
	}
	return count
}
