package ci6ndex

import (
	"ci6ndex/ci6ndex/generated"
	"context"
	"database/sql"
	"github.com/pkg/errors"
	"log/slog"
	"sync"
)

func (c *Ci6ndex) GetOrCreateActiveDraft(guildId uint64) (generated.Draft, error) {
	db, err := c.getDB(guildId)
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

func (c *Ci6ndex) SetPlayersForDraft(guildId uint64, draftId int64,
	players []generated.AddPlayerParams) []error {
	db, err := c.getDB(guildId)
	if err != nil {
		return []error{err}
	}
	err = db.Writes.RemovePlayersFromDraft(context.Background(), draftId)
	if err != nil {
		return []error{errors.Wrap(err, "failed to register players for draft. Unable to delete")}
	}

	if len(players) == 0 {
		slog.Debug("removed all players from draft", "draftId", draftId)
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
				PlayerID: player.ID,
			})
			if err != nil {
				errChan <- errors.Wrapf(err, "failed to register player=%d for draft=%d", player.ID,
					draftId)
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
