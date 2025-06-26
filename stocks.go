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
)

type AssetClass string

const (
	AssetClassEquity AssetClass = "equity"
	AssetClassCrypto AssetClass = "crypto"
)

type Stock struct {
	ID          primitive.ObjectID `json:"id"`
	AssetType   AssetType          `json:"assetType"`
	Name        string             `json:"name"`
	Symbol      string             `json:"symbol"`
	SectorId    primitive.ObjectID `json:"sectorId"`
	IndustryId  primitive.ObjectID `json:"industryId"`
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
	SectorId                  primitive.ObjectID `json:"sectorId"`
	IndustryId                primitive.ObjectID `json:"industryId"`
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
	HistoricalPricePeriodFiveYear   HistoricalPricePeriod = "5Y"
)

type HistoricalPriceInterval string

const (
	HistoricalPriceIntervalOneMinute    HistoricalPriceInterval = "1m"
	HistoricalPriceIntervalFiveMinute   HistoricalPriceInterval = "5m"
	HistoricalPriceIntervalThirtyMinute HistoricalPriceInterval = "30m"
	HistoricalPriceIntervalOneHour      HistoricalPriceInterval = "1h"
	HistoricalPriceIntervalOneDay       HistoricalPriceInterval = "1d"
	HistoricalPriceIntervalFiveDay      HistoricalPriceInterval = "5d"
	HistoricalPriceIntervalSevenDay     HistoricalPriceInterval = "7d"
	HistoricalPriceIntervalThirtyDay    HistoricalPriceInterval = "30d"
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
	Date  int64   `json:"d"`
	Close float64 `json:"c"`
	High  float64 `json:"h"`
	Low   float64 `json:"l"`
	Open  float64 `json:"o"`
}

type StockRestriction struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	StartDate   time.Time `json:"startDate"`
	EndDate     time.Time `json:"endDate"`
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

func (c *Client) GetCustomHistoricalPrices(ctx context.Context, symbol string, region Region, fromDate string, toDate string, interval HistoricalPriceInterval, detail bool) ([]PriceDataPoint, error) {
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

func (c *Client) GetTickRules(ctx context.Context, symbol string, region Region) (TickRule, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/api/v2/stock/rules", c.baseUrl), nil)
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
