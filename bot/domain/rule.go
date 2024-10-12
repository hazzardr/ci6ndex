package domain

import "ci6ndex-bot/domain/generated"

type Rule interface {
	isValid(player generated.Player, leader generated.Leader) bool
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

type NoOpRule struct{}

func (r *NoOpRule) isValid(player generated.Player, leader generated.Leader) bool {
	return true
}

type OfferedPicks map[int64][]generated.Leader
type EligibleLeadersPerPlayer map[int64]map[Rule][]generated.Leader

type RuleSet struct {
	rules EligibleLeadersPerPlayer
}

func NewRulesetWithAtLeastOneTierThree(players []generated.Player, poolSize uint8) RuleSet {
	r := make(map[int64]map[Rule][]generated.Leader)
	for _, player := range players {
		r[player.ID] = make(map[Rule][]generated.Leader)
		for i := 0; i < int(poolSize); i++ {
			if i == 0 {
				r[player.ID][&MinTierRule{minTier: 4}] = make([]generated.Leader, 0)
			}
			r[player.ID][&NoOpRule{}] = make([]generated.Leader, 0)
		}
	}
	return RuleSet{
		rules: r,
	}
}

func (rs *RuleSet) AddLeaders(leaders []generated.Leader) {
	for _, leader := range leaders {
		for _, ruleMap := range rs.rules {
			for rule := range ruleMap {
				if rule.isValid(generated.Player{}, leader) {
					ruleMap[rule] = append(ruleMap[rule], leader)
				}
			}
		}
	}
}

func (rs *RuleSet) Evaluate() (OfferedPicks, error) {
	result := make(map[int64][]generated.Leader)
	for playerId, ruleMap := range rs.rules {
		for rule, leaders := range ruleMap {

		}
	}
	return result, nil
}
