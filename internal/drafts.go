package internal

import (
	"ci6ndex/domain"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"golang.org/x/exp/slices"
	rand2 "math/rand"
	"strings"
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

func CreateDraft(ctx context.Context, cdr CreateDraftRequest, db *DatabaseOperations) (*domain.Ci6ndexDraft, error) {
	strategy := cdr.DraftStrategy

	_, err := db.Queries.GetDraftStrategy(ctx, strategy)
	if err != nil {
		return nil, errors.Join(errors.New("draft strategy does not exist"), err)
	}

	draft, err := db.Queries.CreateDraft(ctx, strategy)
	if err != nil {
		return nil, errors.Join(errors.New("error creating draft"), err)
	}
	return &draft, nil
}

// GetActiveDraft returns the active draft, if one exists. only one draft can be active at a time.
func GetActiveDraft(ctx context.Context, db *DatabaseOperations) (*domain.Ci6ndexDraft, error) {
	drafts, err := db.Queries.GetActiveDrafts(ctx)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, errors.Join(errors.New("error fetching active draft"), err)
	}
	if len(drafts) == 0 {
		return nil, nil
	}
	if len(drafts) > 1 {
		return nil, errors.New("more than one active draft exists! not sure which to chose")
	}
	return &drafts[0], nil
}

func CloseDraft(ctx context.Context, db *DatabaseOperations, draftId int64) (int64, error) {
	return 0, nil
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

// SubmitDraftPickRequest is the request body for submitting a draft pick
type SubmitDraftPickRequest struct {
	Leader      Leader   `json:"leader"`
	DiscordUser string   `json:"discord_user"`
	Offered     []Leader `json:"offered"`
}

func SubmitDraftPick(ctx context.Context, db *DatabaseOperations, sdp SubmitDraftPickRequest, draftId int64) (*domain.Ci6ndexDraftPick, error) {
	_, err := db.Queries.GetDraft(ctx, draftId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.Join(fmt.Errorf("draft_id=%v does not exist", draftId), err)
		}
		return nil, errors.Join(fmt.Errorf("error fetching draft with draft_id=%v", draftId), err)
	}

	user, err := db.Queries.GetUserByDiscordName(ctx, strings.ToLower(sdp.DiscordUser))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.Join(fmt.Errorf("discord user=%v does not exist", sdp.DiscordUser), err)
		}
		return nil, errors.Join(fmt.Errorf("error fetching discord user=%v", sdp.DiscordUser), err)
	}

	leader, err := db.Queries.GetLeaderByNameAndCiv(ctx, domain.GetLeaderByNameAndCivParams{
		LeaderName: strings.ToUpper(sdp.Leader.Name),
		CivName:    strings.ToUpper(sdp.Leader.Civ),
	})

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.Join(fmt.Errorf("leader=%v / civ=%v does not exist", sdp.Leader.Name, sdp.Leader.Civ), err)
		}
		return nil, errors.Join(fmt.Errorf("error fetching leader=%v / civ=%v", sdp.Leader.Name, sdp.Leader.Civ), err)
	}

	var leaderIDs []int64
	for _, offeredLeader := range sdp.Offered {
		ol, err := db.Queries.GetLeaderByNameAndCiv(ctx, domain.GetLeaderByNameAndCivParams{
			LeaderName: strings.ToUpper(offeredLeader.Name),
			CivName:    strings.ToUpper(offeredLeader.Civ),
		})
		if err != nil {
			return nil, errors.Join(fmt.Errorf("error fetching offered leader=%v / civ=%v", offeredLeader.Name, offeredLeader.Civ), err)
		}
		leaderIDs = append(leaderIDs, ol.ID)
	}

	pick, err := db.Queries.SubmitDraftPick(ctx, domain.SubmitDraftPickParams{
		DraftID:  draftId,
		LeaderID: pgtype.Int8{Int64: leader.ID},
		UserID:   user.ID,
		Offered:  leaderIDs,
	})

	if err != nil {
		var sqlErr *pgconn.PgError
		if errors.As(err, &sqlErr) {
			if sqlErr.Code == "23505" {
				return nil, errors.Join(fmt.Errorf("user=%v already has a draft pick for draft=%v", sdp.DiscordUser, draftId), err)
			}
		}
		return nil, errors.Join(fmt.Errorf("error submitting draft pick for user=%v draft=%v", sdp.DiscordUser, draftId), err)

	}

	w.WriteHeader(201)
	err = json.NewEncoder(w).Encode(pick)
	if err != nil {
		w.WriteHeader(500)
		_ = json.NewEncoder(w).Encode(fmt.Sprint("failed to encode draft pick, but it was submitted successfully"))
	}

}

func GetDraftStrategies(ctx context.Context, db *DatabaseOperations) ([]domain.Ci6ndexDraftStrategy, error) {
	strats, err := db.Queries.GetDraftStrategies(ctx)
	if err != nil {
		return nil, errors.Join(errors.New("error fetching draft strategies"), err)
	}
	if strats == nil {
		strats = []domain.Ci6ndexDraftStrategy{}
	}
	return strats, nil
}

func GetDraftStrategy(ctx context.Context, db *DatabaseOperations, name string) (*domain.Ci6ndexDraftStrategy, error) {
	strat, err := db.Queries.GetDraftStrategy(ctx, name)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.Join(fmt.Errorf("draft_strategy=%v does not exist", name), err)
		}
		return nil, errors.Join(fmt.Errorf("error fetching draft_strategy=%v", name), err)
	}

	return &strat, nil
}
