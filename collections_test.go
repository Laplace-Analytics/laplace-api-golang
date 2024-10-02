package laplace

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type CollectionsTestSuite struct {
	*ClientTestSuite
}

func TestCollections(t *testing.T) {
	suite.Run(t, &CollectionsTestSuite{
		NewClientTestSuite(),
	})
}

func (s *CollectionsTestSuite) TestGetAllIndustries() {
	client := newTestClient(s.Config)

	ctx := context.Background()

	resp, err := client.GetAllIndustries(ctx, RegionTr, LocaleTr)
	require.NoError(s.T(), err)

	require.NotEmpty(s.T(), resp)
}

func (s *CollectionsTestSuite) TestGetIndustryDetails() {
	client := newTestClient(s.Config)

	ctx := context.Background()

	resp, err := client.GetIndustryDetail(ctx, "65533e441fa5c7b58afa0944", RegionTr, LocaleTr)
	require.NoError(s.T(), err)

	require.NotEmpty(s.T(), resp)
}

func (s *CollectionsTestSuite) TestGetSectorDetails() {
	client := newTestClient(s.Config)

	ctx := context.Background()

	resp, err := client.GetSectorDetail(ctx, "65533e047844ee7afe9941b9", RegionTr, LocaleTr)
	require.NoError(s.T(), err)

	require.NotEmpty(s.T(), resp)
}
