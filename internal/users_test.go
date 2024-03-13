package internal

import (
	"ci6ndex/internal/testhelper"
	"context"
	"github.com/magiconair/properties/assert"
	"github.com/stretchr/testify/suite"
	"google.golang.org/appengine/log"
	"testing"
)

type UsersSuite struct {
	suite.Suite
	db            *DatabaseOperations
	testContainer *testhelper.TestDatabase
	ctx           context.Context
}

func TestUsersSuite(t *testing.T) {
	suite.Run(t, new(UsersSuite))
}

// Create a database container for testing, and initialize our client against it.
func (suite *UsersSuite) SetupSuite() {
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

func (suite *UsersSuite) TearDownSuite() {
	err := suite.testContainer.Terminate(suite.ctx)
	if err != nil {
		log.Errorf(suite.ctx, "Error close database container: %v", err)
		suite.T().Fatal(err)
	}
	suite.db.Close()
}

func (suite *UsersSuite) TearDownTest() {
	err := suite.db.Queries.WipeTables(suite.ctx)
	if err != nil {
		log.Errorf(suite.ctx, "Error truncating tables, could not clean DB: %v", err)
		suite.T().Fatal(err)
	}
}

func (suite *UsersSuite) TestAddUsersFromFile_Success() {
	t := suite.T()
	exists := fileExists("testhelper/testdata/single_user.json")
	suite.True(exists)
	err := AddUsersFromFile("testhelper/testdata/single_user.json", suite.db)
	suite.NoError(err)
	actual, err := suite.db.Queries.GetUserByName(suite.ctx, "username")
	suite.NoError(err)
	assert.Equal(t, "username", actual.Name)
	assert.Equal(t, "discord_name", actual.DiscordName)
}

func (suite *UsersSuite) TestAddUsersFromFile_MultipleUsers() {
	t := suite.T()
	exists := fileExists("testhelper/testdata/multiple_users.json")
	suite.True(exists)
	err := AddUsersFromFile("testhelper/testdata/multiple_users.json", suite.db)
	suite.NoError(err)
	users, err := suite.db.Queries.GetUsers(suite.ctx)
	suite.NoError(err)
	assert.Equal(t, 2, len(users))
}

func (suite *UsersSuite) TestAddUsersFromFile_UserAlreadyExists() {
	exists := fileExists("testhelper/testdata/single_user.json")
	suite.True(exists)
	err := AddUsersFromFile("testhelper/testdata/single_user.json", suite.db)
	suite.NoError(err)

	err = AddUsersFromFile("testhelper/testdata/single_user.json", suite.db)
	suite.Error(err)
}

func (suite *UsersSuite) TestAddUsersFromFile_FailureCases() {
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
