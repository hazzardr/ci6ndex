package pkg

import (
	"ci6ndex/domain"
	"ci6ndex/pkg/testhelper"
	"context"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"google.golang.org/appengine/log"
	"slices"
	"testing"
)

type RulesSuite struct {
	suite.Suite
	db            *DatabaseOperations
	testContainer *testhelper.TestDatabase
	ctx           context.Context
}

func TestRulesSuite(t *testing.T) {
	suite.Run(t, new(RulesSuite))
}

// Create a database container for testing, and initialize our client against it.
func (suite *RulesSuite) SetupSuite() {
	suite.ctx = context.Background()
	testContainer, err := testhelper.CreateTestDatabase(suite.ctx)
	if err != nil {
		suite.T().Fatal(err)
	}
	suite.testContainer = testContainer
	testURL := testContainer.ConnectionString
	suite.db, err = NewDBConnection(testURL)
	if err != nil {
		suite.T().Fatal(err)
	}
}

func (suite *RulesSuite) TearDownSuite() {
	err := suite.testContainer.Terminate(suite.ctx)
	if err != nil {
		log.Errorf(suite.ctx, "Error close database container: %v", err)
		suite.T().Fatal(err)
	}
	suite.db.Close()
}

func (suite *RulesSuite) TearDownTest() {
	err := suite.db.Queries.WipeTables(suite.ctx)
	if err != nil {
		log.Errorf(suite.ctx, "Error truncating tables, could not clean DB: %v", err)
		suite.T().Fatal(err)
	}
}

func (suite *RulesSuite) TestAllPick() {
	type test struct {
		name     string
		leaders  []domain.Ci6ndexLeader
		players  []string
		strategy domain.Ci6ndexDraftStrategy
	}

	allPick, err := suite.db.Queries.CreateDraftStrategy(
		suite.ctx, domain.CreateDraftStrategyParams{
			Name: "AllPick", Description: "test", Randomize: false,
		})
	if err != nil {
		suite.T().Fatal(errors.Wrap(err, "failed to setup test"))
	}

	err = seedTestLeaders(suite.ctx, suite.db)
	if err != nil {
		suite.T().Fatal(errors.Wrap(err, "failed to seed test leaders"))
	}

	leaders, err := suite.db.Queries.GetLeaders(suite.ctx)
	if err != nil {
		suite.T().Fatal(errors.Wrap(err, "failed to get test leaders"))
	}
	tests := []test{
		{"AllPick", leaders, TestPlayers, allPick},
	}

	for _, tc := range tests {
		suite.T().Run(tc.name, func(t *testing.T) {
			shuffler := NewCivShuffler(tc.leaders, tc.players, tc.strategy, suite.db)
			picks, err := shuffler.Shuffle()
			if err != nil {
				return
			}

			if len(picks) != len(tc.players) {
				suite.FailNow("number of picks did not match number of players")
			}

			// Each player should be offered every leader
			for _, tl := range TestLeaders {
				for _, offer := range picks {
					assert.Truef(suite.T(), slices.ContainsFunc(offer.Leaders,
						func(l domain.Ci6ndexLeader) bool {
							return l.LeaderName == tl.LeaderName && l.CivName == tl.CivName
						}), "leader not found in offering")
				}
			}
		})

	}
}

func (suite *RulesSuite) TestRandomPick() {
	type test struct {
		name     string
		leaders  []domain.Ci6ndexLeader
		players  []string
		strategy domain.Ci6ndexDraftStrategy
	}

	poolSize := 3
	allRand, err := suite.db.Queries.CreateDraftStrategy(
		suite.ctx, domain.CreateDraftStrategyParams{
			Name: "RandomPick", Description: "test", Randomize: true, PoolSize: int32(poolSize),
		})
	if err != nil {
		suite.T().Fatal(errors.Wrap(err, "failed to setup test"))
	}

	err = seedTestLeaders(suite.ctx, suite.db)
	if err != nil {
		suite.T().Fatal(errors.Wrap(err, "failed to seed test leaders"))
	}

	leaders, err := suite.db.Queries.GetLeaders(suite.ctx)
	if err != nil {
		suite.T().Fatal(errors.Wrap(err, "failed to get test leaders"))
	}
	tests := []test{
		{"RandomPick", leaders, TestPlayers, allRand},
	}

	for _, tc := range tests {
		suite.T().Run(tc.name, func(t *testing.T) {
			shuffler := NewCivShuffler(tc.leaders, tc.players, tc.strategy, suite.db)
			picks, err := shuffler.Shuffle()
			if err != nil {
				suite.FailNow(err.Error())
			}

			if len(picks) != len(tc.players) {
				suite.FailNow("number of picks did not match number of players")
			}

			// Each offered leader should be globally unique
			seen := make(map[string]bool)
			for _, offer := range picks {
				assert.Len(suite.T(), offer.Leaders, poolSize)
				for _, l := range offer.Leaders {
					assert.Falsef(suite.T(), seen[l.LeaderName], "leader was offered twice!")
					seen[l.LeaderName] = true
				}
			}
			for name := range seen {
				nameExists := false
				for _, tl := range TestLeaders {
					if tl.LeaderName == name {
						nameExists = true
						break
					}
				}
				assert.Truef(suite.T(), nameExists, "offered leader is not from original pool")
			}
		})

	}
}

func seedTestLeaders(ctx context.Context, db *DatabaseOperations) error {
	var lp []domain.CreateLeadersParams
	for _, leader := range TestLeaders {
		lp = append(lp, domain.CreateLeadersParams{
			LeaderName: leader.LeaderName,
			CivName:    leader.CivName,
		})
	}

	_, err := db.Queries.CreateLeaders(ctx, lp)
	if err != nil {
		return err
	}
	return nil
}

type testLeader struct {
	LeaderName string
	CivName    string
}

var (
	TestLeaders = []testLeader{
		{"BULLMOOSE TEDDY", "AMERICA"},
		{"ROUGH RIDER TEDDY", "AMERICA"},
		{"TEDDY", "AMERICA"},
		{"SALADIN SULTAN", "ARABIA"},
		{"SALADIN VIZIR", "ARABIA"},
		{"JOHN CURTIN", "AUSTRALIA"},
	}

	TestPlayers = []string{"Player1", "Player2"}
)
