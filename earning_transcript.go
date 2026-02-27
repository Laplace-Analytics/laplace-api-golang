package laplace

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// FlexibleTime handles date fields that may arrive in different formats from the API
// (e.g. "2024-01-15", "2024-01-15T10:30:00Z", unix timestamp as number).
type FlexibleTime struct {
	time.Time
}

func (ft *FlexibleTime) UnmarshalJSON(data []byte) error {
	s := strings.Trim(string(data), "\"")
	if s == "null" || s == "" {
		return nil
	}

	// Try unix timestamp (number)
	if ts, err := strconv.ParseInt(s, 10, 64); err == nil {
		ft.Time = time.Unix(ts, 0)
		return nil
	}

	// Try common date/time formats
	formats := []string{
		time.RFC3339,
		"2006-01-02T15:04:05",
		"2006-01-02",
	}
	for _, f := range formats {
		if t, err := time.Parse(f, s); err == nil {
			ft.Time = t
			return nil
		}
	}

	return fmt.Errorf("FlexibleTime: unable to parse %q", s)
}

func (ft FlexibleTime) MarshalJSON() ([]byte, error) {
	if ft.Time.IsZero() {
		return []byte("null"), nil
	}
	return json.Marshal(ft.Time)
}

type EarningsTranscriptWithSummary struct {
	Symbol     string       `json:"symbol"`
	Year       int          `json:"year"`
	Quarter    int          `json:"quarter"`
	Date       FlexibleTime `json:"date"`
	Content    string       `json:"content"`
	Summary    string       `json:"summary,omitempty"`
	HasSummary bool         `json:"has_summary"`
}

type EarningsTranscriptListItem struct {
	Symbol     string       `json:"symbol"`
	Year       int          `json:"year"`
	Quarter    int          `json:"quarter"`
	Date       FlexibleTime `json:"date"`
	FiscalYear int          `json:"fiscal_year"`
}

// GetEarningsTranscriptWithSummary retrieves the earnings transcript with an AI-generated summary for a specific quarter.
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

// GetEarningsTranscriptList retrieves the list of available earnings transcripts for a stock.
func (c *Client) GetEarningsTranscriptList(ctx context.Context, region Region, symbol string) ([]EarningsTranscriptListItem, error) {
	endpoint := fmt.Sprintf("%s/api/v1/earnings/transcripts", c.baseUrl)
	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Add("symbol", symbol)
	q.Add("region", string(region))
	req.URL.RawQuery = q.Encode()

	res, err := sendRequest[[]EarningsTranscriptListItem](ctx, c, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}
