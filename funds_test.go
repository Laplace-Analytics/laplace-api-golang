package laplace

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
)

type FundsTestSuite struct {
	*ClientTestSuite
}

func TestFunds(t *testing.T) {
	suite.Run(t, &FundsTestSuite{
		NewClientTestSuite(),
	})
}

func getAllFundTypes() []FundType {
	return []FundType{
		FundTypeStockUmbrella,
		FundTypeVariableUmbrella,
		FundTypeParticipationUmbrella,
		FundTypeFlexibleUmbrella,
		FundTypeFundBasketUmbrella,
		FundTypeMoneyMarketUmbrella,
		FundTypePreciousMetalsUmbrella,
		FundTypeDebtInstrumentsUmbrella,
		FundTypeMixedUmbrella,
		FundTypeUnknown,
	}
}

func (s *FundsTestSuite) TestGetFunds() {
	client := newTestClient(s.Config)

	ctx := context.Background()

	resp, err := client.GetFunds(ctx, RegionTr, 0, 10)
	s.Require().NoError(err)
	s.Require().NotNil(resp)
	s.Require().Greater(len(resp), 0)

	fund := resp[0]
	s.Require().NotEmpty(fund.Symbol)
	s.Require().NotEmpty(fund.Name)
	s.Require().NotEmpty(fund.OwnerSymbol)
	s.Require().Equal(string(AssetTypeFund), fund.AssetType)
	s.Require().GreaterOrEqual(fund.ManagementFee, 0.0)
	s.Require().Contains(getAllFundTypes(), fund.FundType)
}

func (s *FundsTestSuite) TestGetFundStats() {
	client := newTestClient(s.Config)

	ctx := context.Background()

	resp, err := client.GetFundStats(ctx, "HIM", RegionTr)
	s.Require().NoError(err)
	s.Require().NotNil(resp)
}

func (s *FundsTestSuite) TestGetFundDistribution() {
	client := newTestClient(s.Config)

	ctx := context.Background()

	resp, err := client.GetFundDistribution(ctx, "HIM", RegionTr)
	s.Require().NoError(err)
	s.Require().NotNil(resp)
	s.Require().Greater(len(resp.Categories), 0)

	for _, category := range resp.Categories {
		s.Require().NotEmpty(category.Category)
		s.Require().Greater(category.Percentage, 0)
	}
}

func (s *FundsTestSuite) TestGetHistoricalFundPrices() {
	client := newTestClient(s.Config)

	ctx := context.Background()

	resp, err := client.GetHistoricalFundPrices(ctx, "HIM", RegionTr, HistoricalFundPricePeriodOneMonth)
	s.Require().NoError(err)
	s.Require().NotNil(resp)
	s.Require().Greater(len(resp), 0)

	price := resp[0]
	s.Require().NotEmpty(price.Date)
	s.Require().Greater(price.Price, 0.0)
	s.Require().Greater(price.Aum, 0.0)
	s.Require().Greater(price.ShareCount, 0.0)
	s.Require().Greater(price.InvestorCount, 0)
}