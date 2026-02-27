package laplace

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
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

// RevokeWebSocketConnection revokes an active WebSocket connection by its ID (the UUID segment from the WebSocket URL).
func (c *Client) RevokeWebSocketConnection(ctx context.Context, id string) error {
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/api/v1/ws/user/revoke/%s", c.baseUrl, id), nil)
	if err != nil {
		return err
	}

	_, err = sendRequest[any](ctx, c, req)
	return err
}

type WebSocketMonthlyUsageData struct {
	ExternalUserID      string    `json:"externalUserID"`
	FirstConnectionTime time.Time `json:"firstConnectionTime"`
	UniqueDeviceCount   int64     `json:"uniqueDeviceCount"`
}

// GetWebsocketUsageForMonth retrieves WebSocket usage statistics for a specific month.
func (c *Client) GetWebsocketUsageForMonth(ctx context.Context, month int, year int, feedType FeedType) ([]WebSocketMonthlyUsageData, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/api/v1/ws/report", c.baseUrl), nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Add("month", strconv.Itoa(month))
	q.Add("year", strconv.Itoa(year))
	q.Add("feedType", string(feedType))
	req.URL.RawQuery = q.Encode()

	resp, err := sendRequest[[]WebSocketMonthlyUsageData](ctx, c, req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

type SendWebsocketEventRequest struct {
	ExternalUserID string          `json:"externalUserID,omitempty"`
	Event          json.RawMessage `json:"event"`
	Transient      *bool           `json:"transient,omitempty"`
	BroadCastToAll bool            `json:"broadCastToAll"`
}

// SendWebsocketEvent sends a custom event through the WebSocket connection.
func (c *Client) SendWebsocketEvent(ctx context.Context, params SendWebsocketEventRequest) error {
	jsonData, err := json.Marshal(params)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/api/v1/ws/event", c.baseUrl), bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	_, err = sendRequest[any](ctx, c, req)
	return err
}

