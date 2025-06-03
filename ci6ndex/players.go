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

func (c *Ci6ndex) GetPlayer(
	ctx context.Context,
	guildId uint64,
	playerId int64,
) (*generated.Player, error) {
	db, err := c.getDB(guildId)
	if err != nil {
		return nil, errors.Join(err, errors.New("failed to get db"))
	}
	p, err := db.Queries.GetPlayer(ctx, playerId)
	if err != nil {
		return nil, errors.Join(err, errors.New("failed to get player"))
	}
	return &p, nil
}

func (c *Ci6ndex) GetPlayers(ctx context.Context, guildId uint64) ([]generated.Player, error) {
	db, err := c.getDB(guildId)
	if err != nil {
		return nil, errors.Join(err, errors.New("failed to get db"))
	}
	players, err := db.Queries.GetPlayers(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return make([]generated.Player, 0), nil
		}
		return nil, errors.Join(err, errors.New("failed to get players"))
	}
	return players, nil
}

func (c *Ci6ndex) AddPlayer(ctx context.Context, guildId uint64, params generated.AddPlayerParams) error {
	db, err := c.getDB(guildId)
	if err != nil {
		return errors.Join(err, errors.New("failed to get db"))
	}
	err = db.Writes.AddPlayer(ctx, params)
	if err != nil {
		return errors.Join(err, errors.New("failed to add player"))
	}
	return nil
}
