package domain

import (
	"context"
	"github.com/pkg/errors"
)

func (dbo *DatabaseOperations) RollForPlayers(guildId uint64, poolSize uint8) error {
	ctx := context.TODO()
	db, err := dbo.getDB(guildId)
	if err != nil {
		return errors.Wrap(err, "failed to get database connection")
	}
	players, err := dbo.GetPlayersFromActiveDraft(guildId)
	if err != nil {
		return errors.Wrap(err, "failed to get players from active draft")
	}
	leaders, err := db.Queries.GetEligibleLeaders(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to get eligible leaders")
	}
}
