package laplace

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type StocksTestSuite struct {
	*ClientTestSuite
}

func TestStocks(t *testing.T) {
	suite.Run(t, &StocksTestSuite{
		NewClientTestSuite(),
	})
}

func getAllAssetTypes() []AssetType {
	return []AssetType{
		AssetTypeStock,
		AssetTypeForex,
		AssetTypeIndex,
		AssetTypeEtf,
		AssetTypeCommodity,
		AssetTypeStockRights,
		AssetTypeFund,
	}
}

func (s *StocksTestSuite) TestGetAllStocks() {
	client := newTestClient(s.Config)

	ctx := context.Background()

	resp, err := client.GetAllStocks(ctx, RegionTr, 0, 0)
	s.Require().NoError(err)
	s.Require().NotNil(resp)
	s.Require().Greater(len(resp), 0)

	stock := resp[0]
	s.Require().NotEqual(primitive.NilObjectID, stock.ID)
	s.Require().NotEmpty(stock.Symbol)
	s.Require().NotEmpty(stock.Name)
	s.Require().Contains(getAllAssetTypes(), stock.AssetType)
	s.Require().NotEmpty(stock.SectorId)
	s.Require().NotEmpty(stock.IndustryId)
	s.Require().NotEmpty(stock.UpdatedDate)
}

func (s *StocksTestSuite) TestGetAllStocksPaginated() {
	client := newTestClient(s.Config)

	ctx := context.Background()

	resp, err := client.GetAllStocks(ctx, RegionTr, 0, 10)
	s.Require().NoError(err)
	s.Require().NotNil(resp)
	s.Require().Equal(10, len(resp))

	for _, stock := range resp {
		s.Require().NotEqual(primitive.NilObjectID, stock.ID)
		s.Require().NotEmpty(stock.Symbol)
		s.Require().NotEmpty(stock.Name)
		s.Require().Contains(getAllAssetTypes(), stock.AssetType)
	}
}

func (s *StocksTestSuite) TestGetStockDetailByID() {
	client := newTestClient(s.Config)

	ctx := context.Background()

	resp, err := client.GetStockDetailByID(ctx, "61dd0d960ec2114146343068", LocaleTr)
	s.Require().NoError(err)
	
	s.Require().NotEqual(primitive.NilObjectID, resp.ID)
	s.Require().Equal("TOASO", resp.Symbol)
	s.Require().NotEmpty(resp.Name)
	s.Require().Equal(string(RegionTr), resp.Region)
	s.Require().Equal(AssetClassEquity, resp.AssetClass)
	s.Require().Contains(string(AssetTypeStock), resp.AssetType)
	s.Require().NotEmpty(resp.SectorId)
	s.Require().NotEmpty(resp.IndustryId)
	s.Require().Equal(true, resp.Active)
	s.Require().NotEmpty(resp.Description)
	s.Require().NotEmpty(resp.ShortDescription)
	s.Require().NotNil(resp.LocalizedDescription)
	s.Require().NotNil(resp.LocalizedShortDescription)
}

func (s *StocksTestSuite) TestGetStockDetailBySymbol() {
	client := newTestClient(s.Config)

	ctx := context.Background()

	resp, err := client.GetStockDetailBySymbol(ctx, "TOASO", AssetClassEquity, RegionTr, LocaleTr)
	s.Require().NoError(err)

	s.Require().NotEqual(primitive.NilObjectID, resp.ID)
	s.Require().Equal("TOASO", resp.Symbol)
	s.Require().NotEmpty(resp.Name)
	s.Require().Equal(string(RegionTr), resp.Region)
	s.Require().Equal(AssetClassEquity, resp.AssetClass)
	s.Require().Contains(string(AssetTypeStock), resp.AssetType)
	s.Require().NotEmpty(resp.SectorId)
	s.Require().NotEmpty(resp.IndustryId)
	s.Require().Equal(true, resp.Active)
	s.Require().NotEmpty(resp.Description)
	s.Require().NotEmpty(resp.ShortDescription)
	s.Require().NotNil(resp.LocalizedDescription)
	s.Require().NotNil(resp.LocalizedShortDescription)
}

func (s *StocksTestSuite) TestGetHistoricalPrices() {
	client := newTestClient(s.Config)

	ctx := context.Background()

	resp, err := client.GetHistoricalPrices(ctx, []string{"TUPRS", "SASA"}, RegionTr, []HistoricalPricePeriod{HistoricalPricePeriodOneDay, HistoricalPricePeriodOneWeek, HistoricalPricePeriodOneMonth})
	s.Require().NoError(err)
	s.Require().NotNil(resp)
	s.Require().Greater(len(resp), 0)

	for _, priceGraph := range resp {
		s.Require().NotEmpty(priceGraph.Symbol)
		
		if len(priceGraph.OneDay) > 0 {
			point := priceGraph.OneDay[0]
			s.Require().Greater(point.Date, int64(0))
			s.Require().Greater(point.Close, 0.0)
			s.Require().Greater(point.High, 0.0)
			s.Require().Greater(point.Low, 0.0)
			s.Require().Greater(point.Open, 0.0)
		}
	}
}

func (s *StocksTestSuite) TestGetCustomHistoricalPrices() {
	client := newTestClient(s.Config)

	ctx := context.Background()

	resp, err := client.GetCustomHistoricalPrices(ctx, "TUPRS", RegionTr, "2024-01-01", "2024-03-01", HistoricalPriceIntervalOneDay, false)
	s.Require().NoError(err)
	s.Require().NotNil(resp)
	s.Require().Greater(len(resp), 0)

	for _, price := range resp {
		s.Require().Greater(price.Date, int64(0))
		s.Require().Greater(price.Close, 0.0)
		s.Require().Greater(price.High, 0.0)
		s.Require().Greater(price.Low, 0.0)
		s.Require().Greater(price.Open, 0.0)
	}
}

