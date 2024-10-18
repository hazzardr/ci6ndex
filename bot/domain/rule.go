package domain

import (
	"ci6ndex-bot/domain/generated"
	"fmt"
	"slices"
)

type FullPoolError struct{}

func (e FullPoolError) Error() string {
	return fmt.Sprint("attempted to add rule to full pool")
}

type Rule interface {
	isValid(player generated.Player, leader generated.Leader) bool
	// Evaluate should
	evaluate(player generated.Player, leaders []generated.Leader) []generated.Leader
}

type MinTierRule struct {
	minTier float64
}

func (r *MinTierRule) isValid(player generated.Player, leader generated.Leader) bool {
	if leader.Tier <= r.minTier {
		return true
	}
	return false
}

func (r *MinTierRule) evaluate(player generated.Player, leaders []generated.Leader) []generated.Leader {
	return slices.DeleteFunc(leaders, func(l generated.Leader) bool {
		return !r.isValid(player, l)
	})
}

type NoOpRule struct{}

func (r *NoOpRule) isValid(player generated.Player, leader generated.Leader) bool {
	return true
}
func (r *NoOpRule) evaluate(player generated.Player, leaders []generated.Leader) []generated.Leader {
	return leaders
}
