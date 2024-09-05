package client

import (
	"context"
	"testing"

	"finfree.co/laplace/utilities"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type StocksTestSuite struct {
	*utilities.ClientTestSuite
}

func TestStocks(t *testing.T) {
	suite.Run(t, &StocksTestSuite{
		utilities.NewClientTestSuite(),
	})
}

func (s *StocksTestSuite) TestGetAllStocks() {
	client := NewClient(s.Config, logrus.New())

	ctx := context.Background()

	resp, err := client.GetAllStocks(ctx, RegionTr)
	require.NoError(s.T(), err)

	require.NotEmpty(s.T(), resp)
}

func (s *StocksTestSuite) TestGetStockDetailByID() {
	client := NewClient(s.Config, logrus.New())

	ctx := context.Background()

	resp, err := client.GetStockDetailByID(ctx, "61dd0d6f0ec2114146342fd0", LocaleTr)
	require.NoError(s.T(), err)

	require.NotEmpty(s.T(), resp)
}

func (s *StocksTestSuite) TestGetStockDetailBySymbol() {
	client := NewClient(s.Config, logrus.New())

	ctx := context.Background()

	resp, err := client.GetStockDetailBySymbol(ctx, "TUPRS", AssetClassEquity, RegionTr, LocaleTr)
	require.NoError(s.T(), err)

	require.NotEmpty(s.T(), resp)
}

func (s *StocksTestSuite) TestGetHistoricalPrices() {
	client := NewClient(s.Config, logrus.New())

	ctx := context.Background()

	resp, err := client.GetHistoricalPrices(ctx, []string{"TUPRS", "SASA"}, RegionTr, []HistoricalPricePeriod{HistoricalPricePeriodOneDay, HistoricalPricePeriodOneWeek, HistoricalPricePeriodOneMonth})
	require.NoError(s.T(), err)

	require.NotEmpty(s.T(), resp)

	for _, price := range resp {
		require.NotEmpty(s.T(), price)
	}
}
