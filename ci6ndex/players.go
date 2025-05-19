package ci6ndex

import (
	"ci6ndex/ci6ndex/generated"
	"context"
	"database/sql"
	"errors"
)

func (c *Ci6ndex) GetPlayersFromActiveDraft(guildId uint64) ([]generated.Player, error) {
	db, err := c.getDB(guildId)
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

func (c *Ci6ndex) GetPlayersFromDraft(guildId uint64,
	draftId int64) ([]generated.Player, error) {
	db, err := c.getDB(guildId)
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
