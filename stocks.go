package laplace

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AssetType string

const (
	AssetTypeStock       AssetType = "stock"
	AssetTypeForex       AssetType = "forex"
	AssetTypeIndex       AssetType = "index"
	AssetTypeEtf         AssetType = "etf"
	AssetTypeCommodity   AssetType = "commodity"
	AssetTypeStockRights AssetType = "stock_rights"
	AssetTypeFund        AssetType = "fund"
	AssetTypeAll         AssetType = "all"
)

type AssetClass string

const (
	AssetClassEquity AssetClass = "equity"
	AssetClassCrypto AssetClass = "crypto"
	AssetClassADR    AssetClass = "adr"
	AssetClassETN    AssetClass = "etn"
	AssetClassAll    AssetClass = "all"
)

type Stock struct {
	ID          primitive.ObjectID `json:"id"`
	AssetType   AssetType          `json:"assetType"`
	Name        string             `json:"name"`
	Symbol      string             `json:"symbol"`
	SectorId    string             `json:"sectorId"`
	IndustryId  string             `json:"industryId"`
	UpdatedDate time.Time          `json:"updatedDate"`
	DailyChange float64            `json:"dailyChange,omitempty"`
	Active      bool               `json:"active"`
}

type LocaleString map[Locale]string

type StockDetail struct {
	ID                        primitive.ObjectID `json:"id"`
	AssetType                 AssetType          `json:"assetType"`
	AssetClass                AssetClass         `json:"assetClass"`
	Name                      string             `json:"name"`
	Symbol                    string             `json:"symbol"`
	Description               string             `json:"description"`
	LocalizedDescription      LocaleString       `json:"localized_description"`
	ShortDescription          string             `json:"shortDescription"`
	LocalizedShortDescription LocaleString       `json:"localizedShortDescription"`
	Region                    string             `json:"region"`
	SectorId                  string             `json:"sectorId"`
	IndustryId                string             `json:"industryId"`
	UpdatedDate               time.Time          `json:"updatedDate"`
	Active                    bool               `json:"active"`
	Markets                   []Market           `json:"markets,omitempty"`
}

type Market string

const (
	MarketYildiz      Market = "YILDIZ"
	MarketAna         Market = "ANA"
	MarketAlt         Market = "ALT"
	MarketYakinIzleme Market = "YAKIN_IZLEME"
	MarketPOIP        Market = "POIP"
	MarketFon         Market = "FON"
	MarketGirisim     Market = "GIRISIM"
	MarketEmtia       Market = "EMTIA"
)

type HistoricalPricePeriod string

const (
	HistoricalPricePeriodOneDay     HistoricalPricePeriod = "1D"
	HistoricalPricePeriodOneWeek    HistoricalPricePeriod = "1W"
	HistoricalPricePeriodOneMonth   HistoricalPricePeriod = "1M"
	HistoricalPricePeriodThreeMonth HistoricalPricePeriod = "3M"
	HistoricalPricePeriodOneYear    HistoricalPricePeriod = "1Y"
	HistoricalPricePeriodTwoYear    HistoricalPricePeriod = "2Y"
	HistoricalPricePeriodThreeYear  HistoricalPricePeriod = "3Y"
	HistoricalPricePeriodSixMonth   HistoricalPricePeriod = "6M"
	HistoricalPricePeriodFiveYear   HistoricalPricePeriod = "5Y"
	HistoricalPricePeriodAll        HistoricalPricePeriod = "All"
)

type HistoricalPriceInterval string

const (
	HistoricalPriceIntervalOneMinute     HistoricalPriceInterval = "1m"
	HistoricalPriceIntervalThreeMinute   HistoricalPriceInterval = "3m"
	HistoricalPriceIntervalFiveMinute    HistoricalPriceInterval = "5m"
	HistoricalPriceIntervalFifteenMinute HistoricalPriceInterval = "15m"
	HistoricalPriceIntervalThirtyMinute  HistoricalPriceInterval = "30m"
	HistoricalPriceIntervalOneHour       HistoricalPriceInterval = "1h"
	HistoricalPriceIntervalTwoHour       HistoricalPriceInterval = "2h"
	HistoricalPriceIntervalOneDay        HistoricalPriceInterval = "1d"
	HistoricalPriceIntervalFiveDay       HistoricalPriceInterval = "5d"
	HistoricalPriceIntervalSevenDay      HistoricalPriceInterval = "7d"
	HistoricalPriceIntervalThirtyDay     HistoricalPriceInterval = "30d"
)

type HistoricalPriceDate struct {
	Year   int
	Month  int
	Day    int
	Hour   int
	Minute int
	Second int
}