func (s *StocksTestSuite) TestGetStockRestrictions() {
	client := newTestClient(s.Config)

	ctx := context.Background()

	resp, err := client.GetStockRestrictions(ctx, "TUPRS", RegionTr)
	s.Require().NoError(err)
	s.Require().NotNil(resp)
	
	for _, restriction := range resp {
		s.Require().Greater(restriction.ID, 0)
		s.Require().NotEmpty(restriction.Title)
		s.Require().NotEmpty(restriction.Description)
	}
}

func (s *StocksTestSuite) TestGetAllRestrictions() {
	client := newTestClient(s.Config)

	ctx := context.Background()

	resp, err := client.GetAllRestrictions(ctx)
	s.Require().NoError(err)
	s.Require().NotNil(resp)

	for _, restriction := range resp {
		s.Require().Greater(restriction.ID, 0)
		s.Require().NotEmpty(restriction.Title)
		s.Require().NotEmpty(restriction.Description)
	}
}

func (s *StocksTestSuite) TestGetTickRules() {
	client := newTestClient(s.Config)

	ctx := context.Background()

	resp, err := client.GetTickRules(ctx, "TUPRS", RegionTr)
	s.Require().NoError(err)
	s.Require().NotZero(resp.BasePrice)
	s.Require().NotZero(resp.LowerPriceLimit)
	s.Require().NotZero(resp.UpperPriceLimit)
	s.Require().Greater(len(resp.Rules), 0)

	rule := resp.Rules[0]
	s.Require().GreaterOrEqual(rule.PriceFrom, 0.0)
	s.Require().Greater(rule.PriceTo, 0.0)
	s.Require().Greater(rule.TickSize, 0.0)
}

func (s *StocksTestSuite) TestGetStateOfAllMarkets() {
	client := newTestClient(s.Config)

	ctx := context.Background()

	resp, err := client.GetStateOfAllMarkets(ctx, RegionTr, 0, 10)
	s.Require().NoError(err)
	s.Require().NotNil(resp)
	s.Require().Greater(len(resp.Items), 0)

	state := resp.Items[0]
	s.Require().Greater(state.ID, 0)
	s.Require().NotEmpty(state.State)
	s.Require().NotEmpty(state.LastTimestamp)
}

func (s *StocksTestSuite) TestGetStateOfAllStocks() {
	client := newTestClient(s.Config)

	ctx := context.Background()

	resp, err := client.GetStateOfAllStocks(ctx, RegionTr, 0, 10)
	s.Require().NoError(err)
	s.Require().NotNil(resp)
	s.Require().Greater(len(resp.Items), 0)
}

func (s *StocksTestSuite) TestGetStateForStock() {
	client := newTestClient(s.Config)

	ctx := context.Background()

	resp, err := client.GetStateForStock(ctx, "TUPRS")
	s.Require().NoError(err)
	s.Require().Greater(resp.ID, 0)
	s.Require().NotEmpty(resp.State)
	s.Require().NotEmpty(resp.LastTimestamp)
}

func (s *StocksTestSuite) TestGetStateForMarket() {
	client := newTestClient(s.Config)

	ctx := context.Background()

	allMarkets, err := client.GetStateOfAllMarkets(ctx, RegionTr, 0, 10)
	s.Require().NoError(err)
	s.Require().Greater(len(allMarkets.Items), 0)

	marketSymbol := ""
	for _, m := range allMarkets.Items {
		if m.MarketSymbol != nil {
			marketSymbol = *m.MarketSymbol
			break
		}
	}
	s.Require().NotEmpty(marketSymbol)

	resp, err := client.GetStateForMarket(ctx, marketSymbol)
	s.Require().NoError(err)
	s.Require().Greater(resp.ID, 0)
	s.Require().NotEmpty(resp.State)
	s.Require().NotEmpty(resp.LastTimestamp)
}

func (s *StocksTestSuite) TestGetEarningsTranscriptList() {
	client := newTestClient(s.Config)

	ctx := context.Background()

	resp, err := client.GetEarningsTranscriptList(ctx, RegionUs, "AAPL")
	s.Require().NoError(err)
	s.Require().NotNil(resp)
	s.Require().Greater(len(resp), 0)

	transcript := resp[0]
	s.Require().NotEmpty(transcript.Symbol)
	s.Require().Greater(transcript.Year, 0)
	s.Require().Greater(transcript.Quarter, 0)
	s.Require().Greater(transcript.FiscalYear, 0)
}

func (s *StocksTestSuite) TestGetEarningsTranscriptWithSummary() {
	client := newTestClient(s.Config)

	ctx := context.Background()

	resp, err := client.GetEarningsTranscriptWithSummary(ctx, "AAPL", 2024, 1)
	s.Require().NoError(err)
	s.Require().NotNil(resp)
	s.Require().Equal("AAPL", resp.Symbol)
	s.Require().Greater(resp.Year, 0)
	s.Require().Greater(resp.Quarter, 0)
	s.Require().NotEmpty(resp.Date)
	s.Require().NotEmpty(resp.Content)
	s.Require().IsType(false, resp.HasSummary)
}

func (s *StocksTestSuite) TestGetStockChartImage() {
	client := newTestClient(s.Config)

	ctx := context.Background()

	resp, err := client.GetStockChartImage(ctx, GenerateChartImageRequest{
		Symbol: "TUPRS",
		Region: RegionTr,
	})
	s.Require().NoError(err)
	s.Require().NotNil(resp)
	s.Require().Greater(len(resp), 0)
}
