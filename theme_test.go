package laplace

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
)

type ThemeTestSuite struct {
	*ClientTestSuite
}

func TestTheme(t *testing.T) {
	suite.Run(t, &ThemeTestSuite{
		NewClientTestSuite(),
	})
}

func (s *ThemeTestSuite) TestGetAllThemes() {
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

func (s *ThemeTestSuite) TestGetThemeDetail() {
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
