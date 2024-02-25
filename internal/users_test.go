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
		log.Errorf(suite.ctx, "Error creating database container: %v", err)
		suite.T().Fatal(err)
	}
	suite.testContainer = testContainer
	testURL := testContainer.ConnectionString
	suite.db, err = newDBConnection(testURL)
	if err != nil {
		log.Errorf(suite.ctx, "Error creating database connection: %v", err)
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

func (suite *UsersSuite) TestAddUsersFromFile() {
	t := suite.T()
	err := AddUsersFromFile("testhelper/testdata/single_user.json", suite.db)
	suite.NoError(err)
	actual, err := suite.db.queries.GetUserByName(suite.ctx, "username")
	suite.NoError(err)
	assert.Equal(t, "username", actual.Name)
	assert.Equal(t, "discord_name", actual.DiscordName)
}
