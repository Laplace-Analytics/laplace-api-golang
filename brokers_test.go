package laplace

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
)

type BrokerTestSuite struct {
	*ClientTestSuite
}

func TestBroker(t *testing.T) {
	suite.Run(t, &BrokerTestSuite{
		NewClientTestSuite(),
	})
}

func (s *BrokerTestSuite) TestGetBrokers() {
	client := newTestClient(s.Config)

	ctx := context.Background()

	resp, err := client.GetBrokers(ctx, RegionTr, 0, 10)
	s.Require().NoError(err)
	s.Require().NotNil(resp)

	s.Require().Greater(resp.RecordCount, 0)
	s.Require().NotEmpty(resp.Items)

	if len(resp.Items) > 0 {
		broker := resp.Items[0]
		s.Require().NotEmpty(broker.Symbol)
		s.Require().NotEmpty(broker.Name)
		s.Require().Greater(broker.ID, 0)
		s.Require().NotEmpty(broker.LongName)
	}
}

func (s *BrokerTestSuite) TestGetMarketBrokers() {
	client := newTestClient(s.Config)

	ctx := context.Background()

	resp, err := client.GetMarketBrokers(ctx, RegionTr, BrokerSortTotalVolume, BrokerSortDirectionDesc, "2025-01-28", "2025-05-28", 0, 5)
	s.Require().NoError(err)
	s.Require().NotNil(resp)

	s.Require().Greater(resp.RecordCount, 0)
	s.Require().NotEmpty(resp.Items)

	s.Require().Greater(resp.TotalStats.TotalBuyVolume, 0.0)
	s.Require().Greater(resp.TotalStats.TotalSellVolume, 0.0)
	s.Require().Greater(resp.TotalStats.TotalVolume, 0.0)
	s.Require().Greater(resp.TotalStats.TotalBuyAmount, 0.0)
	s.Require().Greater(resp.TotalStats.TotalSellAmount, 0.0)
	s.Require().Greater(resp.TotalStats.TotalAmount, 0.0)

	if len(resp.Items) > 0 {
		item := resp.Items[0]
		s.Require().Greater(item.TotalBuyVolume, 0.0)
		s.Require().Greater(item.TotalSellVolume, 0.0)
		s.Require().Greater(item.TotalVolume, 0.0)
		s.Require().Greater(item.TotalBuyAmount, 0.0)
		s.Require().Greater(item.TotalSellAmount, 0.0)
		s.Require().Greater(item.TotalAmount, 0.0)
		
		s.Require().NotEmpty(item.Broker.Symbol)
		s.Require().NotEmpty(item.Broker.Name)
		s.Require().Greater(item.Broker.ID, 0)
		s.Require().NotEmpty(item.Broker.LongName)
	}
}

func (s *BrokerTestSuite) TestGetMarketStocks() {
	client := newTestClient(s.Config)

	ctx := context.Background()

	resp, err := client.GetMarketStocks(ctx, RegionTr, BrokerSortTotalVolume, BrokerSortDirectionDesc, "2025-01-28", "2025-05-28", 0, 5)
	s.Require().NoError(err)
	s.Require().NotNil(resp)

	s.Require().Greater(resp.RecordCount, 0)
	s.Require().NotEmpty(resp.Items)

	s.Require().Greater(resp.TotalStats.TotalBuyVolume, 0.0)
	s.Require().Greater(resp.TotalStats.TotalSellVolume, 0.0)
	s.Require().Greater(resp.TotalStats.TotalVolume, 0.0)
	s.Require().Greater(resp.TotalStats.TotalBuyAmount, 0.0)
	s.Require().Greater(resp.TotalStats.TotalSellAmount, 0.0)
	s.Require().Greater(resp.TotalStats.TotalAmount, 0.0)
	s.Require().Greater(resp.TotalStats.AverageCost, 0.0)

	if len(resp.Items) > 0 {
		item := resp.Items[0]
		s.Require().Greater(item.TotalBuyVolume, 0.0)
		s.Require().Greater(item.TotalSellVolume, 0.0)
		s.Require().Greater(item.TotalVolume, 0.0)
		s.Require().Greater(item.TotalBuyAmount, 0.0)
		s.Require().Greater(item.TotalSellAmount, 0.0)
		s.Require().Greater(item.TotalAmount, 0.0)
		s.Require().Greater(item.AverageCost, 0.0)
		if item.Stock != nil {
			s.Require().NotEmpty(item.Stock.Symbol)
			s.Require().NotEmpty(item.Stock.Name)
			s.Require().NotEmpty(item.Stock.AssetId)
			s.Require().NotEmpty(item.Stock.AssetType)
			s.Require().NotEmpty(item.Stock.AssetClass)
		}
	}
}

