package laplace

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
)

type IndustryTestSuite struct {
	*ClientTestSuite
}

func TestIndustry(t *testing.T) {
	suite.Run(t, &IndustryTestSuite{
		NewClientTestSuite(),
	})
}

func (s *IndustryTestSuite) TestGetAllIndustries() {
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

func (s *IndustryTestSuite) TestGetIndustryDetails() {
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
	}
}
