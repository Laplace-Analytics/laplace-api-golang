package laplace

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/samber/lo"
)

type StockPeerFinancialRatioComparison struct {
	MetricName      string                                  `json:"metricName"`
	NormalizedValue float64                                 `json:"normalizedValue"`
	Data            []StockPeerFinancialRatioComparisonData `json:"data"`
}

type StockPeerFinancialRatioComparisonData struct {
	Slug    string  `json:"slug"`
	Value   float64 `json:"value"`
	Average float64 `json:"average"`
}

type PeerType string

const (
	PeerTypeSector   PeerType = "sector"
	PeerTypeIndustry PeerType = "industry"
)

// GetFinancialRatioComparison retrieves financial ratio comparisons for a stock against its sector or industry peers.
func (c *Client) GetFinancialRatioComparison(ctx context.Context, symbol string, region Region, peerType PeerType) ([]StockPeerFinancialRatioComparison, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/api/v2/stock/financial-ratio-comparison", c.baseUrl), nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Add("symbol", symbol)
	q.Add("region", string(region))
	q.Add("peerType", string(peerType))
	req.URL.RawQuery = q.Encode()

	resp, err := sendRequest[[]StockPeerFinancialRatioComparison](ctx, c, req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

type StockHistoricalRatios struct {
	Items            []StockHistoricalRatiosData `json:"items"`
	FinalValue       float64                     `json:"finalValue"`
	ThreeYearGrowth  float64                     `json:"threeYearGrowth"`
	YearGrowth       float64                     `json:"yearGrowth"`
	FinalSectorValue float64                     `json:"finalSectorValue"`
	Slug             string                      `json:"slug"`
	Currency         string                      `json:"currency"`
	Format           string                      `json:"format"`
	Name             string                      `json:"name"`
}

type StockHistoricalRatiosData struct {
	Period     string  `json:"period"`
	Value      float64 `json:"value"`
	SectorMean float64 `json:"sectorMean"`
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
	HistoricalRatiosKeyGrossMargin                          HistoricalRatiosKey = "gross-margin"
	HistoricalRatiosKeyEBITDA                               HistoricalRatiosKey = "ebitda"
	HistoricalRatiosKeyPERatio                              HistoricalRatiosKey = "pe-ratio"
	HistoricalRatiosKeyOperatingMargin                      HistoricalRatiosKey = "favok_marji"
	HistoricalRatiosKeyFreeCashFlowGrowth                   HistoricalRatiosKey = "serbest_nakit_akisi_buyumesi"
	HistoricalRatiosKeyDaysPayable                          HistoricalRatiosKey = "days-payable"
	HistoricalRatiosKeyInventoryTurnover                    HistoricalRatiosKey = "inventory-turnover"
	HistoricalRatiosKeyDepositGrowth                        HistoricalRatiosKey = "mevduat_buyumesi"
	HistoricalRatiosKeyNetInterestMargin                    HistoricalRatiosKey = "net_faiz_marji"
	HistoricalRatiosKeyClaimPaymentsGrowth                  HistoricalRatiosKey = "gerceklesen_tazminatlar_buyumesi"
	HistoricalRatiosKeyClaimsPerPremiumRatio                HistoricalRatiosKey = "prim_basina_tazminat_orani"
	HistoricalRatiosKeyEVToOCF                              HistoricalRatiosKey = "evOcf"
	HistoricalRatiosKeyEVToIC                               HistoricalRatiosKey = "evic"
	HistoricalRatiosKeyEBT                                  HistoricalRatiosKey = "ebt"
	HistoricalRatiosKeyCAPEX                                HistoricalRatiosKey = "capex"
	HistoricalRatiosKeyFinancialInvestments                 HistoricalRatiosKey = "financial_investments"
	HistoricalRatiosKeyRealtimeEPSBasic                     HistoricalRatiosKey = "realtime_eps-basic"
	HistoricalRatiosKeyQuickRatio                           HistoricalRatiosKey = "quick-ratio"
	HistoricalRatiosKeyEVToEBITDA                           HistoricalRatiosKey = "ev-to-ebitda"
	HistoricalRatiosKeyROCE                                 HistoricalRatiosKey = "roce"
	HistoricalRatiosKeyROIC                                 HistoricalRatiosKey = "roic"
	HistoricalRatiosKeyROA                                  HistoricalRatiosKey = "roa"
	HistoricalRatiosKeyDaysSalesOutstanding                 HistoricalRatiosKey = "days-sales-outstanding"
	HistoricalRatiosKeyLoanToAssetRatio                     HistoricalRatiosKey = "kredi_aktif_orani"
	HistoricalRatiosKeyLoanToDepositRatio                   HistoricalRatiosKey = "kredi_mevduat_orani"
	HistoricalRatiosKeyTechnicalProfitGrowth                HistoricalRatiosKey = "teknik_kar_buyumesi"
	HistoricalRatiosKeyNetPremiumEarnedGrowth               HistoricalRatiosKey = "net_kazanilan_prim_buyumesi"
	HistoricalRatiosKeyEBITGrowth                           HistoricalRatiosKey = "ebitGrowth"
	HistoricalRatiosKeyCROIC                                HistoricalRatiosKey = "croic"
	HistoricalRatiosKeyRealtimeMarketValue                  HistoricalRatiosKey = "realtime_piyasa_degeri"
	HistoricalRatiosKeyRealtimePBRatio                      HistoricalRatiosKey = "realtime_pb-ratio"
	HistoricalRatiosKeyRealtimePERatio                      HistoricalRatiosKey = "realtime_pe-ratio"
	HistoricalRatiosKeyCurrentRatio                         HistoricalRatiosKey = "current-ratio"
	HistoricalRatiosKeyDaysInventory                        HistoricalRatiosKey = "days-inventory"
	HistoricalRatiosKeyNetMargin                            HistoricalRatiosKey = "net-margin"
	HistoricalRatiosKeySalesGrowth                          HistoricalRatiosKey = "satis_buyumesi"
	HistoricalRatiosKeyROE                                  HistoricalRatiosKey = "roe"
	HistoricalRatiosKeyAssetTurnover                        HistoricalRatiosKey = "asset-turnover"
	HistoricalRatiosKeyLeverageRatio                        HistoricalRatiosKey = "leverage-ratio"
	HistoricalRatiosKeySales                                HistoricalRatiosKey = "satislar"
	HistoricalRatiosKeyNetProfit                            HistoricalRatiosKey = "net_kar"
	HistoricalRatiosKeyInterestCoverage                     HistoricalRatiosKey = "interestCoverage"
	HistoricalRatiosKeyTotalOperationalExpense              HistoricalRatiosKey = "total_operational_expense"
	HistoricalRatiosKeyTotalOperationalExpenseToGrossProfit HistoricalRatiosKey = "total_operational_expense_gross_profit_ratio"
	HistoricalRatiosKeyCashAndCashEquivalents               HistoricalRatiosKey = "cash_and_cash_equivalents"
	HistoricalRatiosKeyCashToAssets                         HistoricalRatiosKey = "cash_to_assets"
	HistoricalRatiosKeyCAPEXToNetProfit                     HistoricalRatiosKey = "capex_to_net_profit"
	HistoricalRatiosKeyRealtimeEVToEBITDA                   HistoricalRatiosKey = "realtime_ev-to-ebitda"
	HistoricalRatiosKeyReceivablesTurnover                  HistoricalRatiosKey = "alacak_devir_hizi"
	HistoricalRatiosKeyEPSBasic                             HistoricalRatiosKey = "eps-basic"
	HistoricalRatiosKeyNetProfitGrowth                      HistoricalRatiosKey = "net_kar_buyumesi"
	HistoricalRatiosKeyDebtToEquity                         HistoricalRatiosKey = "debt-to-equity"
	HistoricalRatiosKeyNetDebtToEBITDA                      HistoricalRatiosKey = "net_borc_favok"
	HistoricalRatiosKeyPBRatio                              HistoricalRatiosKey = "pb-ratio"
	HistoricalRatiosKeyEBITDAGrowth                         HistoricalRatiosKey = "favok_buyumesi"
	HistoricalRatiosKeyCashConversionCycle                  HistoricalRatiosKey = "cash-conversion-cycle"
	HistoricalRatiosKeyGrossProfitGrowth                    HistoricalRatiosKey = "brut_kar_buyumesi"
	HistoricalRatiosKeyLoanGrowth                           HistoricalRatiosKey = "kredi_buyumesi"
	HistoricalRatiosKeyGrossWrittenPremiumGrowth            HistoricalRatiosKey = "brut_yazilan_prim_buyumesi"
	HistoricalRatiosKeyTechnicalProfitMargin                HistoricalRatiosKey = "teknik_kar_marji"
	HistoricalRatiosKeyCompanyPremiumRetentionRatio         HistoricalRatiosKey = "sirketin_prim_tutma_orani"
	HistoricalRatiosKeyMarketValue                          HistoricalRatiosKey = "piyasa_degeri"
	HistoricalRatiosKeyFinancialExpensesToEBITRatio         HistoricalRatiosKey = "financial_expenses_ebit_ratio"
	HistoricalRatiosKeyShortTermToLongTermObligations       HistoricalRatiosKey = "short_term_obligations_long_term_obligations"
	HistoricalRatiosKeyRetainedEarnings                     HistoricalRatiosKey = "retained_earnings"
	HistoricalRatiosKeyThreeYearCAGRFreeCashFlow            HistoricalRatiosKey = "three_year_cagr_free_cash_flow"
	HistoricalRatiosKeyPOE                                  HistoricalRatiosKey = "poe"
	HistoricalRatiosKeyLongTermLoansToPeriodProfit          HistoricalRatiosKey = "long_term_loans_period_profit_ratio"
	HistoricalRatiosKeyLongTermLoans                        HistoricalRatiosKey = "long_term_loans"
	HistoricalRatiosKeyCommercialReceivablesToCurrentAssets HistoricalRatiosKey = "commercial_receivables_total_current_assets"
	HistoricalRatiosKeyStockGrowth                          HistoricalRatiosKey = "stock_growth"
	HistoricalRatiosKeyFiveYearRetainedEarningsChange       HistoricalRatiosKey = "five_year_retained_earnings_change"
	HistoricalRatiosKeyThreeYearCAGRRetainedEarnings        HistoricalRatiosKey = "three_year_cagr_retained_earnings"
	HistoricalRatiosKeyPOCF                                 HistoricalRatiosKey = "pocf"
	HistoricalRatiosKeyFCFToEV                              HistoricalRatiosKey = "fcfEv"
	HistoricalRatiosKeyDD                                   HistoricalRatiosKey = "dd"
	HistoricalRatiosKeyNetDebt                              HistoricalRatiosKey = "net_borc"
	HistoricalRatiosKeyPaidInCapital                        HistoricalRatiosKey = "odenmis_sermaye"
)

// GetHistoricalRatios fetches historical financial ratios for a stock over time with sector comparisons.
func (c *Client) GetHistoricalRatios(ctx context.Context, symbol string, keys []HistoricalRatiosKey, region Region) ([]StockHistoricalRatios, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/api/v2/stock/historical-ratios", c.baseUrl), nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Add("symbol", symbol)
	q.Add("region", string(region))
	q.Add("slugs", strings.Join(lo.Map(keys, func(key HistoricalRatiosKey, _ int) string {
		return string(key)
	}), ","))
	req.URL.RawQuery = q.Encode()

	resp, err := sendRequest[[]StockHistoricalRatios](ctx, c, req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

type StockHistoricalRatiosDescription struct {
	ID          int       `json:"id"`
	Format      string    `json:"format"`
	Currency    string    `json:"currency"`
	Slug        string    `json:"slug"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Locale      string    `json:"locale"`
	IsRealtime  bool      `json:"isRealtime"`
}

// GetHistoricalRatiosDescriptions retrieves metadata and descriptions for available historical financial ratios.
func (c *Client) GetHistoricalRatiosDescriptions(ctx context.Context, locale Locale, region Region) ([]StockHistoricalRatiosDescription, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/api/v2/stock/historical-ratios/descriptions", c.baseUrl), nil)
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
	Period string      `json:"period"`
	Items  []SheetItem `json:"items"`
}

type SheetItem struct {
	Description     string  `json:"description"`
	Value           float64 `json:"value"`
	SheetLineItemId int     `json:"lineCodeId"`
	Indent          int     `json:"indentLevel"`
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

// GetHistoricalFinancialSheets fetches historical financial statements (income statement, balance sheet, cash flow) for a stock.
func (c *Client) GetHistoricalFinancialSheets(ctx context.Context, symbol string, from FinancialSheetDate, to FinancialSheetDate, sheetType FinancialSheetType, period FinancialSheetPeriod, currency Currency, region Region) (HistoricalFinancialSheets, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/api/v3/stock/historical-financial-sheets", c.baseUrl), nil)
	if err != nil {
		return HistoricalFinancialSheets{}, err
	}

	if sheetType == FinancialSheetBalanceSheet && period != FinancialSheetPeriodCumulative {
		return HistoricalFinancialSheets{}, errors.New("balance sheet is only available for cumulative period")
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
