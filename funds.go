package laplace

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
)

type FundType string

const (
	FundTypeStockUmbrella           FundType = "STOCK_UMBRELLA_FUND"
	FundTypeVariableUmbrella        FundType = "VARIABLE_UMBRELLA_FUND"
	FundTypeParticipationUmbrella   FundType = "PARTICIPATION_UMBRELLA_FUND"
	FundTypeFlexibleUmbrella        FundType = "FLEXIBLE_UMBRELLA_FUND"
	FundTypeFundBasketUmbrella      FundType = "FUND_BASKET_UMBRELLA_FUND"
	FundTypeMoneyMarketUmbrella     FundType = "MONEY_MARKET_UMBRELLA_FUND"
	FundTypePreciousMetalsUmbrella  FundType = "PRECIOUS_METALS_UMBRELLA_FUND"
	FundTypeDebtInstrumentsUmbrella FundType = "DEBT_INSTRUMENTS_UMBRELLA_FUND"
	FundTypeMixedUmbrella           FundType = "MIXED_UMBRELLA_FUND"
	FundTypeUnknown                 FundType = "UNKNOWN_FUND_TYPE"
)

type FundContentType string

const (
	FundContentTypeBistStock  FundContentType = "BIST_STOCK"
	FundContentTypeOtherStock FundContentType = "OTHER_STOCK"
	FundContentTypeUnknown    FundContentType = "UNKNOWN"
)

type FundAssetCategory string

const (
	FundAssetCategoryOther                                  FundAssetCategory = "OTHER"
	FundAssetCategoryEquity                                 FundAssetCategory = "EQUITY"
	FundAssetCategoryLiquidDeposit                          FundAssetCategory = "LIQUID_DEPOSIT"
	FundAssetCategoryFuturesCashCollateral                  FundAssetCategory = "FUTURES_CASH_COLLATERAL"
	FundAssetCategoryInvestmentFunds                        FundAssetCategory = "INVESTMENT_FUNDS"
	FundAssetCategoryParticipationAccount                   FundAssetCategory = "PARTICIPATION_ACCOUNT"
	FundAssetCategoryPreciousMetals                         FundAssetCategory = "PRECIOUS_METALS"
	FundAssetCategoryCorporateBond                          FundAssetCategory = "CORPORATE_BOND"
	FundAssetCategoryCurrency                               FundAssetCategory = "CURRENCY"
	FundAssetCategoryPublicExternalDebtSecurities           FundAssetCategory = "PUBLIC_EXTERNAL_DEBT_SECURITIES"
	FundAssetCategoryPrivateSectorExternalDebtSecurities    FundAssetCategory = "PRIVATE_SECTOR_EXTERNAL_DEBT_SECURITIES"
	FundAssetCategoryPublicLeaseCertificates                FundAssetCategory = "PUBLIC_LEASE_CERTIFICATES"
	FundAssetCategoryPrivateSectorLeaseCertificates         FundAssetCategory = "PRIVATE_SECTOR_LEASE_CERTIFICATES"
	FundAssetCategoryForeignExchangeTradedFunds             FundAssetCategory = "FOREIGN_EXCHANGE_TRADED_FUNDS"
	FundAssetCategoryPublicLeaseCertificatesCurrency        FundAssetCategory = "PUBLIC_LEASE_CERTIFICATES_CURRENCY"
	FundAssetCategoryGovernmentBond                         FundAssetCategory = "GOVERNMENT_BOND"
	FundAssetCategoryPrivateSectorLeaseCertificatesCurrency FundAssetCategory = "PRIVATE_SECTOR_LEASE_CERTIFICATES_CURRENCY"
	FundAssetCategoryUnknown                                FundAssetCategory = "UNKNOWN"
)

type HistoricalFundPricePeriod string

const (
	HistoricalFundPricePeriodOneWeek    HistoricalFundPricePeriod = "1H"
	HistoricalFundPricePeriodOneMonth   HistoricalFundPricePeriod = "1A"
	HistoricalFundPricePeriodThreeMonth HistoricalFundPricePeriod = "3A"
	HistoricalFundPricePeriodOneYear    HistoricalFundPricePeriod = "1Y"
	HistoricalFundPricePeriodThreeYear  HistoricalFundPricePeriod = "3Y"
	HistoricalFundPricePeriodFiveYear   HistoricalFundPricePeriod = "5Y"
)

