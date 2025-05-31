package ci6ndex

import (
	"ci6ndex/ci6ndex/generated"
)

type RuleType string

const (
	All        RuleType = "all"
	AtLeastOne RuleType = "atLeastOne"
)

// Rule is a method of filtering leaders based on player metadata
type Rule interface {
	IsValid(player generated.Player, leader generated.Leader) bool
	Filter(player generated.Player, leaders []generated.Leader) []generated.Leader
	Type() RuleType
}

// MinTierRule filters leaders based on a minimum tier requirement.
type MinTierRule struct {
	MinTier float64
}

func (r *MinTierRule) IsValid(player generated.Player, leader generated.Leader) bool {
	if leader.Tier <= r.MinTier {
		return true
	}
	return false
}

func (r *MinTierRule) Filter(player generated.Player, leaders []generated.Leader) []generated.Leader {
	filtered := make([]generated.Leader, 0)
	for _, leader := range leaders {
		if r.IsValid(player, leader) {
			filtered = append(filtered, leader)
		}
	}
	return filtered
}

func (r *MinTierRule) Type() RuleType {
	return AtLeastOne
}

type NoOpRule struct{}

func (r *NoOpRule) IsValid(player generated.Player, leader generated.Leader) bool {
	return true
}
func (r *NoOpRule) Filter(player generated.Player, leaders []generated.Leader) []generated.Leader {
	return leaders
}
func (r *NoOpRule) Type() RuleType {
	return AtLeastOne
}
