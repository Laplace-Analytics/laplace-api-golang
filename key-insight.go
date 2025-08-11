package laplace

import (
	"context"
	"fmt"
	"net/http"
)

type KeyInsight struct {
	Symbol  string `json:"symbol"`
	Insight string `json:"insight"`
}

// GetKeyInsights fetches key insights and analysis for a specific stock symbol.
func (c *Client) GetKeyInsights(ctx context.Context, symbol string, region Region) (*KeyInsight, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/api/v1/key-insight", c.baseUrl), nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Add("symbol", symbol)
	q.Add("region", string(region))
	req.URL.RawQuery = q.Encode()

	resp, err := sendRequest[KeyInsight](ctx, c, req)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}
