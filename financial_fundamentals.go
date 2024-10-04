package laplace

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/samber/lo"
)

type StockDividend struct {
	Date              time.Time `json:"date"`
	DividendAmount    float32   `json:"dividendAmount"`
	DividendRatio     float32   `json:"dividendRatio"`
	NetDividendAmount float32   `json:"netDividendAmount"`
	NetDividendRatio  float32   `json:"netDividendRatio"`
	PriceThen         float32   `json:"priceThen"`
}

type StockStats struct {
	LatestPrice      float64 `json:"latestPrice"`
	PreviousClose    float64 `json:"previousClose"`
	MarketCap        float64 `json:"marketCap"`
	PeRatio          float64 `json:"peRatio"`
	PbRatio          float64 `json:"pbRatio"`
	DayHigh          float64 `json:"dayHigh"`
	DayLow           float64 `json:"dayLow"`
	YearLow          float64 `json:"yearLow"`
	YearHigh         float64 `json:"yearHigh"`
	DailyChange      float64 `json:"dailyChange"`
	WeeklyReturn     float64 `json:"weeklyReturn"`
	MonthlyReturn    float64 `json:"monthlyReturn"`
	ThreeMonthReturn float64 `json:"3MonthReturn"`
	YtdReturn        float64 `json:"ytdReturn"`
	YearlyReturn     float64 `json:"yearlyReturn"`
	ThreeYear        float64 `json:"3YearReturn"`
	FiveYear         float64 `json:"5YearReturn"`
	Symbol           string  `json:"symbol"`
}

type StockStatsKey string

const (
	StockStatsLatestPrice   StockStatsKey = "latest_price"
	StockStatsPreviousClose StockStatsKey = "previous_close"
	StockStatsMarketCap     StockStatsKey = "market_cap"
	StockStatsFK            StockStatsKey = "fk"
	StockStatsPDDD          StockStatsKey = "pddd"
	StockStatsDayLow        StockStatsKey = "day_low"
	StockStatsDayHigh       StockStatsKey = "day_high"
	StockStatsYearLow       StockStatsKey = "year_low"
	StockStatsYearHigh      StockStatsKey = "year_high"
	StockStatsDailyChange   StockStatsKey = "daily_change"
	StockStatsWeeklyReturn  StockStatsKey = "weekly_return"
	StockStatsMonthlyReturn StockStatsKey = "monthly_return"
	StockStats3MonthReturn  StockStatsKey = "3_month_return"
	StockStatsYtdReturn     StockStatsKey = "ytd_return"
	StockStatsYearlyReturn  StockStatsKey = "yearly_return"
	StockStats3YearReturn   StockStatsKey = "3_year_return"
	StockStats5YearReturn   StockStatsKey = "5_year_return"
)

type TopMover struct {
	Symbol        string  `json:"symbol"`
	PercentChange float64 `json:"percent_change"`
}

func (c *Client) GetStockDividends(ctx context.Context, symbol string, region Region) ([]StockDividend, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/api/v1/stock/dividends", c.baseUrl), nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Add("symbol", symbol)
	q.Add("region", string(region))
	req.URL.RawQuery = q.Encode()

	resp, err := sendRequest[[]StockDividend](ctx, c, req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *Client) GetStockStats(ctx context.Context, symbols []string, keys []StockStatsKey, region Region) ([]StockStats, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/api/v1/stock/stats", c.baseUrl), nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Add("symbols", strings.Join(symbols, ","))
	q.Add("region", string(region))
	q.Add("keys", strings.Join(lo.Map(keys, func(key StockStatsKey, _ int) string { return string(key) }), ","))
	req.URL.RawQuery = q.Encode()

	resp, err := sendRequest[[]StockStats](ctx, c, req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *Client) GetTopMovers(ctx context.Context, region Region) ([]TopMover, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/api/v1/stock/top-movers", c.baseUrl), nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Add("region", string(region))
	req.URL.RawQuery = q.Encode()

	resp, err := sendRequest[[]TopMover](ctx, c, req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
