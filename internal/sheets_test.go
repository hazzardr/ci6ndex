package internal

import (
	"context"
	"github.com/stretchr/testify/suite"
	"testing"
)

type SheetsSuite struct {
	suite.Suite
	db  *DatabaseOperations
	ctx context.Context
}

func TestSheetsSuite(t *testing.T) {
	suite.Run(t, new(SheetsSuite))
}

func (suite *SheetsSuite) TestSuccessfullyReadFromSheets() {
	suite.ctx = context.Background()
}
