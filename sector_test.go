package laplace

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
)

type SectorTestSuite struct {
	*ClientTestSuite
}

func TestSector(t *testing.T) {
	suite.Run(t, &SectorTestSuite{
		NewClientTestSuite(),
	})
}

func (s *SectorTestSuite) TestGetAllSectors() {
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

func (s *SectorTestSuite) TestGetSectorDetails() {
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