func (s *BrokerTestSuite) TestGetBrokersByStock() {
	client := newTestClient(s.Config)

	ctx := context.Background()

	resp, err := client.GetBrokersByStock(ctx, "SASA", RegionTr, BrokerSortTotalVolume, BrokerSortDirectionDesc, "2025-01-28", "2025-05-28", 0, 5)
	s.Require().NoError(err)
	s.Require().NotNil(resp)

	s.Require().Greater(resp.RecordCount, 0)
	s.Require().NotEmpty(resp.Items)

	s.Require().Greater(resp.TotalStats.TotalBuyVolume, 0.0)
	s.Require().Greater(resp.TotalStats.TotalSellVolume, 0.0)
	s.Require().Greater(resp.TotalStats.TotalVolume, 0.0)
	s.Require().Greater(resp.TotalStats.TotalBuyAmount, 0.0)
	s.Require().Greater(resp.TotalStats.TotalSellAmount, 0.0)
	s.Require().Greater(resp.TotalStats.TotalAmount, 0.0)
	s.Require().Greater(resp.TotalStats.AverageCost, 0.0)

	if len(resp.Items) > 0 {
		item := resp.Items[0]
		s.Require().Greater(item.TotalBuyVolume, 0.0)
		s.Require().Greater(item.TotalSellVolume, 0.0)
		s.Require().Greater(item.TotalVolume, 0.0)
		s.Require().Greater(item.TotalBuyAmount, 0.0)
		s.Require().Greater(item.TotalSellAmount, 0.0)
		s.Require().Greater(item.TotalAmount, 0.0)
		s.Require().Greater(item.AverageCost, 0.0)

		s.Require().NotEmpty(item.Broker.Symbol)
		s.Require().Greater(item.Broker.ID, 0)
		s.Require().NotEmpty(item.Broker.Name)
		s.Require().NotEmpty(item.Broker.LongName)
	}
}

func (s *BrokerTestSuite) TestGetStocksByBroker() {
	client := newTestClient(s.Config)

	ctx := context.Background()

	resp, err := client.GetStocksByBroker(ctx, "BIYKR", RegionTr, BrokerSortTotalVolume, BrokerSortDirectionDesc, "2025-01-28", "2025-05-28", 0, 5)
	s.Require().NoError(err)
	s.Require().NotNil(resp)

	s.Require().Greater(resp.RecordCount, 0)
	s.Require().NotEmpty(resp.Items)

	s.Require().Greater(resp.TotalStats.TotalBuyVolume, 0.0)
	s.Require().Greater(resp.TotalStats.TotalSellVolume, 0.0)
	s.Require().Greater(resp.TotalStats.TotalVolume, 0.0)
	s.Require().Greater(resp.TotalStats.TotalBuyAmount, 0.0)
	s.Require().Greater(resp.TotalStats.TotalSellAmount, 0.0)
	s.Require().Greater(resp.TotalStats.TotalAmount, 0.0)

	if len(resp.Items) > 0 {
		item := resp.Items[0]
		if item.Stock != nil {
			s.Require().NotEmpty(item.Stock.Symbol)
			s.Require().NotEmpty(item.Stock.Name)
			s.Require().NotEmpty(item.Stock.AssetId)
			s.Require().NotEmpty(item.Stock.AssetType)
			s.Require().NotEmpty(item.Stock.AssetClass)
		}

		s.Require().Greater(item.TotalBuyVolume, 0.0)
		s.Require().Greater(item.TotalSellVolume, 0.0)
		s.Require().Greater(item.TotalVolume, 0.0)
		s.Require().Greater(item.TotalBuyAmount, 0.0)
		s.Require().Greater(item.TotalSellAmount, 0.0)
		s.Require().Greater(item.TotalAmount, 0.0)
	}
}

