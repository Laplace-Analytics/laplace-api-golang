package laplace

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
)

type BrokerSort string

const (
	BrokerSortNetBuy  BrokerSort = "netBuy"
	BrokerSortNetSell BrokerSort = "netSell"
	BrokerSortVolume  BrokerSort = "volume"
)

type Broker struct {
	ID       int    `json:"id"`
	Symbol   string `json:"symbol"`
	Name     string `json:"name"`
	LongName string `json:"longName"`
	Logo     string `json:"logo"`
}

type BrokerStock struct {
	Symbol     string     `json:"symbol"`
	Name       string     `json:"name"`
	ID         string     `json:"id"`
	AssetType  AssetType  `json:"assetType"`
	AssetClass AssetClass `json:"assetClass"`
	Region     Region     `json:"region"`
}

type BaseBrokerStats struct {
	TotalBuyAmount   float64 `json:"totalBuyAmount"`
	TotalSellAmount  float64 `json:"totalSellAmount"`
	NetAmount        float64 `json:"netAmount"`
	TotalBuyVolume   int64   `json:"totalBuyVolume"`
	TotalSellVolume  int64   `json:"totalSellVolume"`
	TotalVolume      int64   `json:"totalVolume"`
	TotalAmount      float64 `json:"totalAmount"`
}

type BrokerStats struct {
	BaseBrokerStats
	Broker Broker `json:"broker"`
}

type MarketBrokersResponse struct {
	RecordCount int             `json:"recordCount"`
	TotalStats  BaseBrokerStats `json:"totalStats"`
	Items       []BrokerStats   `json:"items"`
}

type TopBrokersResponse struct {
	TopStats  BaseBrokerStats `json:"topStats"`
	RestStats BaseBrokerStats `json:"restStats"`
	TopItems  []BrokerStats   `json:"topItems"`
}

type StockBrokerStats struct {
	BaseBrokerStats
	AverageCost float64 `json:"averageCost"`
	Broker      Broker  `json:"broker"`
}

type StockOverallStats struct {
	BaseBrokerStats
	AverageCost float64 `json:"averageCost"`
}

type StockBrokersResponse struct {
	RecordCount int                 `json:"recordCount"`
	TotalStats  StockOverallStats   `json:"totalStats"`
	Items       []StockBrokerStats  `json:"items"`
}

type TopStockBrokersResponse struct {
	TopStats  StockOverallStats  `json:"topStats"`
	RestStats StockOverallStats  `json:"restStats"`
	TopItems  []StockBrokerStats `json:"topItems"`
}

type BrokerStockStats struct {
	BaseBrokerStats
	Stock BrokerStock `json:"stock"`
}

type TopStocksForBrokerResponse struct {
	TopStats  BaseBrokerStats    `json:"topStats"`
	RestStats BaseBrokerStats    `json:"restStats"`
	TopItems  []BrokerStockStats `json:"topItems"`
}

func (c *Client) GetMarketBrokers(ctx context.Context, region Region, fromDate, toDate string, sortBy BrokerSort, page, size int) (MarketBrokersResponse, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/api/v1/brokers/market", c.baseUrl), nil)
	if err != nil {
		return MarketBrokersResponse{}, err
	}

	q := req.URL.Query()
	q.Add("region", string(region))
	q.Add("fromDate", fromDate)
	q.Add("toDate", toDate)
	q.Add("sortBy", string(sortBy))
	q.Add("page", strconv.Itoa(page))
	q.Add("size", strconv.Itoa(size))
	req.URL.RawQuery = q.Encode()

	resp, err := sendRequest[MarketBrokersResponse](ctx, c, req)
	if err != nil {
		return MarketBrokersResponse{}, err
	}

	return resp, nil
}

