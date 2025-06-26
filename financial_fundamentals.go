package laplace

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type StockDividend struct {
	Date           time.Time `json:"date"`
	NetAmount      float64   `json:"netAmount"`
	NetRatio       float64   `json:"netRatio"`
	GrossAmount    float64   `json:"grossAmount"`
	GrossRatio     float64   `json:"grossRatio"`
	PriceThen      float64   `json:"priceThen"`
	StoppageRatio  float64   `json:"stoppageRatio"`
	StoppageAmount float64   `json:"stoppageAmount"`
}

type StockStats struct {
	PreviousClose    float64 `json:"previousClose,omitempty"`
	YtdReturn        float64 `json:"ytdReturn,omitempty"`
	YearlyReturn     float64 `json:"yearlyReturn,omitempty"`
	MarketCap        float64 `json:"marketCap,omitempty"`
	PeRatio          float64 `json:"peRatio,omitempty"`
	PbRatio          float64 `json:"pbRatio,omitempty"`
	YearLow          float64 `json:"yearLow,omitempty"`
	YearHigh         float64 `json:"yearHigh,omitempty"`
	ThreeYearReturn  float64 `json:"3YearReturn,omitempty"`
	FiveYearReturn   float64 `json:"5YearReturn,omitempty"`
	ThreeMonthReturn float64 `json:"3MonthReturn,omitempty"`
	MonthlyReturn    float64 `json:"monthlyReturn,omitempty"`
	WeeklyReturn     float64 `json:"weeklyReturn,omitempty"`
	Symbol           string  `json:"symbol"`
	LatestPrice      float64 `json:"latestPrice,omitempty"`
	DailyChange      float64 `json:"dailyChange,omitempty"`
	DayHigh          float64 `json:"dayHigh,omitempty"`
	DayLow           float64 `json:"dayLow,omitempty"`
	LowerPriceLimit  Price   `json:"lowerPriceLimit,omitempty"`
	UpperPriceLimit  Price   `json:"upperPriceLimit,omitempty"`
	DayOpen          float64 `json:"dayOpen,omitempty"`
	Eps              float64 `json:"eps,omitempty"`
}

type Price float64

func (p Price) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("%.2f", float64(p))), nil
}

type TopMover struct {
	Symbol     string  `json:"symbol"`
	AssetClass string  `json:"assetClass,omitempty"`
	AssetType  string  `json:"assetType,omitempty"`
	Change     float64 `json:"change"`
}

type TopMoversDirection string

const (
	TopMoversDirectionGainers TopMoversDirection = "gainers"
	TopMoversDirectionLosers  TopMoversDirection = "losers"
)

func (c *Client) GetStockDividends(ctx context.Context, symbol string, region Region) ([]StockDividend, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/api/v2/stock/dividends", c.baseUrl), nil)
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

func (c *Client) GetStockStats(ctx context.Context, symbols []string, region Region) ([]StockStats, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/api/v2/stock/stats", c.baseUrl), nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Add("symbols", strings.Join(symbols, ","))
	q.Add("region", string(region))
	req.URL.RawQuery = q.Encode()

	resp, err := sendRequest[[]StockStats](ctx, c, req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *Client) GetTopMovers(ctx context.Context, direction TopMoversDirection, assetClass *AssetClass, assetType *AssetType, page int, pageSize int, region Region) ([]TopMover, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/api/v2/stock/top-movers", c.baseUrl), nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Add("direction", string(direction))
	q.Add("pageSize", strconv.Itoa(pageSize))
	q.Add("page", strconv.Itoa(page))
	q.Add("region", string(region))
	
	if assetClass != nil {
		q.Add("assetClass", string(*assetClass))
	}
	if assetType != nil {
		q.Add("assetType", string(*assetType))
	}
	
	req.URL.RawQuery = q.Encode()

	resp, err := sendRequest[[]TopMover](ctx, c, req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
