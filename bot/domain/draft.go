package domain

import (
	"ci6ndex-bot/domain/generated"
	"context"
	"database/sql"
	"github.com/pkg/errors"
)

func (dbo *DatabaseOperations) GetOrCreateActiveDraft(guildId uint64) (generated.Draft, error) {
	db, err := dbo.getDB(guildId)
	if err != nil {
		return generated.Draft{}, err
	}
	d, err := db.Queries.GetActiveDraft(context.Background())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			d, err = db.Writes.CreateDraft(context.Background(), generated.CreateDraftParams{
				Active: true,
			})
			if err != nil {
				return generated.Draft{}, errors.Wrap(err, "failed to create new draft")
			}
			return d, nil
		}
		return generated.Draft{}, err
	}
	return d, nil
}
