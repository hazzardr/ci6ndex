package domain

import (
	"ci6ndex-bot/domain/generated"
	"context"
	"database/sql"
	"github.com/pkg/errors"
	"sync"
)

func (dbo *DatabaseOperations) GetOrCreateActiveDraft(guildId uint64) (generated.Draft, error) {
	db, err := dbo.getDB(guildId)
	if err != nil {
		return generated.Draft{}, err
	}
	d, err := db.Queries.GetActiveDraft(context.Background())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			d, err = db.Writes.CreateActiveDraft(context.Background())
			if err != nil {
				return generated.Draft{}, errors.Wrap(err, "failed to create new draft")
			}
			return d, nil
		}
		return generated.Draft{}, err
	}
	return d, nil
}

func (dbo *DatabaseOperations) SetPlayersForDraft(guildId uint64, draftId int64,
	players []generated.AddPlayerParams) []error {
	db, err := dbo.getDB(guildId)
	if err != nil {
		return []error{err}
	}
	err = db.Writes.RemovePlayersFromDraft(context.Background(), draftId)
	if err != nil {
		return []error{errors.Wrap(err, "failed to register players for draft. Unable to delete")}
	}

	if len(players) == 0 {
		dbo.logger.Debug("removed all players from draft", "draftId", draftId)
		return nil
	}

	var wg sync.WaitGroup
	errChan := make(chan error, len(players))

	for _, p := range players {
		wg.Add(1)
		go func(player generated.AddPlayerParams) {
			defer wg.Done()
			ctx := context.Background()
			_, err := db.Writes.AddPlayerToDraft(ctx, generated.AddPlayerToDraftParams{
				DraftID:  draftId,
				PlayerID: p.ID,
			})
			if err != nil {
				errChan <- errors.Wrapf(err, "failed to register player=%d for draft=%d", player.ID,
					draftId)
			} else {
				err = db.Writes.AddPlayer(ctx, player)
				if err != nil {
					errChan <- errors.Wrapf(
						err,
						"failed to add details for player=%d to database",
						player.ID,
					)
				}
			}
		}(p)
	}

	wg.Wait()
	close(errChan)

	var errs []error
	for err := range errChan {
		errs = append(errs, err)
	}

	if len(errs) > 0 {
		return errs
	}

	return nil
}

func (dbo *DatabaseOperations) GetPlayersFromActiveDraft(guildId uint64) ([]generated.Player, error) {
	db, err := dbo.getDB(guildId)
	if err != nil {
		return nil, err
	}
	players, err := db.Queries.GetPlayersFromActiveDraft(context.TODO())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return make([]generated.Player, 0), nil
		}
		return nil, err
	}
	return players, nil
}

func (dbo *DatabaseOperations) GetPlayersFromDraft(guildId uint64,
	draftId int64) ([]generated.Player, error) {
	db, err := dbo.getDB(guildId)
	if err != nil {
		return nil, err
	}
	players, err := db.Queries.GetPlayersFromDraft(context.Background(), draftId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return make([]generated.Player, 0), nil
		}
		return nil, err
	}
	return players, nil
}
