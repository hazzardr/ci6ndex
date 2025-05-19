package ci6ndex

import (
	"ci6ndex/ci6ndex/generated"
	"context"
	"github.com/pkg/errors"
)

func (c *Ci6ndex) GetLeaders(playerID uint64, guildId uint64, rule Rule) ([]generated.Leader,
	error) {
	c.Logger.Info("Getting leaders for player", "playerID", playerID)
	ctx := context.TODO()
	db, err := c.getDB(guildId)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get database connection")
	}

	leaders, err := db.Queries.GetEligibleLeaders(ctx)
	filtered := rule.evaluate(player, leaders)
	return filtered, nil
}
