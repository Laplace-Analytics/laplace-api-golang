package laplace

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type ScreenerSortBy string

const (
	ScreenerSortBySymbol           ScreenerSortBy = "symbol"
	ScreenerSortByPrice            ScreenerSortBy = "price"
	ScreenerSortByDailyChange      ScreenerSortBy = "dailyChange"
	ScreenerSortByMarketCap        ScreenerSortBy = "marketCap"
	ScreenerSortByPERatio          ScreenerSortBy = "peRatio"
	ScreenerSortByPBRatio          ScreenerSortBy = "pbRatio"
	ScreenerSortByWeeklyReturn     ScreenerSortBy = "weeklyReturn"
	ScreenerSortByMonthlyReturn    ScreenerSortBy = "monthlyReturn"
	ScreenerSortByThreeMonthReturn ScreenerSortBy = "threeMonthReturn"
	ScreenerSortByYearlyReturn     ScreenerSortBy = "yearlyReturn"
	ScreenerSortByThreeYearReturn  ScreenerSortBy = "threeYearReturn"
	ScreenerSortByFiveYearReturn   ScreenerSortBy = "fiveYearReturn"
	ScreenerSortByYTDReturn        ScreenerSortBy = "ytdReturn"
)

type ScreenerRange struct {
	Min *float64 `json:"min,omitempty"`
	Max *float64 `json:"max,omitempty"`
}

type ScreenerFilters struct {
	Price            *ScreenerRange `json:"price,omitempty"`
	DailyChange      *ScreenerRange `json:"dailyChange,omitempty"`
	PERatio          *ScreenerRange `json:"peRatio,omitempty"`
	PBRatio          *ScreenerRange `json:"pbRatio,omitempty"`
	MarketCap        *ScreenerRange `json:"marketCap,omitempty"`
	WeeklyReturn     *ScreenerRange `json:"weeklyReturn,omitempty"`
	MonthlyReturn    *ScreenerRange `json:"monthlyReturn,omitempty"`
	ThreeMonthReturn *ScreenerRange `json:"threeMonthReturn,omitempty"`
	YearlyReturn     *ScreenerRange `json:"yearlyReturn,omitempty"`
	ThreeYearReturn  *ScreenerRange `json:"threeYearReturn,omitempty"`
	FiveYearReturn   *ScreenerRange `json:"fiveYearReturn,omitempty"`
	YTDReturn        *ScreenerRange `json:"ytdReturn,omitempty"`
}

type ScreenerRequest struct {
	Filters   *ScreenerFilters `json:"filters,omitempty"`
	SortBy    ScreenerSortBy   `json:"sortBy,omitempty"`
	SortOrder SortDirection    `json:"sortOrder,omitempty"`
	Page      int              `json:"page,omitempty"`
	PageSize  int              `json:"pageSize,omitempty"`
}

type ScreenerItem struct {
	Symbol           string   `json:"symbol"`
	Price            *float64 `json:"price"`
	DailyChange      *float64 `json:"dailyChange"`
	MarketCap        *float64 `json:"marketCap"`
	PERatio          *float64 `json:"peRatio"`
	PBRatio          *float64 `json:"pbRatio"`
	WeeklyReturn     *float64 `json:"weeklyReturn"`
	MonthlyReturn    *float64 `json:"monthlyReturn"`
	ThreeMonthReturn *float64 `json:"threeMonthReturn"`
	YearlyReturn     *float64 `json:"yearlyReturn"`
	ThreeYearReturn  *float64 `json:"threeYearReturn"`
	FiveYearReturn   *float64 `json:"fiveYearReturn"`
	YTDReturn        *float64 `json:"ytdReturn"`
}

type ScreenerResponse struct {
	Items       []ScreenerItem `json:"items"`
	RecordCount int            `json:"recordCount"`
}

// Screener returns a filtered and sorted list of stocks for the given region.
// Region defaults to "tr" if empty. US is not currently supported.
func (c *Client) Screener(ctx context.Context, region Region, params ScreenerRequest) (ScreenerResponse, error) {
	bodyJSON, err := json.Marshal(params)
	if err != nil {
		return ScreenerResponse{}, err
	}

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/api/v1/screener", c.baseUrl), bytes.NewBuffer(bodyJSON))
	if err != nil {
		return ScreenerResponse{}, err
	}

	if region != "" {
		q := req.URL.Query()
		q.Add("region", string(region))
		req.URL.RawQuery = q.Encode()
	}

	req.Header.Set("Content-Type", "application/json")

	return sendRequest[ScreenerResponse](ctx, c, req)
}
