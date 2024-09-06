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

	resp, err := client.GetStockStats(ctx, []string{"TUPRS"}, []StockStatsKey{
		StockStatsPreviousClose,
		StockStatsYtdReturn,
		StockStatsYearlyReturn,
		StockStatsMarketCap,
		StockStatsFK,
		StockStatsPDDD,
		StockStatsYearLow,
		StockStatsYearHigh,
		StockStats3YearReturn,
		StockStats5YearReturn,
		StockStatsLatestPrice,
	}, RegionTr)
	require.NoError(s.T(), err)
	require.NotEmpty(s.T(), resp)
}

func (s *FinancialFundamentalsTestSuite) TestGetTopMovers() {
	client := NewClient(s.Config, logrus.New())

	ctx := context.Background()

	resp, err := client.GetTopMovers(ctx, RegionTr)
	require.NoError(s.T(), err)
	require.NotEmpty(s.T(), resp)
}
