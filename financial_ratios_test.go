package laplace

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type FinancialRatiosTestSuite struct {
	*ClientTestSuite
}

func TestFinancialRatios(t *testing.T) {
	suite.Run(t, &FinancialRatiosTestSuite{
		NewClientTestSuite(),
	})
}

func (s *FinancialRatiosTestSuite) TestGetFinancialRatioComparison() {
	client := newTestClient(s.Config)

	ctx := context.Background()

	resp, err := client.GetFinancialRatioComparison(ctx, "TUPRS", RegionTr)
	require.NoError(s.T(), err)
	require.NotEmpty(s.T(), resp)
}

func (s *FinancialRatiosTestSuite) TestGetHistoricalRatios() {
	client := newTestClient(s.Config)

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
	client := newTestClient(s.Config)

	ctx := context.Background()

	resp, err := client.GetHistoricalRatiosDescriptions(ctx, LocaleTr, RegionTr)
	require.NoError(s.T(), err)
	require.NotEmpty(s.T(), resp)
}

func (s *FinancialRatiosTestSuite) TestGetHistoricalFinancialSheets() {
	client := newTestClient(s.Config)

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
