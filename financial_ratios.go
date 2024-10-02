package laplace

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/samber/lo"
)

type StockSectorFinancialRatioComparison struct {
	MetricName      string                                      `json:"metric_name"`
	NormalizedValue float64                                     `json:"normalizedValue"`
	Details         []StockSectorFinancialRatioComparisonDetail `json:"details"`
}

type StockSectorFinancialRatioComparisonDetail struct {
	Slug          string  `json:"slug"`
	Value         float64 `json:"value"`
	SectorAverage float64 `json:"sectorAverage"`
}

func (c *Client) GetFinancialRatioComparison(ctx context.Context, symbol string, region Region) ([]StockSectorFinancialRatioComparison, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/api/v1/stock/financial-ratio-comparison", c.baseUrl), nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Add("symbol", symbol)
	q.Add("region", string(region))
	req.URL.RawQuery = q.Encode()

	resp, err := sendRequest[[]StockSectorFinancialRatioComparison](ctx, c, req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

type StockHistoricalRatios struct {
	Symbol     string                                     `json:"symbol"`
	Data       []StockHistoricalRatiosData                `json:"data"`
	Formatting map[string]StockHistoricalRatiosFormatting `json:"formatting"`
}

type StockHistoricalRatiosData struct {
	FiscalYear    int                                   `json:"fiscalYear"`
	FiscalQuarter int                                   `json:"fiscalQuarter"`
	Values        map[string]StockHistoricalRatiosValue `json:"values"`
}

type StockHistoricalRatiosValue struct {
	Value         float64 `json:"value"`
	SectorAverage float64 `json:"sectorAverage"`
}

type StockHistoricalRatiosFormatting struct {
	Name        string  `json:"name"`
	Slug        string  `json:"slug"`
	Precision   int     `json:"precision"`
	Multiplier  float64 `json:"multiplier"`
	Suffix      string  `json:"suffix"`
	Prefix      string  `json:"prefix"`
	Interval    string  `json:"interval"`
	Description string  `json:"description"`
}

type HistoricalRatiosKey string

const (
	HistoricalRatiosKeyPriceToEarningsRatio HistoricalRatiosKey = "pe-ratio"
	HistoricalRatiosKeyReturnOnEquity       HistoricalRatiosKey = "roe"
	HistoricalRatiosKeyReturnOnAssets       HistoricalRatiosKey = "roa"
	HistoricalRatiosKeyReturnOnCapital      HistoricalRatiosKey = "roic"
)

func (c *Client) GetHistoricalRatios(ctx context.Context, symbol string, keys []HistoricalRatiosKey, region Region) (StockHistoricalRatios, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/api/v1/stock/historical-ratios", c.baseUrl), nil)
	if err != nil {
		return StockHistoricalRatios{}, err
	}

	q := req.URL.Query()
	q.Add("symbol", symbol)
	q.Add("region", string(region))
	q.Add("slugs", strings.Join(lo.Map(keys, func(key HistoricalRatiosKey, _ int) string {
		return string(key)
	}), ","))
	req.URL.RawQuery = q.Encode()

	resp, err := sendRequest[StockHistoricalRatios](ctx, c, req)
	if err != nil {
		return StockHistoricalRatios{}, err
	}

	return resp, nil
}

type StockHistoricalRatiosDescription struct {
	Slug        string  `json:"slug"`
	Name        string  `json:"name"`
	Suffix      string  `json:"suffix"`
	Prefix      string  `json:"prefix"`
	Display     bool    `json:"display"`
	Precision   int     `json:"precision"`
	Multiplier  float64 `json:"multiplier"`
	Description string  `json:"description"`
	Interval    string  `json:"interval"`
}

func (c *Client) GetHistoricalRatiosDescriptions(ctx context.Context, locale Locale, region Region) ([]StockHistoricalRatiosDescription, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/api/v1/stock/historical-ratios/descriptions", c.baseUrl), nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Add("locale", string(locale))
	q.Add("region", string(region))

	req.URL.RawQuery = q.Encode()

	resp, err := sendRequest[[]StockHistoricalRatiosDescription](ctx, c, req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

type HistoricalFinancialSheets struct {
	Sheets []HistoricalFinancialSheet `json:"sheets"`
}

type HistoricalFinancialSheet struct {
	Period string                        `json:"period"`
	Rows   []HistoricalFinancialSheetRow `json:"rows"`
}

type HistoricalFinancialSheetRow struct {
	Description             string  `json:"description"`
	Value                   float64 `json:"value"`
	LineCodeId              int     `json:"lineCodeId"`
	IndentLevel             int     `json:"indentLevel"`
	FirstAncestorLineCodeId int     `json:"firstAncestorLineCodeId"`
	SectionLineCodeId       int     `json:"sectionLineCodeId"`
}

type FinancialSheetType string

const (
	FinancialSheetIncomeStatement FinancialSheetType = "incomeStatement"
	FinancialSheetBalanceSheet    FinancialSheetType = "balanceSheet"
	FinancialSheetCashFlow        FinancialSheetType = "cashFlowStatement"
)

type FinancialSheetPeriod string

const (
	FinancialSheetPeriodAnnual     FinancialSheetPeriod = "annual"
	FinancialSheetPeriodQuarterly  FinancialSheetPeriod = "quarterly"
	FinancialSheetPeriodCumulative FinancialSheetPeriod = "cumulative"
)

type Currency string

const (
	CurrencyUSD Currency = "USD"
	CurrencyTRY Currency = "TRY"
	CurrencyEUR Currency = "EUR"
)

type FinancialSheetDate struct {
	Day   int `json:"day"`
	Month int `json:"month"`
	Year  int `json:"year"`
}

func (c *Client) GetHistoricalFinancialSheets(ctx context.Context, symbol string, from FinancialSheetDate, to FinancialSheetDate, sheetType FinancialSheetType, period FinancialSheetPeriod, currency Currency, region Region) (HistoricalFinancialSheets, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/api/v1/stock/historical-financial-sheets", c.baseUrl), nil)
	if err != nil {
		return HistoricalFinancialSheets{}, err
	}

	q := req.URL.Query()
	q.Add("symbol", symbol)
	q.Add("from", fmt.Sprintf("%04d-%02d-%02d", from.Year, from.Month, from.Day))
	q.Add("to", fmt.Sprintf("%04d-%02d-%02d", to.Year, to.Month, to.Day))
	q.Add("sheetType", string(sheetType))
	q.Add("periodType", string(period))
	q.Add("currency", string(currency))
	q.Add("region", string(region))
	req.URL.RawQuery = q.Encode()

	resp, err := sendRequest[HistoricalFinancialSheets](ctx, c, req)
	if err != nil {
		return HistoricalFinancialSheets{}, err
	}

	return resp, nil
}
