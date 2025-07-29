package laplace

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
)

type BrokerSort string

const (
	BrokerSortNetAmount       BrokerSort = "netAmount"
	BrokerSortTotalAmount     BrokerSort = "totalAmount"
	BrokerSortTotalVolume     BrokerSort = "totalVolume"
	BrokerSortTotalBuyAmount  BrokerSort = "totalBuyAmount"
	BrokerSortTotalBuyVolume  BrokerSort = "totalBuyVolume"
	BrokerSortTotalSellAmount BrokerSort = "totalSellAmount"
	BrokerSortTotalSellVolume BrokerSort = "totalSellVolume"
)

type BrokerSortDirection string

const (
	BrokerSortDirectionDesc BrokerSortDirection = "desc"
	BrokerSortDirectionAsc  BrokerSortDirection = "asc"
)

type Broker struct {
	ID       int    `json:"id"`
	Symbol   string `json:"symbol"`
	Name     string `json:"name"`
	LongName string `json:"longName"`
	Logo     string `json:"logo"`
}

type BrokerStock struct {
	Symbol     string `json:"symbol"`
	Name       string `json:"name"`
	AssetId    string `json:"id"`
	AssetType  string `json:"assetType"`
	AssetClass string `json:"assetClass"`
	LogoUrl    string `json:"logoUrl,omitempty"`
	Exchange   string `json:"exchange,omitempty"`
}

type BrokerStats struct {
	TotalBuyAmount  float64 `json:"totalBuyAmount"`
	TotalSellAmount float64 `json:"totalSellAmount"`
	NetAmount       float64 `json:"netAmount"`
	TotalBuyVolume  float64 `json:"totalBuyVolume"`
	TotalSellVolume float64 `json:"totalSellVolume"`
	TotalVolume     float64 `json:"totalVolume"`
	TotalAmount     float64 `json:"totalAmount"`
	AverageCost     float64 `json:"averageCost,omitempty"`
}

type BrokerListResponse struct {
	PaginatedResponse[*BrokerResponseItem]
	TotalStats BrokerStats `json:"totalStats"`
}

type BrokerResponseItem struct {
	BrokerStats
	Broker *Broker      `json:"broker,omitempty"`
	Stock  *BrokerStock `json:"stock,omitempty"`
}

func (c *Client) GetBrokers(ctx context.Context, region Region, page, size int) (PaginatedResponse[*Broker], error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/api/v1/brokers", c.baseUrl), nil)
	if err != nil {
		return PaginatedResponse[*Broker]{}, err
	}

	q := req.URL.Query()
	q.Add("region", string(region))
	q.Add("page", strconv.Itoa(page))
	q.Add("size", strconv.Itoa(size))
	req.URL.RawQuery = q.Encode()

	resp, err := sendRequest[PaginatedResponse[*Broker]](ctx, c, req)
	if err != nil {
		return PaginatedResponse[*Broker]{}, err
	}

	return resp, nil
}

func (c *Client) GetMarketStocks(ctx context.Context, region Region, sortBy BrokerSort, sortDirection BrokerSortDirection, fromDate, toDate string, page, size int) (BrokerListResponse, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/api/v1/brokers/market/stock", c.baseUrl), nil)
	if err != nil {
		return BrokerListResponse{}, err
	}

	q := req.URL.Query()
	q.Add("region", string(region))
	q.Add("sortBy", string(sortBy))
	q.Add("sortDirection", string(sortDirection))
	q.Add("fromDate", fromDate)
	q.Add("toDate", toDate)
	q.Add("page", strconv.Itoa(page))
	q.Add("size", strconv.Itoa(size))
	req.URL.RawQuery = q.Encode()

	resp, err := sendRequest[BrokerListResponse](ctx, c, req)
	if err != nil {
		return BrokerListResponse{}, err
	}

	return resp, nil
}

func (c *Client) GetMarketBrokers(ctx context.Context, region Region, sortBy BrokerSort, sortDirection BrokerSortDirection, fromDate, toDate string, page, size int) (BrokerListResponse, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/api/v1/brokers/market", c.baseUrl), nil)
	if err != nil {
		return BrokerListResponse{}, err
	}

	q := req.URL.Query()
	q.Add("region", string(region))
	q.Add("sortBy", string(sortBy))
	q.Add("sortDirection", string(sortDirection))
	q.Add("fromDate", fromDate)
	q.Add("toDate", toDate)
	q.Add("page", strconv.Itoa(page))
	q.Add("size", strconv.Itoa(size))
	req.URL.RawQuery = q.Encode()

	resp, err := sendRequest[BrokerListResponse](ctx, c, req)
	if err != nil {
		return BrokerListResponse{}, err
	}

	return resp, nil
}

func (c *Client) GetBrokersByStock(ctx context.Context, symbol string, region Region, sortBy BrokerSort, sortDirection BrokerSortDirection, fromDate, toDate string, page, size int) (BrokerListResponse, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/api/v1/brokers/%s", c.baseUrl, symbol), nil)
	if err != nil {
		return BrokerListResponse{}, err
	}

	q := req.URL.Query()
	q.Add("symbol", symbol)
	q.Add("region", string(region))
	q.Add("sortBy", string(sortBy))
	q.Add("sortDirection", string(sortDirection))
	q.Add("fromDate", fromDate)
	q.Add("toDate", toDate)
	q.Add("page", strconv.Itoa(page))
	q.Add("size", strconv.Itoa(size))
	req.URL.RawQuery = q.Encode()

	resp, err := sendRequest[BrokerListResponse](ctx, c, req)
	if err != nil {
		return BrokerListResponse{}, err
	}

	return resp, nil
}

func (c *Client) GetStocksByBroker(ctx context.Context, symbol string, region Region, sortBy BrokerSort, sortDirection BrokerSortDirection, fromDate, toDate string, page, size int) (BrokerListResponse, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/api/v1/brokers/stock/%s", c.baseUrl, symbol), nil)
	if err != nil {
		return BrokerListResponse{}, err
	}

	q := req.URL.Query()
	q.Add("symbol", symbol)
	q.Add("region", string(region))
	q.Add("sortBy", string(sortBy))
	q.Add("sortDirection", string(sortDirection))
	q.Add("fromDate", fromDate)
	q.Add("toDate", toDate)
	q.Add("page", strconv.Itoa(page))
	q.Add("size", strconv.Itoa(size))
	req.URL.RawQuery = q.Encode()

	resp, err := sendRequest[BrokerListResponse](ctx, c, req)
	if err != nil {
		return BrokerListResponse{}, err
	}

	return resp, nil
}