type Fund struct {
	AssetType     AssetType `json:"assetType"`
	Name          string    `json:"name"`
	Symbol        string    `json:"symbol"`
	Active        bool      `json:"active"`
	ManagementFee float64   `json:"managementFee"`
	RiskLevel     int       `json:"riskLevel"`
	FundType      FundType  `json:"fundType"`
	OwnerSymbol   string    `json:"ownerSymbol"`
}

type FundStats struct {
	YearBeta         float64 `json:"yearBeta"`
	YearStdev        float64 `json:"yearStdev"`
	YtdReturn        float64 `json:"ytdReturn"`
	YearMomentum     float64 `json:"yearMomentum"`
	YearlyReturn     float64 `json:"yearlyReturn"`
	MonthlyReturn    float64 `json:"monthlyReturn"`
	FiveYearReturn   float64 `json:"fiveYearReturn"`
	SixMonthReturn   float64 `json:"sixMonthReturn"`
	ThreeYearReturn  float64 `json:"threeYearReturn"`
	ThreeMonthReturn float64 `json:"threeMonthReturn"`
}

type FundAsset struct {
	Type               FundContentType `json:"type"`
	Symbol             string          `json:"symbol"`
	WholePercentage    float64         `json:"wholePercentage"`
	CategoryPercentage float64         `json:"categoryPercentage"`
}

type FundDistribution struct {
	Categories []FundDistributionCategory `json:"categories"`
}

type FundDistributionCategory struct {
	Category   FundAssetCategory `json:"category"`
	Percentage float64           `json:"percentage"`
	Assets     []FundAsset       `json:"assets,omitempty"`
}

type FundHistoricalPrice struct {
	Aum           float64 `json:"aum"`
	Date          string  `json:"date"`
	Price         float64 `json:"price"`
	ShareCount    float64 `json:"shareCount"`
	InvestorCount int     `json:"investorCount"`
}

// GetFunds retrieves a paginated list of investment funds for the specified region with basic fund information.
func (c *Client) GetFunds(ctx context.Context, region Region, page int, pageSize int) ([]Fund, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/api/v1/fund", c.baseUrl), nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Add("region", string(region))
	q.Add("page", strconv.Itoa(page))
	q.Add("pageSize", strconv.Itoa(pageSize))
	req.URL.RawQuery = q.Encode()

	resp, err := sendRequest[[]Fund](ctx, c, req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// GetFundStats fetches comprehensive statistical data for a specific fund including returns, risk metrics, and performance indicators.
func (c *Client) GetFundStats(ctx context.Context, symbol string, region Region) (*FundStats, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/api/v1/fund/stats", c.baseUrl), nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Add("symbol", symbol)
	q.Add("region", string(region))
	req.URL.RawQuery = q.Encode()

	resp, err := sendRequest[FundStats](ctx, c, req)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

// GetFundDistribution retrieves detailed asset allocation and distribution information for a specific fund.
func (c *Client) GetFundDistribution(ctx context.Context, symbol string, region Region) (*FundDistribution, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/api/v1/fund/distribution", c.baseUrl), nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Add("symbol", symbol)
	q.Add("region", string(region))
	req.URL.RawQuery = q.Encode()

	resp, err := sendRequest[FundDistribution](ctx, c, req)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

// GetHistoricalFundPrices retrieves historical price data for a fund over the specified time period.
func (c *Client) GetHistoricalFundPrices(ctx context.Context, symbol string, region Region, period HistoricalFundPricePeriod) ([]FundHistoricalPrice, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/api/v1/fund/price", c.baseUrl), nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Add("symbol", symbol)
	q.Add("region", string(region))
	q.Add("period", string(period))
	req.URL.RawQuery = q.Encode()

	resp, err := sendRequest[[]FundHistoricalPrice](ctx, c, req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
