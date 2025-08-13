package laplace

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

type MarketState struct {
	ID            int       `json:"id"`
	MarketSymbol  *string   `json:"marketSymbol,omitempty"`
	State         string    `json:"state"`
	LastTimestamp time.Time `json:"lastTimestamp"`
	StockSymbol   *string   `json:"stockSymbol,omitempty"`
}

// GetStateOfAllMarkets returns the state of all markets for a given region.
func (c *Client) GetStateOfAllMarkets(ctx context.Context, region Region, page, size int) (PaginatedResponse[*MarketState], error) {
	endpoint := fmt.Sprintf("%s/api/v1/state/all", c.baseUrl)
	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return PaginatedResponse[*MarketState]{}, err
	}

	q := req.URL.Query()
	q.Add("region", string(region))
	q.Add("page", strconv.Itoa(page))
	q.Add("size", strconv.Itoa(size))
	req.URL.RawQuery = q.Encode()

	resp, err := sendRequest[PaginatedResponse[*MarketState]](ctx, c, req)
	if err != nil {
		return PaginatedResponse[*MarketState]{}, err
	}

	return resp, nil
}
