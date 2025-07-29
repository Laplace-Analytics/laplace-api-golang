package laplace

import (
	"context"
	"testing"

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
	s.Require().NoError(err)
	s.Require().NotNil(resp)
	s.Require().Greater(len(resp), 0)

	comparison := resp[0]
	s.Require().NotEmpty(comparison.MetricName)
	s.Require().NotEmpty(comparison.Data)

	data := comparison.Data[0]
	s.Require().NotEmpty(data.Slug)
}

func (s *FinancialRatiosTestSuite) TestGetHistoricalRatios() {
	client := newTestClient(s.Config)

	ctx := context.Background()

	resp, err := client.GetHistoricalRatios(ctx, "TUPRS", []HistoricalRatiosKey{HistoricalRatiosKeyPERatio}, RegionTr)
	s.Require().NoError(err)
	s.Require().NotNil(resp)
	s.Require().Greater(len(resp), 0)

	ratio := resp[0]
	s.Require().NotEmpty(ratio.Items)
	s.Require().NotEmpty(ratio.Slug)
	s.Require().NotEmpty(ratio.Name)
	s.Require().NotEmpty(ratio.Format)
	s.Require().NotEmpty(ratio.Currency)

	item := ratio.Items[0]
	s.Require().NotEmpty(item.Period)
}

func (s *FinancialRatiosTestSuite) TestGetHistoricalRatiosDescriptions() {
	client := newTestClient(s.Config)

	ctx := context.Background()

	resp, err := client.GetHistoricalRatiosDescriptions(ctx, LocaleTr, RegionTr)
	s.Require().NoError(err)
	s.Require().NotNil(resp)
	s.Require().Greater(len(resp), 0)

	desc := resp[0]
	s.Require().Greater(desc.ID, 0)
	s.Require().NotEmpty(desc.Slug)
	s.Require().NotEmpty(desc.Name)
	s.Require().NotEmpty(desc.Currency)
	s.Require().NotEmpty(desc.Format)
	s.Require().NotEmpty(desc.Description)
	s.Require().Equal(string(LocaleTr), desc.Locale)
}

func (s *FinancialRatiosTestSuite) TestGetHistoricalFinancialSheets() {
	client := newTestClient(s.Config)

	ctx := context.Background()

	resp, err := client.GetHistoricalFinancialSheets(ctx, "TUPRS", FinancialSheetDate{
		Year:  2022,
		Month: 1,
		Day:   1,
	}, FinancialSheetDate{
		Year:  2024,
		Month: 1,
		Day:   1,
	}, FinancialSheetBalanceSheet, FinancialSheetPeriodCumulative, CurrencyTRY, RegionTr)
	s.Require().NoError(err)
	s.Require().Greater(len(resp.Sheets), 0)

	sheet := resp.Sheets[0]
	s.Require().NotEmpty(sheet.Period)
	s.Require().Greater(len(sheet.Items), 0)

	item := sheet.Items[0]
	s.Require().NotEmpty(item.Description)
}
