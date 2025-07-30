# Laplace Go SDK

[![Go Version](https://img.shields.io/badge/Go-1.21+-blue.svg)](https://golang.org/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)
[![Go Report Card](https://goreportcard.com/badge/github.com/Laplace-Analytics/laplace-api-golang)](https://goreportcard.com/report/github.com/Laplace-Analytics/laplace-api-golang)

The official Go SDK for the Laplace stock data platform. Get easy access to stock data, collections, financials, funds, and AI-powered insights.

## Features

- 🚀 **Easy to use**: Simple, intuitive API with Go idioms
- 📊 **Comprehensive data**: Stocks, collections, financials, funds, and AI insights
- 🔧 **Well-typed**: Full Go type safety with comprehensive structs
- 🧪 **Well-tested**: Comprehensive test coverage with real API integration
- 🌍 **Multi-region**: Support for US and Turkish markets
- ⚡ **Fast**: Built on Go's high-performance HTTP client
- 📚 **Well-documented**: Complete GoDoc documentation for all public methods

## Installation

```bash
go get github.com/Laplace-Analytics/laplace-api-golang
```

## Quick Start

```go
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/Laplace-Analytics/laplace-api-golang"
)

func main() {
	// Initialize the client
	client := laplace.NewClient(laplace.LaplaceConfiguration{
		APIKey:  "your-api-key-here",
	})

	// Create a context
	ctx := context.Background()

	// Get stock details
	stock, err := client.GetStockDetailBySymbol(ctx, "AAPL", laplace.AssetClassEquity, laplace.RegionUs, laplace.LocaleEn)
	if err != nil {
		log.Fatal("Error fetching stock:", err)
	}
	fmt.Printf("%s: %s\n", stock.Name, stock.Description)

	// Get all stocks in a region
	stocks, err := client.GetAllStocks(ctx, laplace.RegionUs, 1, 10)
	if err != nil {
		log.Fatal("Error fetching stocks:", err)
	}
	for _, stock := range stocks {
		fmt.Printf("%s: %s\n", stock.Symbol, stock.Name)
	}

	// Get collections
	collections, err := client.GetAllCollections(ctx, laplace.RegionTr, laplace.LocaleEn)
	if err != nil {
		log.Fatal("Error fetching collections:", err)
	}
	for _, collection := range collections {
		fmt.Printf("%s: %d stocks\n", collection.Title, collection.NumStocks)
	}

	// Get collection details
	collectionDetail, err := client.GetCollectionDetail(ctx, "620f455a0187ade00bb0d55f", laplace.RegionTr, laplace.LocaleEn)
	if err != nil {
		log.Fatal("Error fetching collection detail:", err)
	}
	fmt.Printf("Stocks in %s:\n", collectionDetail.Title)
	for _, stock := range collectionDetail.Stocks {
		fmt.Printf("  %s: %s\n", stock.Symbol, stock.Name)
	}
}
```

## API Reference

### Stocks Client

```go
// Get all stocks with pagination
stocks, err := client.GetAllStocks(ctx, laplace.RegionUs, 1, 10)

// Get stock detail by symbol
stock, err := client.GetStockDetailBySymbol(ctx, "AAPL", laplace.AssetClassEquity, laplace.RegionUs, laplace.LocaleEn)

// Get stock detail by ID
stock, err := client.GetStockDetailByID(ctx, "stock-id", laplace.LocaleEn)

// Get historical prices
prices, err := client.GetHistoricalPrices(ctx, []string{"AAPL", "GOOGL"}, laplace.RegionUs, []laplace.HistoricalPricePeriod{laplace.HistoricalPricePeriodOneDay, laplace.HistoricalPricePeriodOneWeek})

// Get historical prices with custom interval
prices, err := client.GetCustomHistoricalPrices(ctx, "AAPL", laplace.RegionUs, "2024-01-01", "2024-01-31", laplace.HistoricalPriceIntervalOneMinute, true)

// Get tick rules (Turkey only)
rules, err := client.GetTickRules(ctx, "THYAO", laplace.RegionTr)

// Get restrictions (Turkey only)
restrictions, err := client.GetStockRestrictions(ctx, "THYAO", laplace.RegionTr)
```

### Collections Client

```go
// Get all collections
collections, err := client.GetAllCollections(ctx, laplace.RegionTr, laplace.LocaleEn)

// Get collection detail
detail, err := client.GetCollectionDetail(ctx, "collection-id", laplace.RegionTr, laplace.LocaleEn)

// Get themes
themes, err := client.GetAllThemes(ctx, laplace.RegionTr, laplace.LocaleEn)

// Get theme detail
themeDetail, err := client.GetThemeDetail(ctx, "theme-id", laplace.RegionTr, laplace.LocaleEn)

// Get industries
industries, err := client.GetAllIndustries(ctx, laplace.RegionTr, laplace.LocaleEn)

// Get industry detail
industryDetail, err := client.GetIndustryDetail(ctx, "industry-id", laplace.RegionTr, laplace.LocaleEn)

// Get sectors
sectors, err := client.GetAllSectors(ctx, laplace.RegionTr, laplace.LocaleEn)

// Get sector detail
sectorDetail, err := client.GetSectorDetail(ctx, "sector-id", laplace.RegionTr, laplace.LocaleEn)
```

### Funds Client

```go
// Get all funds
funds, err := client.GetFunds(ctx, laplace.RegionTr, 1, 10)

// Get fund statistics
stats, err := client.GetFundStats(ctx, "fund-symbol", laplace.RegionTr)

// Get fund distribution
distribution, err := client.GetFundDistribution(ctx, "fund-symbol", laplace.RegionTr)

// Get historical fund prices
prices, err := client.GetHistoricalFundPrices(ctx, "fund-symbol", laplace.RegionTr, laplace.HistoricalFundPricePeriodOneYear)
```

### Financial Data Client

```go
// Get financial ratios
ratios, err := client.GetHistoricalRatios(ctx, "AAPL", []laplace.HistoricalRatiosKey{laplace.HistoricalRatiosKeyPERatio}, laplace.RegionUs)

// Get financial ratio comparisons
comparisons, err := client.GetFinancialRatioComparison(ctx, "AAPL", laplace.RegionUs, laplace.PeerTypeSector)

// Get financial statements
statements, err := client.GetHistoricalFinancialSheets(ctx, "AAPL", laplace.FinancialSheetDate{Year: 2024, Month: 1, Day: 1}, laplace.FinancialSheetDate{Year: 2024, Month: 12, Day: 31}, laplace.FinancialSheetIncomeStatement, laplace.FinancialSheetPeriodAnnual, laplace.CurrencyUSD, laplace.RegionUs)

// Get stock dividends
dividends, err := client.GetStockDividends(ctx, "AAPL", laplace.RegionUs)

// Get stock statistics
stats, err := client.GetStockStats(ctx, []string{"AAPL", "GOOGL"}, laplace.RegionUs)

// Get top movers
movers, err := client.GetTopMovers(ctx, laplace.TopMoversDirectionGainers, laplace.AssetClassEquity, laplace.AssetTypeStock, 1, 10, laplace.RegionUs)
```

### Live Price Client

```go
// Get live prices for BIST stocks
lc, err := client.GetLivePriceForBIST(ctx, []string{"THYAO", "GARAN"})

// Get live prices for US stocks
lc, err := client.GetLivePriceForUS(ctx, []string{"AAPL", "GOOGL"})


for data := range lc.Receive() {
	fmt.Printf("Received data: %+v\n", data.Data)
}

```

### Brokers Client

```go
// Get all brokers
brokers, err := client.GetBrokers(ctx, laplace.RegionTr, 1, 10)

// Get market stocks with broker statistics
marketStocks, err := client.GetMarketStocks(ctx, laplace.RegionTr, laplace.BrokerSortNetAmount, laplace.BrokerSortDirectionDesc, "2024-01-01", "2024-01-31", 1, 10)

// Get brokers by stock
brokersByStock, err := client.GetBrokersByStock(ctx, "THYAO", laplace.RegionTr, laplace.BrokerSortNetAmount, laplace.BrokerSortDirectionDesc, "2024-01-01", "2024-01-31", 1, 10)
```

### Search Client

```go
// Search across stocks, collections, sectors, and industries
results, err := client.Search(ctx, "technology", []laplace.SearchType{laplace.SearchTypeStock, laplace.SearchTypeCollection}, laplace.RegionUs, laplace.LocaleEn, 1, laplace.PageSize20)
```

### WebSocket Client

```go
// Get WebSocket URL for real-time data
url, err := client.GetWebSocketUrl(ctx, "user-id", []laplace.FeedType{laplace.FeedTypeLivePriceTR}, laplace.RegionTr)

// Update user details
err = client.UpdateUserDetails(ctx, laplace.UpdateUserDetailsParams{
	ExternalUserID: "user-id",
	FirstName:      "John",
	LastName:       "Doe",
	Address:        "123 Main St",
	City:           "New York",
	CountryCode:    "US",
	AccessorType:   laplace.AccessorTypeUser,
	Active:         true,
})
```

### Capital Increase Client

```go
// Get all capital increases
increases, err := client.GetAllCapitalIncreases(ctx, 1, 10, laplace.RegionTr)

// Get capital increases for a specific instrument
instrumentIncreases, err := client.GetCapitalIncreasesForInstrument(ctx, "THYAO", 1, 10, laplace.RegionTr)

// Get active rights for an instrument
rights, err := client.GetActiveRightsForInstrument(ctx, "THYAO", "2024-01-15", laplace.RegionTr)
```

### Custom Themes Client

```go
// Get all custom themes
themes, err := client.GetAllCustomThemes(ctx, laplace.LocaleEn)

// Get custom theme detail
themeDetail, err := client.GetCustomThemeDetail(ctx, "theme-id", laplace.RegionTr, laplace.LocaleEn, laplace.SortByPriceChange)

// Create a custom theme
id, err := client.CreateCustomTheme(ctx, laplace.CreateCustomThemeParams{
	Title:       laplace.LocaleString{"en": "My Tech Portfolio"},
	Description: laplace.LocaleString{"en": "Technology stocks portfolio"},
	Region:      []laplace.Region{laplace.RegionUs},
	Stocks:      []primitive.ObjectID{/* stock IDs */},
	Status:      laplace.CollectionStatusActive,
})

// Update a custom theme
err = client.UpdateCustomTheme(ctx, *id, laplace.UpdateCustomThemeParams{
	Title: laplace.LocaleString{"en": "Updated Tech Portfolio"},
	Stocks: []primitive.ObjectID{/* updated stock IDs */},
})

// Delete a custom theme
err = client.DeleteCustomTheme(ctx, *id)
```

### Key Insights Client

```go
// Get key insights for a stock
insights, err := client.GetKeyInsights(ctx, "AAPL", laplace.RegionUs)
```

## Supported Regions

- **US**: United States stock market
- **TR**: Turkey stock market (Borsa Istanbul)

## Error Handling

```go
import (
	"fmt"
	"github.com/Laplace-Analytics/laplace-api-golang"
)

client := laplace.NewClient(laplace.LaplaceConfiguration{
	APIKey:  "your-api-key",
})

ctx := context.Background()

stock, err := client.GetStockDetailBySymbol(ctx, "INVALID", laplace.AssetClassEquity, laplace.RegionUs, laplace.LocaleEn)
if err != nil {
	if laplaceErr, ok := err.(*laplace.LaplaceHTTPError); ok {
		fmt.Printf("API Error: %s\n", laplaceErr.Message)
		fmt.Printf("Status Code: %d\n", laplaceErr.HTTPStatus)
	} else {
		fmt.Printf("Error: %v\n", err)
	}
	return
}
```

## Authentication

Get your API key from the Laplace platform and initialize the client:

```go
client := laplace.NewClient(laplace.LaplaceConfiguration{
	APIKey:  "your-api-key-here",
})
```

## Configuration

You can also load configuration from environment variables:

```go
config, err := laplace.LoadGlobal("")
if err != nil {
	log.Fatal("Error loading config:", err)
}

client := laplace.NewClient(*config)
```

Environment variables:

- `LAPLACE_API_KEY`: Your API key
- `LAPLACE_BASE_URL`: API base URL

## Development

### Setup

```bash
git clone https://github.com/Laplace-Analytics/laplace-api-golang.git
cd laplace-api-golang
go mod download
```

### Running Tests

```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Run integration tests (requires API key)
LAPLACE_API_KEY=your-key go test -tags=integration ./...
```

## Requirements

- Go 1.21+
- Standard library only (no external dependencies)

## Documentation

Full API documentation is available at [laplace.finfree.co/en/docs](https://laplace.finfree.co/en/docs)

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
