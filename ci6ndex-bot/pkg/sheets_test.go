package pkg

import (
	"context"
	"github.com/stretchr/testify/suite"
	"testing"
)

type SheetsSuite struct {
	suite.Suite
	config *AppConfig
	db     *DatabaseOperations
	ctx    context.Context
}

func TestSheetsSuite(t *testing.T) {
	suite.Run(t, new(SheetsSuite))
}

func (suite *SheetsSuite) SetupSuite() {
	suite.ctx = context.Background()
	suite.config = &AppConfig{
		GoogleCloudCredentialsLocation: "testhelper/testdata/creds.json",
		CivRankingSheetId:              "1",
	}
}

//
//func (suite *SheetsSuite) TestSuccessfullyReadFromSheets() {
//	GetRankingsFromSheets(suite.config, suite.ctx)
//}
