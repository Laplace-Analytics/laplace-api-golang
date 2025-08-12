package laplace

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

type EarningsTranscriptWithSummary struct {
	Symbol     string    `json:"symbol"`
	Year       int       `json:"year"`
	Quarter    int       `json:"quarter"`
	Date       time.Time `json:"date"`
	Content    string    `json:"content"`
	Summary    string    `json:"summary,omitempty"`
	HasSummary bool      `json:"has_summary"`
}

func (c *Client) GetEarningsTranscriptWithSummary(ctx context.Context, symbol string, year, quarter int) (*EarningsTranscriptWithSummary, error) {
	endpoint := fmt.Sprintf("%s/api/v1/earnings/transcript", c.baseUrl)
	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Add("symbol", symbol)
	q.Add("year", strconv.Itoa(year))
	q.Add("quarter", strconv.Itoa(quarter))
	req.URL.RawQuery = q.Encode()

	res, err := sendRequest[*EarningsTranscriptWithSummary](ctx, c, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}
