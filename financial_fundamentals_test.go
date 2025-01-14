package laplace

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
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

	resp, err := client.GetStockDividends(ctx, "TUPRS", RegionTr)
	require.NoError(s.T(), err)
	require.NotEmpty(s.T(), resp)
}

func (s *FinancialFundamentalsTestSuite) TestGetStockStats() {
	client := newTestClient(s.Config)

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

func (s *FinancialFundamentalsTestSuite) TestGetStockStatsV2() {
	client := newTestClient(s.Config)

	ctx := context.Background()

	resp, err := client.GetStockStatsV2(ctx, []string{"TUPRS"}, RegionTr)
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
	require.NotEqual(s.T(), currentStockStats.LowerPriceLimit, 0.0)
	require.NotEqual(s.T(), currentStockStats.UpperPriceLimit, 0.0)
}

func (s *FinancialFundamentalsTestSuite) TestGetTopMovers() {
	client := newTestClient(s.Config)

	ctx := context.Background()

	resp, err := client.GetTopMovers(ctx, RegionTr)
	require.NoError(s.T(), err)
	require.NotEmpty(s.T(), resp)
}
