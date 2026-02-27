package laplace

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
)

type AggregateTestSuite struct {
	*ClientTestSuite
}

func TestAggregate(t *testing.T) {
	suite.Run(t, &AggregateTestSuite{
		NewClientTestSuite(),
	})
}

func (s *AggregateTestSuite) TestGetAggregateGraphBySector() {
	client := newTestClient(s.Config)
	ctx := context.Background()

	sectors, err := client.GetAllSectors(ctx, RegionTr, LocaleTr)
	s.Require().NoError(err)
	s.Require().NotEmpty(sectors)

	sectorID := sectors[0].ID.Hex()

	resp, err := client.GetAggregateGraph(ctx, AggregatePricePeriodOneDay, RegionTr, sectorID, "", "")
	s.Require().NoError(err)
	s.Require().Greater(len(resp.Graph), 0)

	point := resp.Graph[0]
	s.Require().Greater(point.Date, int64(0))
	s.Require().Greater(point.Open, 0.0)
	s.Require().Greater(point.High, 0.0)
	s.Require().Greater(point.Low, 0.0)
	s.Require().Greater(point.Close, 0.0)
}

func (s *AggregateTestSuite) TestGetAggregateGraphByIndustry() {
	client := newTestClient(s.Config)
	ctx := context.Background()

	industries, err := client.GetAllIndustries(ctx, RegionTr, LocaleTr)
	s.Require().NoError(err)
	s.Require().NotEmpty(industries)

	industryID := industries[0].ID.Hex()

	resp, err := client.GetAggregateGraph(ctx, AggregatePricePeriodOneWeek, RegionTr, "", industryID, "")
	s.Require().NoError(err)
	s.Require().Greater(len(resp.Graph), 0)

	point := resp.Graph[0]
	s.Require().Greater(point.Date, int64(0))
	s.Require().Greater(point.Open, 0.0)
	s.Require().Greater(point.High, 0.0)
	s.Require().Greater(point.Low, 0.0)
	s.Require().Greater(point.Close, 0.0)
}

func (s *AggregateTestSuite) TestGetAggregateGraphByCollection() {
	client := newTestClient(s.Config)
	ctx := context.Background()

	collections, err := client.GetAllCollections(ctx, RegionTr, LocaleTr)
	s.Require().NoError(err)
	s.Require().NotEmpty(collections)

	collectionID := collections[0].ID.Hex()

	resp, err := client.GetAggregateGraph(ctx, AggregatePricePeriodOneMonth, RegionTr, "", "", collectionID)
	s.Require().NoError(err)
	s.Require().Greater(len(resp.Graph), 0)

	point := resp.Graph[0]
	s.Require().Greater(point.Date, int64(0))
	s.Require().Greater(point.Open, 0.0)
	s.Require().Greater(point.High, 0.0)
	s.Require().Greater(point.Low, 0.0)
	s.Require().Greater(point.Close, 0.0)
}
