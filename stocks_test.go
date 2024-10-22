package laplace

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type StocksTestSuite struct {
	*ClientTestSuite
}

func TestStocks(t *testing.T) {
	suite.Run(t, &StocksTestSuite{
		NewClientTestSuite(),
	})
}

func (s *StocksTestSuite) TestGetAllStocks() {
	client := newTestClient(s.Config)

	ctx := context.Background()

	resp, err := client.GetAllStocks(ctx, RegionTr)
	require.NoError(s.T(), err)

	require.NotEmpty(s.T(), resp)
}

func (s *StocksTestSuite) TestGetStockDetailByID() {
	client := newTestClient(s.Config)

	ctx := context.Background()

	resp, err := client.GetStockDetailByID(ctx, "61dd0d6f0ec2114146342fd0", LocaleTr)
	require.NoError(s.T(), err)

	require.NotEmpty(s.T(), resp)
}

func (s *StocksTestSuite) TestGetStockDetailBySymbol() {
	client := newTestClient(s.Config)

	ctx := context.Background()

	resp, err := client.GetStockDetailBySymbol(ctx, "TUPRS", AssetClassEquity, RegionTr, LocaleTr)
	require.NoError(s.T(), err)

	require.NotEmpty(s.T(), resp)
}

func (s *StocksTestSuite) TestGetHistoricalPrices() {
	client := newTestClient(s.Config)

	ctx := context.Background()

	resp, err := client.GetHistoricalPrices(ctx, []string{"TUPRS", "SASA"}, RegionTr, []HistoricalPricePeriod{HistoricalPricePeriodOneDay, HistoricalPricePeriodOneWeek, HistoricalPricePeriodOneMonth})
	require.NoError(s.T(), err)

	require.NotEmpty(s.T(), resp)

	for _, price := range resp {
		require.NotEmpty(s.T(), price)
	}
}

func (s *StocksTestSuite) TestGetCustomHistoricalPrices() {
	client := newTestClient(s.Config)

	ctx := context.Background()

	resp, err := client.GetCustomHistoricalPrices(ctx, "TUPRS", RegionTr, "2024-01-01", "2024-03-01", HistoricalPriceIntervalOneDay, false)
	require.NoError(s.T(), err)

	require.NotEmpty(s.T(), resp)

	for _, price := range resp {
		require.NotEmpty(s.T(), price)
	}

	resp, err = client.GetCustomHistoricalPrices(ctx, "TUPRS", RegionTr, "2024-01-01 10:00:00", "2024-01-02 10:00:00", HistoricalPriceIntervalOneHour, true)
	require.NoError(s.T(), err)

	require.NotEmpty(s.T(), resp)

	for _, price := range resp {
		require.NotEmpty(s.T(), price)
	}
}

func (s *StocksTestSuite) TestGetStockRestrictions() {
	client := newTestClient(s.Config)

	ctx := context.Background()

	resp, err := client.GetStockRestrictions(ctx, "TUPRS", RegionTr)
	require.NoError(s.T(), err)

	require.NotNil(s.T(), resp)
}
