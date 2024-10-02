package client

import (
	"context"
	"testing"

	"finfree.co/laplace/utilities"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type FinancialRatiosTestSuite struct {
	*utilities.ClientTestSuite
}

func TestFinancialRatios(t *testing.T) {
	suite.Run(t, &FinancialRatiosTestSuite{
		utilities.NewClientTestSuite(),
	})
}

func (s *FinancialRatiosTestSuite) TestGetFinancialRatioComparison() {
	client := NewClient(s.Config, logrus.New())

	ctx := context.Background()

	resp, err := client.GetFinancialRatioComparison(ctx, "TUPRS", RegionTr)
	require.NoError(s.T(), err)
	require.NotEmpty(s.T(), resp)
}

func (s *FinancialRatiosTestSuite) TestGetHistoricalRatios() {
	client := NewClient(s.Config, logrus.New())

	ctx := context.Background()

	resp, err := client.GetHistoricalRatios(ctx, "TUPRS", []HistoricalRatiosKey{HistoricalRatiosKeyPriceToEarningsRatio}, RegionTr)
	require.NoError(s.T(), err)
	require.NotEmpty(s.T(), resp)
	for _, formatting := range resp.Formatting {
		require.NotEmpty(s.T(), formatting)
		require.NotEmpty(s.T(), formatting.Name)
	}
}

func (s *FinancialRatiosTestSuite) TestGetHistoricalRatiosDescriptions() {
	client := NewClient(s.Config, logrus.New())

	ctx := context.Background()

	resp, err := client.GetHistoricalRatiosDescriptions(ctx, LocaleTr, RegionTr)
	require.NoError(s.T(), err)
	require.NotEmpty(s.T(), resp)
}

func (s *FinancialRatiosTestSuite) TestGetHistoricalFinancialSheets() {
	client := NewClient(s.Config, logrus.New())

	ctx := context.Background()

	resp, err := client.GetHistoricalFinancialSheets(ctx, "TUPRS", FinancialSheetDate{
		Year:  2022,
		Month: 1,
		Day:   1,
	}, FinancialSheetDate{
		Year:  2023,
		Month: 1,
		Day:   1,
	}, FinancialSheetBalanceSheet, FinancialSheetPeriodAnnual, CurrencyEUR, RegionTr)
	require.NoError(s.T(), err)
	require.NotEmpty(s.T(), resp)
}