type StockPriceGraph struct {
	Symbol     string           `json:"symbol"`
	OneDay     []PriceDataPoint `json:"1D"`
	OneWeek    []PriceDataPoint `json:"1W"`
	OneMonth   []PriceDataPoint `json:"1M"`
	ThreeMonth []PriceDataPoint `json:"3M"`
	OneYear    []PriceDataPoint `json:"1Y"`
	TwoYear    []PriceDataPoint `json:"2Y"`
	ThreeYear  []PriceDataPoint `json:"3Y"`
	FiveYear   []PriceDataPoint `json:"5Y"`
}

type PriceDataPoint struct {
	Date            int64   `json:"d"`
	Open            float64 `json:"o"`
	UnadjustedOpen  float64 `json:"uo,omitempty"`
	High            float64 `json:"h"`
	UnadjustedHigh  float64 `json:"uh,omitempty"`
	Low             float64 `json:"l"`
	UnadjustedLow   float64 `json:"ul,omitempty"`
	Close           float64 `json:"c"`
	UnadjustedClose float64 `json:"uc,omitempty"`
	Volume          float64 `json:"v,omitempty"`
	UnadjustedVol   float64 `json:"uv,omitempty"`
}

type StockRestriction struct {
	ID          int        `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Symbol      string     `json:"symbol,omitempty"`
	StartDate   *time.Time `json:"startDate"`
	EndDate     *time.Time `json:"endDate"`
	Market      string     `json:"market,omitempty"`
}

type TickRule struct {
	BasePrice       float64        `json:"basePrice"`
	AdditionalPrice int            `json:"additionalPrice"`
	LowerPriceLimit float64        `json:"lowerPriceLimit"`
	UpperPriceLimit float64        `json:"upperPriceLimit"`
	Rules           []TickSizeRule `json:"rules"`
}

type TickSizeRule struct {
	PriceFrom float64 `json:"priceFrom"`
	PriceTo   float64 `json:"priceTo"`
	TickSize  float64 `json:"tickSize"`
}

// GetAllStocks retrieves a paginated list of all stocks for the specified region.
func (c *Client) GetAllStocks(ctx context.Context, region Region, page int, pageSize int) ([]Stock, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/api/v2/stock/all", c.baseUrl), nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Add("region", string(region))
	if pageSize > 0 {
		q.Add("page", strconv.Itoa(page))
		q.Add("pageSize", strconv.Itoa(pageSize))
	}
	req.URL.RawQuery = q.Encode()

	resp, err := sendRequest[[]Stock](ctx, c, req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// GetStockDetailByID fetches detailed information about a stock using its unique ID.
func (c *Client) GetStockDetailByID(ctx context.Context, id string, locale Locale) (StockDetail, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/api/v1/stock/%s", c.baseUrl, id), nil)
	if err != nil {
		return StockDetail{}, err
	}

	q := req.URL.Query()
	q.Add("locale", string(locale))
	req.URL.RawQuery = q.Encode()

	resp, err := sendRequest[StockDetail](ctx, c, req)
	if err != nil {
		return StockDetail{}, err
	}

	return resp, nil
}

// GetStockDetailBySymbol fetches detailed information about a stock using its symbol, asset class, region, and locale.
func (c *Client) GetStockDetailBySymbol(ctx context.Context, symbol string, assetClass AssetClass, region Region, locale Locale) (StockDetail, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/api/v1/stock/detail", c.baseUrl), nil)
	if err != nil {
		return StockDetail{}, err
	}

	q := req.URL.Query()
	q.Add("symbol", symbol)
	q.Add("asset_class", string(assetClass))
	q.Add("region", string(region))
	q.Add("locale", string(locale))
	req.URL.RawQuery = q.Encode()

	resp, err := sendRequest[StockDetail](ctx, c, req)
	if err != nil {
		return StockDetail{}, err
	}

	return resp, nil
}

// GetHistoricalPrices retrieves historical price data for multiple stocks over specified time periods.
func (c *Client) GetHistoricalPrices(ctx context.Context, symbols []string, region Region, keys []HistoricalPricePeriod) ([]StockPriceGraph, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/api/v1/stock/price", c.baseUrl), nil)
	if err != nil {
		return []StockPriceGraph{}, err
	}

	q := req.URL.Query()
	q.Add("symbols", strings.Join(symbols, ","))
	q.Add("region", string(region))
	q.Add("keys", strings.Join(lo.Map(keys, func(key HistoricalPricePeriod, _ int) string {
		return string(key)
	}), ","))
	req.URL.RawQuery = q.Encode()

	resp, err := sendRequest[[]StockPriceGraph](ctx, c, req)
	if err != nil {
		return []StockPriceGraph{}, err
	}

	return resp, nil
}

// GetCustomHistoricalPrices retrieves custom historical price data for a stock within a specific date range and interval.
func (c *Client) GetCustomHistoricalPrices(ctx context.Context, symbol string, region Region, fromDate string, toDate string, interval HistoricalPriceInterval, detail bool, numIntervals ...int) ([]PriceDataPoint, error) {
	if err := validateCustomHistoricalPriceDate(fromDate); err != nil {
		return []PriceDataPoint{}, err
	}

	if err := validateCustomHistoricalPriceDate(toDate); err != nil {
		return []PriceDataPoint{}, err
	}

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/api/v1/stock/price/interval", c.baseUrl), nil)
	if err != nil {
		return []PriceDataPoint{}, err
	}

	q := req.URL.Query()
	q.Add("stock", symbol)
	q.Add("region", string(region))
	q.Add("fromDate", fromDate)
	q.Add("toDate", toDate)
	q.Add("interval", string(interval))
	q.Add("detail", strconv.FormatBool(detail))
	if len(numIntervals) > 0 {
		q.Add("numIntervals", strconv.Itoa(numIntervals[0]))
	}
	req.URL.RawQuery = q.Encode()

	resp, err := sendRequest[[]PriceDataPoint](ctx, c, req)
	if err != nil {
		return []PriceDataPoint{}, err
	}

	return resp, nil
}

func validateCustomHistoricalPriceDate(date string) error {
	pattern := `^\d{4}-\d{2}-\d{2}( \d{2}:\d{2}:\d{2})?$`
	matched, err := regexp.MatchString(pattern, date)
	if err != nil || !matched {
		return fmt.Errorf("invalid date format, allowed formats: YYYY-MM-DD, YYYY-MM-DD HH:MM:SS")
	}

	return nil
}

// GetStockRestrictions fetches trading restrictions and limitations for a specific stock in the given region.
func (c *Client) GetStockRestrictions(ctx context.Context, symbol string, region Region) ([]StockRestriction, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/api/v1/stock/restrictions", c.baseUrl), nil)
	if err != nil {
		return []StockRestriction{}, err
	}

	q := req.URL.Query()
	q.Add("symbol", symbol)
	q.Add("region", string(region))
	req.URL.RawQuery = q.Encode()

	resp, err := sendRequest[[]StockRestriction](ctx, c, req)
	if err != nil {
		return []StockRestriction{}, err
	}

	return resp, nil
}

// GetAllRestrictions retrieves all trading restrictions and limitations.
func (c *Client) GetAllRestrictions(ctx context.Context) ([]StockRestriction, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/api/v1/stock/restrictions/all", c.baseUrl), nil)
	if err != nil {
		return []StockRestriction{}, err
	}

	resp, err := sendRequest[[]StockRestriction](ctx, c, req)
	if err != nil {
		return []StockRestriction{}, err
	}

	return resp, nil
}

// GetTickRules retrieves tick size rules and price limits for a stock in the specified region.
func (c *Client) GetTickRules(ctx context.Context, symbol string, region Region) (TickRule, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/api/v1/stock/rules", c.baseUrl), nil)
	if err != nil {
		return TickRule{}, err
	}

	q := req.URL.Query()
	q.Add("symbol", symbol)
	q.Add("region", string(region))
	req.URL.RawQuery = q.Encode()

	resp, err := sendRequest[TickRule](ctx, c, req)
	if err != nil {
		return TickRule{}, err
	}

	return resp, nil
}

type GenerateChartImageRequest struct {
	Symbol     string
	Period     HistoricalPricePeriod
	Region     Region
	Resolution HistoricalPriceInterval
	Indicators []string
	ChartType  *int
}

// GetStockChartImage generates a chart image for a stock and returns the raw image bytes.
func (c *Client) GetStockChartImage(ctx context.Context, params GenerateChartImageRequest) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/api/v1/stock/chart", c.baseUrl), nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Add("symbol", params.Symbol)
	q.Add("region", string(params.Region))
	if params.Period != "" {
		q.Add("period", string(params.Period))
	}
	if params.Resolution != "" {
		q.Add("resolution", string(params.Resolution))
	}
	if len(params.Indicators) > 0 {
		q.Add("indicators", strings.Join(params.Indicators, ","))
	}
	if params.ChartType != nil {
		q.Add("chartType", strconv.Itoa(*params.ChartType))
	}
	req.URL.RawQuery = q.Encode()

	resp, err := sendRawRequest(ctx, c, req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
