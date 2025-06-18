package ci6ndex

import (
	"ci6ndex/ci6ndex/generated"
	"context"
	"database/sql"
	"errors"
	"fmt"
)

func (c *Ci6ndex) GetLeadersInRange(guildId, offset, limit uint64) ([]generated.Leader, error) {
	if offset < 0 || limit < 0 {
		return nil, errors.New(fmt.Sprintf("invalid offset or limit. offset: %d limit: %d", offset, limit))
	}
	db, err := c.getDB(guildId)
	if err != nil {
		return nil, err
	}
	ctx := context.Background()
	//TODO:
	leaders, err := db.Queries.GetLeadersByLimitAndOffset(ctx, generated.GetLeadersByLimitAndOffsetParams{
		Limit: int64(limit), Offset: int64(offset),
	})

	if err != nil {
		return nil, errors.Join(err, errors.New(fmt.Sprintf("failed to get leaders between offset and limit %d and %d", offset, limit)))
	}
	return leaders, nil
}

func (c *Ci6ndex) GetLeaderById(guildId uint64, leaderId uint64) (generated.Leader, error) {
	db, err := c.getDB(guildId)
	if err != nil {
		return generated.Leader{}, err
	}
	
	ctx := context.Background()
	leader, err := db.Queries.GetLeaderById(ctx, int64(leaderId))
	
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return generated.Leader{}, errors.New(fmt.Sprintf("leader with ID %d not found", leaderId))
		}
		return generated.Leader{}, errors.Join(err, errors.New(fmt.Sprintf("failed to fetch leader with ID %d", leaderId)))
	}
	
	return leader, nil
}
