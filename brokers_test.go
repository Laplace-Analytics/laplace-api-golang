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

func (s *BrokerTestSuite) TestGetMarketBrokers() {
	client := newTestClient(s.Config)

	ctx := context.Background()

	resp, err := client.GetMarketBrokers(ctx, RegionTr, "2025-01-28", "2025-05-28", BrokerSortVolume, 0, 5)
	s.Require().NoError(err)
	s.Require().NotNil(resp)

	s.Require().Greater(resp.RecordCount, 0)
	s.Require().NotEmpty(resp.Items)

	if len(resp.Items) > 0 {
		item := resp.Items[0]
		s.Require().NotEmpty(item.Broker.Symbol)
		s.Require().NotEmpty(item.Broker.Name)
		s.Require().Greater(item.Broker.ID, 0)
		s.Require().NotEmpty(item.Broker.LongName)
	}

	s.Require().Greater(resp.TotalStats.TotalBuyVolume, int64(0))
	s.Require().Greater(resp.TotalStats.TotalSellVolume, int64(0))
	s.Require().Greater(resp.TotalStats.TotalVolume, int64(0))
	s.Require().Greater(resp.TotalStats.TotalBuyAmount, 0.0)
	s.Require().Greater(resp.TotalStats.TotalSellAmount, 0.0)
	s.Require().Greater(resp.TotalStats.TotalAmount, 0.0)
}

func (s *BrokerTestSuite) TestGetTopMarketBrokers() {
	client := newTestClient(s.Config)

	ctx := context.Background()

	resp, err := client.GetTopMarketBrokers(ctx, RegionTr, "2025-01-28", "2025-05-28", BrokerSortVolume, 5)
	s.Require().NoError(err)
	s.Require().NotNil(resp)

	s.Require().NotEmpty(resp.TopItems)

	for _, item := range resp.TopItems {
		s.Require().Greater(item.TotalBuyVolume, int64(0))
		s.Require().Greater(item.TotalSellVolume, int64(0))
		s.Require().Greater(item.TotalVolume, int64(0))
		s.Require().Greater(item.TotalBuyAmount, 0.0)
		s.Require().Greater(item.TotalSellAmount, 0.0)
		s.Require().Greater(item.TotalAmount, 0.0)
		s.Require().NotEmpty(item.Broker.Symbol)
		s.Require().NotEmpty(item.Broker.LongName)
		s.Require().Greater(item.Broker.ID, 0)
	}
}

func (s *BrokerTestSuite) TestGetStockBrokers() {
	client := newTestClient(s.Config)

	ctx := context.Background()

	resp, err := client.GetStockBrokers(ctx, RegionTr, "2025-01-28", "2025-05-28", BrokerSortVolume, "SASA", 0, 5)
	s.Require().NoError(err)
	s.Require().NotNil(resp)

	s.Require().Greater(resp.RecordCount, 0)
	s.Require().NotEmpty(resp.Items)

	s.Require().Greater(resp.TotalStats.TotalBuyVolume, int64(0))
	s.Require().Greater(resp.TotalStats.TotalSellVolume, int64(0))
	s.Require().Greater(resp.TotalStats.TotalVolume, int64(0))
	s.Require().Greater(resp.TotalStats.TotalBuyAmount, 0.0)
	s.Require().Greater(resp.TotalStats.TotalSellAmount, 0.0)
	s.Require().Greater(resp.TotalStats.TotalAmount, 0.0)
	s.Require().Greater(resp.TotalStats.AverageCost, 0.0)

	if len(resp.Items) > 0 {
		item := resp.Items[0]
		s.Require().NotEmpty(item.Broker.Symbol)
		s.Require().Greater(item.Broker.ID, 0)
		s.Require().NotEmpty(item.Broker.Name)
		s.Require().NotEmpty(item.Broker.LongName)

		s.Require().Greater(item.TotalVolume, int64(0))
		s.Require().Greater(item.TotalAmount, 0.0)
		s.Require().Greater(item.AverageCost, 0.0)
	}
}

func (s *BrokerTestSuite) TestGetTopStockBrokers() {
	client := newTestClient(s.Config)

	ctx := context.Background()

	resp, err := client.GetTopStockBrokers(ctx, RegionTr, "2025-01-28", "2025-05-28", BrokerSortVolume, "TUPRS", 5)
	s.Require().NoError(err)
	s.Require().NotNil(resp)

	s.Require().NotEmpty(resp.TopItems)

	for _, item := range resp.TopItems {
		s.Require().Greater(item.TotalBuyVolume, int64(0))
		s.Require().Greater(item.TotalSellVolume, int64(0))
		s.Require().Greater(item.TotalVolume, int64(0))
		s.Require().Greater(item.TotalBuyAmount, 0.0)
		s.Require().Greater(item.TotalSellAmount, 0.0)
		s.Require().Greater(item.TotalAmount, 0.0)
		s.Require().Greater(item.AverageCost, 0.0)
		s.Require().NotEmpty(item.Broker.Symbol)
		s.Require().NotEmpty(item.Broker.Name)
		s.Require().NotEmpty(item.Broker.LongName)
		s.Require().Greater(item.Broker.ID, 0)
	}
}

func (s *BrokerTestSuite) TestGetTopStocksForBroker() {
	client := newTestClient(s.Config)

	ctx := context.Background()

	resp, err := client.GetTopStocksForBroker(ctx, RegionTr, "2025-01-28", "2025-05-28", BrokerSortVolume, "BIYKR", 5)
	s.Require().NoError(err)
	s.Require().NotNil(resp)

	s.Require().NotEmpty(resp.TopItems)

	for _, item := range resp.TopItems {
		s.Require().NotEmpty(item.Stock.Symbol)
		s.Require().NotEmpty(item.Stock.Name)
		s.Require().NotEmpty(item.Stock.ID)
		s.Require().NotEmpty(string(item.Stock.AssetType))
		s.Require().NotEmpty(string(item.Stock.AssetClass))
		s.Require().Greater(item.TotalBuyVolume, int64(0))
		s.Require().Greater(item.TotalSellVolume, int64(0))
		s.Require().Greater(item.TotalVolume, int64(0))
		s.Require().Greater(item.TotalBuyAmount, 0.0)
		s.Require().Greater(item.TotalSellAmount, 0.0)
		s.Require().Greater(item.TotalAmount, 0.0)
	}
}