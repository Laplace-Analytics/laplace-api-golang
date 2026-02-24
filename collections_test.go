package laplace

import (
	"context"
	"testing"

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

func (s *CollectionsTestSuite) TestGetAllCollections() {
	client := newTestClient(s.Config)

	ctx := context.Background()

	resp, err := client.GetAllCollections(ctx, RegionTr, LocaleTr)
	s.Require().NoError(err)
	s.Require().NotEmpty(resp)

	hasStocks := false
	for _, collection := range resp {
		s.Require().NotEmpty(collection.ID)
		s.Require().NotEmpty(collection.Title)
		s.Require().NotEmpty(collection.ImageUrl)
		s.Require().NotEmpty(collection.AvatarUrl)
		if collection.NumStocks > 0 {
			hasStocks = true
		}
	}
	s.Require().True(hasStocks)
}

func (s *CollectionsTestSuite) TestGetCollectionDetail() {
	client := newTestClient(s.Config)

	ctx := context.Background()

	resp, err := client.GetCollectionDetail(ctx, "620f455a0187ade00bb0d55f", RegionTr, LocaleTr)
	s.Require().NoError(err)
	s.Require().NotEmpty(resp)
	s.Require().NotEmpty(resp.ID)
	s.Require().NotEmpty(resp.Title)
	s.Require().NotEmpty(resp.Region)
	s.Require().Greater(resp.NumStocks, 0)

	hasValidRegion := false
	for _, region := range resp.Region {
		if region == RegionTr || region == RegionUs {
			hasValidRegion = true
			break
		}
	}
	s.Require().True(hasValidRegion)

	for _, stock := range resp.Stocks {
		s.Require().NotEmpty(stock.ID)
		s.Require().NotEmpty(stock.Name)
		s.Require().NotEmpty(stock.Symbol)
		s.Require().NotEmpty(stock.SectorId)
		s.Require().NotEmpty(stock.IndustryId)
		s.Require().NotEmpty(stock.AssetType)
		s.Require().NotEmpty(stock.UpdatedDate)
	}
}
