package client

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AssetType string

const (
	AssetTypeStock     AssetType = "stock"
	AssetTypeForex     AssetType = "forex"
	AssetTypeIndex     AssetType = "index"
	AssetTypeEtf       AssetType = "etf"
	AssetTypeCommodity AssetType = "commodity"
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
}

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

func (c *Client) GetAllStocks(ctx context.Context, region Region) ([]Stock, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/api/v1/stock/all", c.baseUrl), nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Add("region", string(region))
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
