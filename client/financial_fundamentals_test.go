package client

import (
	"context"
	"testing"

	"finfree.co/laplace/utilities"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type FinancialFundamentalsTestSuite struct {
	*utilities.ClientTestSuite
}

func TestFinancialFundamentals(t *testing.T) {
	suite.Run(t, &FinancialFundamentalsTestSuite{
		utilities.NewClientTestSuite(),
	})
}

func (s *FinancialFundamentalsTestSuite) TestGetStockDividends() {
	client := NewClient(s.Config, logrus.New())

	ctx := context.Background()

	resp, err := client.GetStockDividends(ctx, "TUPRS", RegionTr)
	require.NoError(s.T(), err)
	require.NotEmpty(s.T(), resp)
}

func (s *FinancialFundamentalsTestSuite) TestGetStockStats() {
	client := NewClient(s.Config, logrus.New())

	ctx := context.Background()

	statKeys := []StockStatsKey{
		StockStatsPreviousClose,
		StockStatsMarketCap,
		StockStatsFK,
		StockStatsPDDD,
		StockStatsYearLow,
		StockStatsYearHigh,
		StockStatsWeeklyReturn,
		StockStatsMonthlyReturn,
		StockStats3MonthReturn,
		StockStatsYtdReturn,
		StockStatsYearlyReturn,
		StockStats3YearReturn,
		StockStats5YearReturn,
		StockStatsLatestPrice,
	}

	resp, err := client.GetStockStats(ctx, []string{"TUPRS"}, statKeys, RegionTr)
	require.NoError(s.T(), err)
	require.NotEmpty(s.T(), resp)
	currentStockStats := resp[0]
	require.NotEmpty(s.T(), currentStockStats)
	require.Equal(s.T(), "TUPRS", currentStockStats.Symbol)
	require.Greater(s.T(), currentStockStats.PreviousClose, 0.0)
	require.Greater(s.T(), currentStockStats.MarketCap, 0.0)
	require.NotEqual(s.T(), currentStockStats.PeRatio, 0.0)
	require.NotEqual(s.T(), currentStockStats.PbRatio, 0.0)
	require.Greater(s.T(), currentStockStats.YearLow, 0.0)
	require.Greater(s.T(), currentStockStats.YearHigh, 0.0)
	require.NotEqual(s.T(), currentStockStats.WeeklyReturn, 0.0)
	require.NotEqual(s.T(), currentStockStats.MonthlyReturn, 0.0)
	require.NotEqual(s.T(), currentStockStats.ThreeMonthReturn, 0.0)
	require.NotEqual(s.T(), currentStockStats.YtdReturn, 0.0)
	require.NotEqual(s.T(), currentStockStats.YearlyReturn, 0.0)
	require.NotEqual(s.T(), currentStockStats.ThreeYear, 0.0)
	require.NotEqual(s.T(), currentStockStats.FiveYear, 0.0)
}

func (s *FinancialFundamentalsTestSuite) TestGetTopMovers() {
	client := NewClient(s.Config, logrus.New())

	ctx := context.Background()

	resp, err := client.GetTopMovers(ctx, RegionTr)
	require.NoError(s.T(), err)
	require.NotEmpty(s.T(), resp)
}