func (c *Client) GetTopMarketBrokers(ctx context.Context, region Region, fromDate, toDate string, sortBy BrokerSort, top int) (TopBrokersResponse, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/api/v1/brokers/market/top", c.baseUrl), nil)
	if err != nil {
		return TopBrokersResponse{}, err
	}

	q := req.URL.Query()
	q.Add("region", string(region))
	q.Add("fromDate", fromDate)
	q.Add("toDate", toDate)
	q.Add("sortBy", string(sortBy))
	q.Add("top", strconv.Itoa(top))
	req.URL.RawQuery = q.Encode()

	resp, err := sendRequest[TopBrokersResponse](ctx, c, req)
	if err != nil {
		return TopBrokersResponse{}, err
	}

	return resp, nil
}

func (c *Client) GetStockBrokers(ctx context.Context, region Region, fromDate, toDate string, sortBy BrokerSort, symbol string, page, size int) (StockBrokersResponse, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/api/v1/brokers/stock", c.baseUrl), nil)
	if err != nil {
		return StockBrokersResponse{}, err
	}

	q := req.URL.Query()
	q.Add("region", string(region))
	q.Add("fromDate", fromDate)
	q.Add("toDate", toDate)
	q.Add("sortBy", string(sortBy))
	q.Add("symbol", symbol)
	q.Add("page", strconv.Itoa(page))
	q.Add("size", strconv.Itoa(size))
	req.URL.RawQuery = q.Encode()

	resp, err := sendRequest[StockBrokersResponse](ctx, c, req)
	if err != nil {
		return StockBrokersResponse{}, err
	}

	return resp, nil
}

func (c *Client) GetTopStockBrokers(ctx context.Context, region Region, fromDate, toDate string, sortBy BrokerSort, symbol string, top int) (TopStockBrokersResponse, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/api/v1/brokers/stock/top", c.baseUrl), nil)
	if err != nil {
		return TopStockBrokersResponse{}, err
	}

	q := req.URL.Query()
	q.Add("region", string(region))
	q.Add("fromDate", fromDate)
	q.Add("toDate", toDate)
	q.Add("sortBy", string(sortBy))
	q.Add("symbol", symbol)
	q.Add("top", strconv.Itoa(top))
	req.URL.RawQuery = q.Encode()

	resp, err := sendRequest[TopStockBrokersResponse](ctx, c, req)
	if err != nil {
		return TopStockBrokersResponse{}, err
	}

	return resp, nil
}

func (c *Client) GetTopBrokersForBroker(ctx context.Context, region Region, fromDate, toDate string, sortBy BrokerSort, brokerSymbol string, top int) (TopBrokersResponse, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/api/v1/brokers/top", c.baseUrl), nil)
	if err != nil {
		return TopBrokersResponse{}, err
	}

	q := req.URL.Query()
	q.Add("region", string(region))
	q.Add("fromDate", fromDate)
	q.Add("toDate", toDate)
	q.Add("sortBy", string(sortBy))
	q.Add("brokerSymbol", brokerSymbol)
	q.Add("top", strconv.Itoa(top))
	req.URL.RawQuery = q.Encode()

	resp, err := sendRequest[TopBrokersResponse](ctx, c, req)
	if err != nil {
		return TopBrokersResponse{}, err
	}

	return resp, nil
}

func (c *Client) GetTopStocksForBroker(ctx context.Context, region Region, fromDate, toDate string, sortBy BrokerSort, brokerSymbol string, top int) (TopStocksForBrokerResponse, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/api/v1/brokers/top", c.baseUrl), nil)
	if err != nil {
		return TopStocksForBrokerResponse{}, err
	}

	q := req.URL.Query()
	q.Add("region", string(region))
	q.Add("fromDate", fromDate)
	q.Add("toDate", toDate)
	q.Add("sortBy", string(sortBy))
	q.Add("brokerSymbol", brokerSymbol)
	q.Add("top", strconv.Itoa(top))
	req.URL.RawQuery = q.Encode()

	resp, err := sendRequest[TopStocksForBrokerResponse](ctx, c, req)
	if err != nil {
		return TopStocksForBrokerResponse{}, err
	}

	return resp, nil
}