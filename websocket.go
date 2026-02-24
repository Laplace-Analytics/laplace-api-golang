package laplace

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type MessageCode string

const (
	MessageCodeNewUser            MessageCode = "new_user"
	MessageCodeHasNoAccessToLevel MessageCode = "no_access_to_level"
)

type FeedType string

const (
	FeedTypeLivePriceTR       FeedType = "live_price_tr"
	FeedTypeDelayedPriceTR    FeedType = "delayed_price_tr"
	FeedTypeLivePriceUS       FeedType = "live_price_us"
	FeedTypeDelayedPriceUS    FeedType = "delayed_price_us"
	FeedTypeDepthTR           FeedType = "depth_tr"
	FeedTypeStateUS           FeedType = "state_us"
	FeedTypeLiveAskBidPriceTR FeedType = "live_ask_bid_price_tr"
	FeedTypeCustom            FeedType = "custom"
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

// GetWebSocketUrl generates a WebSocket URL for accessing real-time market data feeds including live prices and depth data.
func (c *Client) GetWebSocketUrl(ctx context.Context, externalUserId string, feeds []FeedType) (string, error) {
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

	resp, err := sendRequest[WebSocketUrlResponse](ctx, c, req)
	if err != nil {
		return "", err
	}

	return resp.URL, nil
}

