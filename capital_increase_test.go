package laplace

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type CapitalIncreaseTestSuite struct {
	*ClientTestSuite
}

func TestCapitalIncrease(t *testing.T) {
	suite.Run(t, &CapitalIncreaseTestSuite{
		NewClientTestSuite(),
	})
}

func (s *CapitalIncreaseTestSuite) TestGetAllCapitalIncreases() {
	client := newTestClient(s.Config)

	ctx := context.Background()

	resp, err := client.GetAllCapitalIncreases(ctx, 1, 10, RegionTr)
	require.NoError(s.T(), err)

	require.NotEmpty(s.T(), resp.Items)
}

func (s *CapitalIncreaseTestSuite) TestGetCapitalIncreasesForInstrument() {
	client := newTestClient(s.Config)

	ctx := context.Background()

	resp, err := client.GetCapitalIncreasesForInstrument(ctx, "CANTE", 1, 10, RegionTr)
	require.NoError(s.T(), err)

	require.NotNil(s.T(), resp)
}

func (s *CapitalIncreaseTestSuite) TestGetActiveRightsForInstrument() {
	client := newTestClient(s.Config)

	ctx := context.Background()

	resp, err := client.GetActiveRightsForInstrument(ctx, "AKBNK", "2024-07-20", RegionTr)
	require.NoError(s.T(), err)
	require.NotNil(s.T(), resp)
}
