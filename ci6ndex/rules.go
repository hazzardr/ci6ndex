package ci6ndex

import (
	"ci6ndex/ci6ndex/generated"
)

// Rule is a method of filtering leaders based on player metadata
type Rule interface {
	isValid(player generated.Player, leader generated.Leader) bool
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
	filtered := make([]generated.Leader, 0)
	for _, leader := range leaders {
		if r.isValid(player, leader) {
			filtered = append(filtered, leader)
		}
	}
	return filtered
}

type NoOpRule struct{}

func (r *NoOpRule) isValid(player generated.Player, leader generated.Leader) bool {
	return true
}
func (r *NoOpRule) evaluate(player generated.Player, leaders []generated.Leader) []generated.Leader {
	return leaders
}
