package internal

import (
	"ci6ndex/internal/testhelper"
	"context"
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

func (suite *RulesSuite) TestAddUsersFromFile_FailureCases() {
	type test struct {
		name string
		path string
	}

	tests := []test{
		{"EmptyFile", "testhelper/testdata/empty.json"},
		{"FileDoesNotExist", "testhelper/testdata/nonexistent.json"},
		{"InvalidFileType", "testhelper/testdata/invalid.yaml"},
		{"MalformedJSON", "testhelper/testdata/malformed.json"},
	}

	for _, tc := range tests {
		suite.T().Run(tc.name, func(t *testing.T) {
			err := AddUsersFromFile(tc.path, suite.db)
			suite.Error(err)
		})

	}
}
