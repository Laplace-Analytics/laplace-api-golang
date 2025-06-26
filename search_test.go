package laplace

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	s.Require().NoError(err)
	s.Require().NotNil(resp)
	s.Require().Greater(len(resp.Stocks), 0)

	stock := resp.Stocks[0]
	s.Require().NotEqual(primitive.NilObjectID, stock.ID)
	s.Require().NotEmpty(stock.Name)
	s.Require().NotEmpty(stock.Symbol)
	s.Require().Equal(string(RegionTr), stock.Region)
	s.Require().NotEmpty(stock.AssetClass)
	s.Require().NotEmpty(stock.AssetType)
}

func (s *SearchTestSuite) TestSearchIndustry() {
	client := newTestClient(s.Config)

	ctx := context.Background()

	resp, err := client.Search(ctx, "Hava Taşımacılığı", []SearchType{SearchTypeIndustry}, RegionTr, LocaleTr, 0, PageSize10)
	s.Require().NoError(err)
	s.Require().NotNil(resp)
	s.Require().Greater(len(resp.Industries), 0)

	industry := resp.Industries[0]
	s.Require().NotEqual(primitive.NilObjectID, industry.ID)
	s.Require().NotEmpty(industry.Title)

	hasValidRegion := false
		for _, region := range industry.Region {
			if region == string(RegionTr) || region == string(RegionUs) {
				hasValidRegion = true
				break
			}
		}
		s.Require().True(hasValidRegion)
}

func (s *SearchTestSuite) TestSearchAllTypes() {
	client := newTestClient(s.Config)

	ctx := context.Background()

	resp, err := client.Search(ctx, "A", []SearchType{
		SearchTypeStock,
		SearchTypeIndustry,
		SearchTypeSector,
		SearchTypeCollection,
	}, RegionTr, LocaleTr, 0, PageSize10)
	s.Require().NoError(err)
	s.Require().NotNil(resp)

	hasResults := len(resp.Stocks) > 0 || 
		len(resp.Industries) > 0 || 
		len(resp.Sectors) > 0 || 
		len(resp.Collections) > 0
	s.Require().True(hasResults)

	for _, stock := range resp.Stocks {
		s.Require().NotEqual(primitive.NilObjectID, stock.ID)
		s.Require().NotEmpty(stock.Name)
		s.Require().NotEmpty(stock.Symbol)
		s.Require().Equal(string(RegionTr), stock.Region)
	}

	for _, collection := range resp.Collections {
		s.Require().NotEqual(primitive.NilObjectID, collection.ID)
		s.Require().NotEmpty(collection.Title)
		s.Require().NotEmpty(collection.Region)
	}

	for _, sector := range resp.Sectors {
		s.Require().NotEqual(primitive.NilObjectID, sector.ID)
		s.Require().NotEmpty(sector.Title)
		s.Require().NotEmpty(sector.Region)
	}

	for _, industry := range resp.Industries {
		s.Require().NotEqual(primitive.NilObjectID, industry.ID)
		s.Require().NotEmpty(industry.Title)
		s.Require().NotEmpty(industry.Region)
	}
}