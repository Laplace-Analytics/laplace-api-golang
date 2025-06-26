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
	s.Require().NotEqual(primitive.NilObjectID, stock.SectorId)
	s.Require().NotEqual(primitive.NilObjectID, stock.IndustryId)
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
	s.Require().NotEqual(primitive.NilObjectID, resp.SectorId)
	s.Require().NotEqual(primitive.NilObjectID, resp.IndustryId)
	s.Require().Equal(true, resp.Active)
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
	s.Require().NotEqual(primitive.NilObjectID, resp.SectorId)
	s.Require().NotEqual(primitive.NilObjectID, resp.IndustryId)
	s.Require().Equal(true, resp.Active)
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
		s.Require().NotEmpty(restriction.Market)
		s.Require().NotEmpty(restriction.StartDate)
		s.Require().NotEmpty(restriction.EndDate)
	}
}

func (s *StocksTestSuite) TestGetAllRestrictions() {
	client := newTestClient(s.Config)

	ctx := context.Background()

	resp, err := client.GetAllRestrictions(ctx, RegionTr)
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
	s.Require().NotNil(resp)
}