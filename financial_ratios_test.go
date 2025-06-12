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

	resp, err := client.GetFinancialRatioComparison(ctx, "TUPRS", RegionTr, PeerTypeSector)
	require.NoError(s.T(), err)
	require.NotEmpty(s.T(), resp)
}

func (s *FinancialRatiosTestSuite) TestGetHistoricalRatios() {
	client := newTestClient(s.Config)

	ctx := context.Background()

	resp, err := client.GetHistoricalRatios(ctx, "TUPRS", []HistoricalRatiosKey{HistoricalRatiosKeyPERatio}, RegionTr)
	require.NoError(s.T(), err)
	require.NotEmpty(s.T(), resp)

	for _, item := range resp {
		require.NotEmpty(s.T(), item.Items)
		require.NotEmpty(s.T(), item.FinalValue)
		require.NotEmpty(s.T(), item.ThreeYearGrowth)
		require.NotEmpty(s.T(), item.YearGrowth)
		require.NotEmpty(s.T(), item.FinalSectorValue)
		require.NotEmpty(s.T(), item.Slug)
		require.NotEmpty(s.T(), item.Currency)
		require.NotEmpty(s.T(), item.Format)
		require.NotEmpty(s.T(), item.Name)
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
	}, FinancialSheetBalanceSheet, FinancialSheetPeriodAnnual, CurrencyTRY, RegionTr)
	require.NoError(s.T(), err)
	require.NotEmpty(s.T(), resp)
}
