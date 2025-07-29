package laplace

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
)

type FinancialFundamentalsTestSuite struct {
	*ClientTestSuite
}

func TestFinancialFundamentals(t *testing.T) {
	suite.Run(t, &FinancialFundamentalsTestSuite{
		NewClientTestSuite(),
	})
}

func (s *FinancialFundamentalsTestSuite) TestGetStockDividends() {
	client := newTestClient(s.Config)

	ctx := context.Background()

	resp, err := client.GetStockDividends(ctx, "AKBNK", RegionTr)
	s.Require().NoError(err)
	s.Require().NotNil(resp)
	s.Require().Greater(len(resp), 0)

	dividend := resp[0]
	s.Require().NotEmpty(dividend.Date)
	s.Require().Greater(dividend.GrossAmount, 0.0)
	s.Require().Greater(dividend.GrossRatio, 0.0)
	s.Require().Greater(dividend.PriceThen, 0.0)
}

func (s *FinancialFundamentalsTestSuite) TestGetStockStats() {
	client := newTestClient(s.Config)

	ctx := context.Background()

	resp, err := client.GetStockStats(ctx, []string{"TUPRS"}, RegionTr)
	s.Require().NoError(err)
	s.Require().NotEmpty(resp)

	currentStockStats := resp[0]
	s.Require().Equal("TUPRS", currentStockStats.Symbol)
	s.Require().Greater(currentStockStats.PreviousClose, 0.0)
	s.Require().Greater(currentStockStats.MarketCap, 0.0)
	s.Require().Greater(currentStockStats.YearLow, 0.0)
	s.Require().Greater(currentStockStats.YearHigh, 0.0)
	s.Require().NotEqual(currentStockStats.LowerPriceLimit, 0.0)
	s.Require().NotEqual(currentStockStats.UpperPriceLimit, 0.0)
	s.Require().NotEqual(currentStockStats.DayOpen, 0.0)
}

func (s *FinancialFundamentalsTestSuite) TestGetTopMovers() {
	client := newTestClient(s.Config)

	ctx := context.Background()

	assetClass := AssetClassEquity
	assetType := AssetTypeStock

	respGainers, err := client.GetTopMovers(ctx, TopMoversDirectionGainers, assetClass, assetType, 0, 10, RegionTr)
	s.Require().NoError(err)
	s.Require().NotNil(respGainers)
	s.Require().Greater(len(respGainers), 0)

	gainer := respGainers[0]
	s.Require().NotEmpty(gainer.Symbol)
	s.Require().Greater(gainer.Change, 0.0)
	s.Require().Equal(assetClass, gainer.AssetClass)
	s.Require().Equal(assetType, gainer.AssetType)

	respLosers, err := client.GetTopMovers(ctx, TopMoversDirectionLosers, assetClass, assetType, 0, 10, RegionTr)
	s.Require().NoError(err)
	s.Require().NotNil(respLosers)
	s.Require().Greater(len(respLosers), 0)

	loser := respLosers[0]
	s.Require().NotEmpty(loser.Symbol)
	s.Require().Less(loser.Change, 0.0)
	s.Require().Equal(assetClass, gainer.AssetClass)
	s.Require().Equal(assetType, gainer.AssetType)
}
