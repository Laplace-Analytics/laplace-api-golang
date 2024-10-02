package laplace

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type SearchTestSuite struct {
	*ClientTestSuite
}

func TestSearch(t *testing.T) {
	suite.Run(t, &SearchTestSuite{
		NewClientTestSuite(),
	})
}

func (s *SearchTestSuite) TestSearchStock() {
	client := newTestClient(s.Config)

	ctx := context.Background()

	resp, err := client.Search(ctx, "TUPRS", []SearchType{SearchTypeStock}, RegionTr, LocaleTr, 0, PageSize10)
	require.NoError(s.T(), err)

	require.NotEmpty(s.T(), resp)
}

func (s *SearchTestSuite) TestSearchIndustry() {
	client := newTestClient(s.Config)

	ctx := context.Background()

	resp, err := client.Search(ctx, "Hava Taşımacılığı", []SearchType{SearchTypeIndustry}, RegionTr, LocaleTr, 0, PageSize10)
	require.NoError(s.T(), err)

	require.NotEmpty(s.T(), resp)
}

func (s *SearchTestSuite) TestSearchAllTypes() {
	client := newTestClient(s.Config)

	ctx := context.Background()

	resp, err := client.Search(ctx, "Ab", []SearchType{
		SearchTypeStock,
		SearchTypeIndustry,
		SearchTypeSector,
		SearchTypeCollection,
	}, RegionUs, LocaleTr, 0, PageSize10)
	require.NoError(s.T(), err)

	require.NotEmpty(s.T(), resp)
}
