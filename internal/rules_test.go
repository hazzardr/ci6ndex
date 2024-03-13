package internal

import (
	"ci6ndex/domain"
	"ci6ndex/internal/testhelper"
	"context"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"google.golang.org/appengine/log"
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

			for _, l := range TestLeaders {
				assert.Contains(t, picks, l)
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
		{"AMERICA", "ABE"},
		{"AMERICA", "BULLMOOSE TEDDY"},
		{"AMERICA", "ROUGH RIDER TEDDY"},
		{"AMERICA", "TEDDY"},
		{"ARABIA", "SALADIN SULTAN"},
		{"ARABIA", "SALADIN VIZIR"},
		{"AUSTRALIA", "JOHN CURTIN"},
	}

	TestPlayers = []string{"Player1", "Player2", "Player3", "Player4", "Player5"}
)
