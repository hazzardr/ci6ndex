package internal

import (
	"ci6ndex/domain"
	"context"
	"encoding/json"
	"errors"
	"golang.org/x/exp/slices"
	rand2 "math/rand"
)

var BannedLeaders = []string{
	"ABE",
	"TOMYRIS",
	"GILGAMESH",
	"HAMMURABI",
}

// TODO: apply min-tiers rule
func OfferPicks(db *DatabaseOperations, draft domain.Ci6ndexDraft, numPlayers int) ([]DraftOffering, error) {
	ds, err := db.queries.GetDraftStrategy(context.Background(), draft.DraftStrategy)
	if err != nil {
		return nil, errors.Join(errors.New("error fetching draft strategy"), err)
	}

	var rules map[string]interface{}

	err = json.Unmarshal(ds.Rules, &rules)
	if err != nil {
		return nil, errors.Join(errors.New("error decoding draft strategy rules"), err)
	}

	leaders, err := db.queries.GetLeaders(context.Background())
	if err != nil {
		return nil, errors.Join(errors.New("error fetching leaders when constructing picks"), err)
	}
	var validLeaders []domain.Ci6ndexLeader

	for _, l := range leaders {
		if !slices.Contains(BannedLeaders, l.LeaderName) {
			validLeaders = append(validLeaders, l)
		}
	}

	allOffers := make([]DraftOffering, 0, numPlayers)

	rand := rules["randomize"]
	if rand == nil || rand.(bool) == false {
		// All pick
		i := 1
		for i <= numPlayers {
			do := DraftOffering{
				Leaders:  validLeaders,
				Strategy: ds,
			}
			allOffers = append(allOffers, do)
		}

		return allOffers, nil
	} else {
		// Random pick
		// Determine pool size, default is 1
		var psize int
		if rules["pool_size"] == nil {
			psize = 1
		} else {
			psize = int(rules["pool_size"].(float64))
		}

		i := 1
		for i <= numPlayers {
			var offeredLeaders []domain.Ci6ndexLeader
			pickNum := 1
			for pickNum <= psize {
				randIndex := rand2.Intn(len(validLeaders))
				offeredLeaders = append(offeredLeaders, validLeaders[randIndex])
				// We remove leaders so we don't have duplicates in the draft
				validLeaders = append(validLeaders[:randIndex], validLeaders[randIndex+1:]...)
				pickNum++
			}
			allOffers = append(allOffers, DraftOffering{
				Leaders:  offeredLeaders,
				Strategy: ds,
			})
			i++
		}
		return allOffers, nil
	}
}

type DraftOffering struct {
	Leaders  []domain.Ci6ndexLeader
	Strategy domain.Ci6ndexDraftStrategy
}
