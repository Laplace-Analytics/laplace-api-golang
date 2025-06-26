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

func (s *CollectionsTestSuite) TestGetAllIndustries() {
	client := newTestClient(s.Config)

	ctx := context.Background()

	resp, err := client.GetAllIndustries(ctx, RegionTr, LocaleTr)
	s.Require().NoError(err)
	s.Require().NotEmpty(resp)

	hasStocks := false
	for _, industry := range resp {
		s.Require().NotEmpty(industry.ID)
		s.Require().NotEmpty(industry.Title)
		if industry.NumStocks > 0 {
			hasStocks = true
		}
	}
	s.Require().True(hasStocks)
}

func (s *CollectionsTestSuite) TestGetIndustryDetails() {
	client := newTestClient(s.Config)

	ctx := context.Background()

	resp, err := client.GetIndustryDetail(ctx, "65533e441fa5c7b58afa0944", RegionTr, LocaleTr)
	s.Require().NoError(err)
	s.Require().NotEmpty(resp)
	s.Require().NotEmpty(resp.ID)
	s.Require().NotEmpty(resp.Title)
	s.Require().NotEmpty(resp.Region)
	s.Require().Greater(resp.NumStocks, 0)

	hasValidRegion := false
	for _, region := range resp.Region {
		if region == string(RegionTr) || region == string(RegionUs) {
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
	}
}

func (s *CollectionsTestSuite) TestGetAllSectors() {
	client := newTestClient(s.Config)

	ctx := context.Background()

	resp, err := client.GetAllSectors(ctx, RegionTr, LocaleTr)
	s.Require().NoError(err)
	s.Require().NotEmpty(resp)

	hasStocks := false
	for _, sector := range resp {
		s.Require().NotEmpty(sector.ID)
		s.Require().NotEmpty(sector.Title)
		if sector.NumStocks > 0 {
			hasStocks = true
		}
	}
	s.Require().True(hasStocks)
}

func (s *CollectionsTestSuite) TestGetSectorDetails() {
	client := newTestClient(s.Config)

	ctx := context.Background()

	resp, err := client.GetSectorDetail(ctx, "65533e047844ee7afe9941b9", RegionTr, LocaleTr)
	s.Require().NoError(err)
	s.Require().NotEmpty(resp)
	s.Require().NotEmpty(resp.ID)
	s.Require().NotEmpty(resp.Title)
	s.Require().NotEmpty(resp.Region)
	s.Require().Greater(resp.NumStocks, 0)

	hasValidRegion := false
	for _, region := range resp.Region {
		if region == string(RegionTr) || region == string(RegionUs) {
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
	}
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
		if region == string(RegionTr) || region == string(RegionUs) {
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
	}
}

func (s *CollectionsTestSuite) TestGetAllThemes() {
	client := newTestClient(s.Config)

	ctx := context.Background()

	resp, err := client.GetAllThemes(ctx, RegionTr, LocaleTr)
	s.Require().NoError(err)
	s.Require().NotEmpty(resp)

	hasStocks := false
	for _, theme := range resp {
		s.Require().NotEmpty(theme.ID)
		s.Require().NotEmpty(theme.Title)
		if theme.NumStocks > 0 {
			hasStocks = true
		}
	}
	s.Require().True(hasStocks)
}

func (s *CollectionsTestSuite) TestGetThemeDetail() {
	client := newTestClient(s.Config)

	ctx := context.Background()

	resp, err := client.GetThemeDetail(ctx, "64ff31e14ee6ea1024a76e73", RegionTr, LocaleTr)
	s.Require().NoError(err)
	s.Require().NotEmpty(resp)
	s.Require().NotEmpty(resp.ID)
	s.Require().NotEmpty(resp.Title)
	s.Require().NotEmpty(resp.Region)
	s.Require().Greater(resp.NumStocks, 0)

	hasValidRegion := false
	for _, region := range resp.Region {
		if region == string(RegionTr) || region == string(RegionUs) {
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
	}
}
