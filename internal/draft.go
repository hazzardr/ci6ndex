package internal

import (
	"ci6ndex/domain"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"golang.org/x/exp/slices"
	rand2 "math/rand"
	"net/http"
)

var BannedLeaders = []string{
	"ABE",
	"TOMYRIS",
	"GILGAMESH",
	"HAMMURABI",
}

type DraftOffering struct {
	Leaders  []domain.Ci6ndexLeader
	Strategy domain.Ci6ndexDraftStrategy
}

type CreateDraftRequest struct {
	DraftStrategy string `json:"draft_strategy"`
}

func CreateDraft(w http.ResponseWriter, req *http.Request) {
	var cdr CreateDraftRequest
	err := json.NewDecoder(req.Body).Decode(&cdr)
	if err != nil {
		w.WriteHeader(400)
		_ = json.NewEncoder(w).Encode("could not parse strategy from request body")
		return
	}

	strategy := cdr.DraftStrategy

	_, err = db.Queries.GetDraftStrategy(req.Context(), strategy)
	if err != nil {
		w.WriteHeader(422)
		_ = json.NewEncoder(w).Encode(fmt.Sprintf("draft_strategy=%v does not exist", strategy))
		return
	}

	draft, err := db.Queries.CreateDraft(req.Context(), strategy)
	if err != nil {
		w.WriteHeader(500)
		_ = json.NewEncoder(w).Encode(err)
		return

	}

	w.WriteHeader(201)
	err = json.NewEncoder(w).Encode(draft)
	if err != nil {
		w.WriteHeader(500)
		_ = json.NewEncoder(w).Encode(err)
		return
	}
}

// TODO: apply min-tiers rule
func OfferPicks(db *DatabaseOperations, draft domain.Ci6ndexDraft, numPlayers int) ([]DraftOffering, error) {
	ds, err := db.Queries.GetDraftStrategy(context.Background(), draft.DraftStrategy)
	if err != nil {
		return nil, errors.Join(errors.New("error fetching draft strategy"), err)
	}

	var rules map[string]interface{}

	err = json.Unmarshal(ds.Rules, &rules)
	if err != nil {
		return nil, errors.Join(errors.New("error decoding draft strategy rules"), err)
	}

	leaders, err := db.Queries.GetLeaders(context.Background())
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
