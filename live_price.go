package laplace

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

func getLivePrice[T StockLiveData](c *Client, ctx context.Context, symbols []string, region Region) (data <-chan T, errors <-chan error, close func(), err error) {
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

type StockLiveData interface{}

type BISTStockLiveData struct {
	Symbol             string  `json:"s"`
	DailyPercentChange float64 `json:"ch"`
	ClosePrice         float64 `json:"p"`
	StockLiveData
}

func (c *Client) GetLivePriceForBIST(ctx context.Context, symbols []string) (data <-chan BISTStockLiveData, errors <-chan error, close func(), err error) {
	return getLivePrice[BISTStockLiveData](c, ctx, symbols, RegionTr)
}

type USStockLiveData struct {
	Symbol   string  `json:"s"`
	BidPrice float64 `json:"bp"`
	AskPrice float64 `json:"ap"`
	StockLiveData
}

func (c *Client) GetLivePriceForUS(ctx context.Context, symbols []string) (data <-chan USStockLiveData, errors <-chan error, close func(), err error) {
	return getLivePrice[USStockLiveData](c, ctx, symbols, RegionUs)
}

type MessageCode string

const (
	MessageCodeNewUser            MessageCode = "new_user"
	MessageCodeHasNoAccessToLevel MessageCode = "no_access_to_level"
)

type FeedType string

const (
	FeedTypeLivePriceTR    FeedType = "live_price_tr"
	FeedTypeDelayedPriceTR FeedType = "delayed_price_tr"
	FeedTypeLivePriceUS    FeedType = "live_price_us"
	FeedTypeDelayedPriceUS FeedType = "delayed_price_us"
	FeedTypeDepthTR        FeedType = "depth_tr"
)

type WebSocketUrlParams struct {
	ExternalUserId string     `json:"externalUserId"`
	Feeds          []FeedType `json:"feeds"`
}

type WebSocketUrlResponse struct {
	URL         string      `json:"url,omitempty"`
	Message     string      `json:"message,omitempty"`
	Code        MessageCode `json:"code,omitempty"`
	ExampleBody any         `json:"exampleBody,omitempty"`
}

type AccessorType string

const (
	AccessorTypeUser AccessorType = "user"
)

type UpdateUserDetailsParams struct {
	ExternalUserID string       `json:"externalUserID"`
	FirstName      string       `json:"firstName"`
	LastName       string       `json:"lastName"`
	Address        string       `json:"address"`
	City           string       `json:"city"`
	CountryCode    string       `json:"countryCode"`
	AccessorType   AccessorType `json:"accessorType"`
	Active         bool         `json:"active"`
}

func (c *Client) GetWebSocketUrl(ctx context.Context, externalUserId string, feeds []FeedType, region Region) (string, error) {
	params := WebSocketUrlParams{
		ExternalUserId: externalUserId,
		Feeds:          feeds,
	}

	jsonData, err := json.Marshal(params)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/api/v2/ws/url", c.baseUrl), bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")

	q := req.URL.Query()
	q.Add("region", string(region))
	req.URL.RawQuery = q.Encode()

	resp, err := sendRequest[WebSocketUrlResponse](ctx, c, req)
	if err != nil {
		return "", err
	}

	return resp.URL, nil
}

func (c *Client) UpdateUserDetails(ctx context.Context, params UpdateUserDetailsParams) error {
	jsonData, err := json.Marshal(params)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("%s/api/v1/ws/user", c.baseUrl), bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	_, err = sendRequest[any](ctx, c, req)
	if err != nil {
		return err
	}

	return nil
}
