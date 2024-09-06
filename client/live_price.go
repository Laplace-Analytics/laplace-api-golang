package client

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

func getLivePrice[T any](c *Client, ctx context.Context, symbols []string, region Region) (data <-chan T, errors <-chan error, close func(), err error) {
	// Construct the URL
	streamID := uuid.New().String()
	url := fmt.Sprintf("%s/api/v1/stock/price/live?filter=%s&region=%s&stream=%s", c.baseUrl, strings.Join(symbols, ","), string(region), streamID)

	// Create a new request
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, nil, nil, err
	}
	return sendSSERequest[T](ctx, c, req)
}

type BISTStockLiveData struct {
	Symbol             string  `json:"s"`
	DailyPercentChange float64 `json:"ch"`
	ClosePrice         float64 `json:"p"`
}

func (c *Client) GetLivePriceForBIST(ctx context.Context, symbols []string, region Region) (data <-chan BISTStockLiveData, errors <-chan error, close func(), err error) {
	return getLivePrice[BISTStockLiveData](c, ctx, symbols, region)
}

type USStockLiveData struct {
	Symbol   string  `json:"s"`
	BidPrice float64 `json:"bp"`
	AskPrice float64 `json:"ap"`
}

func (c *Client) GetLivePriceForUS(ctx context.Context, symbols []string, region Region) (data <-chan USStockLiveData, errors <-chan error, close func(), err error) {
	return getLivePrice[USStockLiveData](c, ctx, symbols, region)
}
