package internal

import (
	"ci6ndex/domain"
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
)

var BannedLeaders = []string{
	"ABE",
	"TOMYRIS",
	"GILGAMESH",
	"HAMMURABI",
}

type DraftOffering struct {
	Leaders []domain.Ci6ndexLeader
	User    string
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

func GetDrafts(ctx context.Context, db *DatabaseOperations, active bool) ([]domain.Ci6ndexDraft, error) {
	var drafts []domain.Ci6ndexDraft
	var err error
	if active {
		drafts, err = db.Queries.GetActiveDrafts(ctx)
	} else {
		drafts, err = db.Queries.GetDrafts(ctx)
	}
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, errors.Join(errors.New("error fetching active draft"), err)
	}
	if len(drafts) == 0 {
		return []domain.Ci6ndexDraft{}, nil
	}
	if len(drafts) > 1 {
		return nil, errors.New("more than one active draft exists! not sure which to chose")
	}
	if drafts == nil {
		return []domain.Ci6ndexDraft{}, nil
	}
	return drafts, nil
}

func CloseDraft(ctx context.Context, db *DatabaseOperations, draftId int64) (int64, error) {
	return 0, nil
}

// SubmitDraftPickRequest is the request body for submitting a draft pick
type SubmitDraftPickRequest struct {
	Leader      Leader   `json:"leader"`
	DiscordUser string   `json:"discord_user"`
	Offered     []Leader `json:"offered"`
	DraftId     int64    `json:"draft_id"`
}

//
//func SubmitDraftPick(ctx context.Context, db *DatabaseOperations, sdp *SubmitDraftPickRequest) (*domain.Ci6ndexDraftPick, error) {
//	_, err := db.Queries.GetDraft(ctx, sdp.DraftId)
//	if err != nil {
//		if errors.Is(err, pgx.ErrNoRows) {
//			return nil, errors.Join(fmt.Errorf("draft_id=%v does not exist", sdp.DraftId), err)
//		}
//		return nil, errors.Join(fmt.Errorf("error fetching draft with draft_id=%v", sdp.DraftId),
//			err)
//	}
//
//	user, err := db.Queries.GetUserByDiscordName(ctx, strings.ToLower(sdp.DiscordUser))
//	if err != nil {
//		if errors.Is(err, pgx.ErrNoRows) {
//			return nil, errors.Join(fmt.Errorf("discord user=%v does not exist", sdp.DiscordUser), err)
//		}
//		return nil, errors.Join(fmt.Errorf("error fetching discord user=%v", sdp.DiscordUser), err)
//	}
//
//	leader, err := db.Queries.GetLeaderByNameAndCiv(ctx, domain.GetLeaderByNameAndCivParams{
//		LeaderName: strings.ToUpper(sdp.Leader.Name),
//		CivName:    strings.ToUpper(sdp.Leader.Civ),
//	})
//
//	if err != nil {
//		if errors.Is(err, pgx.ErrNoRows) {
//			return nil, errors.Join(fmt.Errorf("leader=%v / civ=%v does not exist", sdp.Leader.Name, sdp.Leader.Civ), err)
//		}
//		return nil, errors.Join(fmt.Errorf("error fetching leader=%v / civ=%v", sdp.Leader.Name, sdp.Leader.Civ), err)
//	}
//
//	var leaderIDs []int64
//	for _, offeredLeader := range sdp.Offered {
//		ol, err := db.Queries.GetLeaderByNameAndCiv(ctx, domain.GetLeaderByNameAndCivParams{
//			LeaderName: strings.ToUpper(offeredLeader.Name),
//			CivName:    strings.ToUpper(offeredLeader.Civ),
//		})
//		if err != nil {
//			return nil, errors.Join(fmt.Errorf("error fetching offered leader=%v / civ=%v", offeredLeader.Name, offeredLeader.Civ), err)
//		}
//		leaderIDs = append(leaderIDs, ol.ID)
//	}
//
//	pick, err := db.Queries.SubmitDraftPick(ctx, domain.SubmitDraftPickParams{
//		DraftID:  sdp.DraftId,
//		LeaderID: pgtype.Int8{Int64: leader.ID},
//		UserID:   user.ID,
//		Offered:  leaderIDs,
//	})
//
//	if err != nil {
//		var sqlErr *pgconn.PgError
//		if errors.As(err, &sqlErr) {
//			if sqlErr.Code == "23505" {
//				return nil, errors.Join(fmt.Errorf("user=%v already has a draft pick for draft=%v", sdp.DiscordUser, draftId), err)
//			}
//		}
//		return nil, errors.Join(fmt.Errorf("error submitting draft pick for user=%v draft=%v", sdp.DiscordUser, draftId), err)
//
//	}
//	return &pick, nil
//
//}

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
