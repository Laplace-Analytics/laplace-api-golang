package laplace

import (
	"context"
	"testing"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestReadme(t *testing.T) {
	cfg, err := loadConfig()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	client, err := NewClient(*cfg)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	ctx := context.Background()

	// Test Stocks Client methods
	t.Run("Stocks Client", func(t *testing.T) {
		// Get all stocks with pagination
		_, err := client.GetAllStocks(ctx, RegionUs, 1, 10)
		if err != nil {
			t.Errorf("GetAllStocks failed: %v", err)
		}

		// Get stock detail by symbol
		_, err = client.GetStockDetailBySymbol(ctx, "AAPL", AssetClassEquity, RegionUs, LocaleEn)
		if err != nil {
			t.Errorf("GetStockDetailBySymbol failed: %v", err)
		}

		// Get stock detail by ID (using a valid ID format)
		_, err = client.GetStockDetailByID(ctx, "648ab66e38daf3102a5a7401", LocaleEn)
		if err != nil {
			t.Errorf("GetStockDetailByID failed: %v", err)
		}

		// Get historical prices
		_, err = client.GetHistoricalPrices(ctx, []string{"THYAO"}, RegionTr, []HistoricalPricePeriod{HistoricalPricePeriodOneDay, HistoricalPricePeriodOneWeek})
		if err != nil {
			t.Errorf("GetHistoricalPrices failed: %v", err)
		}

		// Get custom historical prices
		_, err = client.GetCustomHistoricalPrices(ctx, "THYAO", RegionTr, "2024-01-01", "2024-01-10", HistoricalPriceIntervalOneMinute, true)
		if err != nil {
			t.Errorf("GetCustomHistoricalPrices failed: %v", err)
		}

		// Get tick rules (Turkey only)
		_, err = client.GetTickRules(ctx, "THYAO", RegionTr)
		if err != nil {
			t.Errorf("GetTickRules failed: %v", err)
		}

		// Get restrictions (Turkey only)
		_, err = client.GetStockRestrictions(ctx, "THYAO", RegionTr)
		if err != nil {
			t.Errorf("GetStockRestrictions failed: %v", err)
		}
	})

	// Test Collections Client methods
	t.Run("Collections Client", func(t *testing.T) {
		// Get all collections
		_, err := client.GetAllCollections(ctx, RegionTr, LocaleEn)
		if err != nil {
			t.Errorf("GetAllCollections failed: %v", err)
		}

		// Get collection detail
		_, err = client.GetCollectionDetail(ctx, "620f455a0187ade00bb0d55f", RegionTr, LocaleEn)
		if err != nil {
			t.Errorf("GetCollectionDetail failed: %v", err)
		}

		// Get themes
		_, err = client.GetAllThemes(ctx, RegionTr, LocaleEn)
		if err != nil {
			t.Errorf("GetAllThemes failed: %v", err)
		}

		// Get theme detail
		_, err = client.GetThemeDetail(ctx, "620f455a0187ade00bb0d55f", RegionTr, LocaleEn)
		if err != nil {
			t.Errorf("GetThemeDetail failed: %v", err)
		}

		// Get industries
		_, err = client.GetAllIndustries(ctx, RegionTr, LocaleEn)
		if err != nil {
			t.Errorf("GetAllIndustries failed: %v", err)
		}

		// Get industry detail
		_, err = client.GetIndustryDetail(ctx, "65533e441fa5c7b58afa0957", RegionTr, LocaleEn)
		if err != nil {
			t.Errorf("GetIndustryDetail failed: %v", err)
		}

		// Get sectors
		_, err = client.GetAllSectors(ctx, RegionTr, LocaleEn)
		if err != nil {
			t.Errorf("GetAllSectors failed: %v", err)
		}

		// Get sector detail
		_, err = client.GetSectorDetail(ctx, "65533e047844ee7afe9941bf", RegionTr, LocaleEn)
		if err != nil {
			t.Errorf("GetSectorDetail failed: %v", err)
		}
	})

	// Test Funds Client methods
	t.Run("Funds Client", func(t *testing.T) {
		// Get all funds
		_, err := client.GetFunds(ctx, RegionTr, 1, 10)
		if err != nil {
			t.Errorf("GetFunds failed: %v", err)
		}

		// Get fund statistics
		_, err = client.GetFundStats(ctx, "fund-symbol", RegionTr)
		if err != nil {
			t.Errorf("GetFundStats failed: %v", err)
		}

		// Get fund distribution
		_, err = client.GetFundDistribution(ctx, "fund-symbol", RegionTr)
		if err != nil {
			t.Errorf("GetFundDistribution failed: %v", err)
		}

		// Get historical fund prices
		_, err = client.GetHistoricalFundPrices(ctx, "fund-symbol", RegionTr, HistoricalFundPricePeriodOneYear)
		if err != nil {
			t.Errorf("GetHistoricalFundPrices failed: %v", err)
		}
	})

	// Test Financial Data Client methods
	t.Run("Financial Data Client", func(t *testing.T) {
		// Get financial ratios
		_, err := client.GetHistoricalRatios(ctx, "THYAO", []HistoricalRatiosKey{HistoricalRatiosKeyPERatio}, RegionTr)
		if err != nil {
			t.Errorf("GetHistoricalRatios failed: %v", err)
		}

		// Get financial ratio comparisons
		_, err = client.GetFinancialRatioComparison(ctx, "TUPRS", RegionTr, PeerTypeSector)
		if err != nil {
			t.Errorf("GetFinancialRatioComparison failed: %v", err)
		}

		// Get financial statements
		_, err = client.GetHistoricalFinancialSheets(ctx, "THYAO", FinancialSheetDate{Year: 2024, Month: 1, Day: 1}, FinancialSheetDate{Year: 2024, Month: 12, Day: 31}, FinancialSheetIncomeStatement, FinancialSheetPeriodAnnual, CurrencyUSD, RegionUs)
		if err != nil {
			t.Errorf("GetHistoricalFinancialSheets failed: %v", err)
		}

		// Get stock dividends
		_, err = client.GetStockDividends(ctx, "AAPL", RegionUs)
		if err != nil {
			t.Errorf("GetStockDividends failed: %v", err)
		}

		// Get stock statistics
		_, err = client.GetStockStats(ctx, []string{"AAPL", "GOOGL"}, RegionUs)
		if err != nil {
			t.Errorf("GetStockStats failed: %v", err)
		}

		// Get top movers
		_, err = client.GetTopMovers(ctx, TopMoversDirectionGainers, AssetClassEquity, AssetTypeStock, 1, 10, RegionTr)
		if err != nil {
			t.Errorf("GetTopMovers failed: %v", err)
		}
	})

	// Test Live Price Client methods
	t.Run("Live Price Client", func(t *testing.T) {
		// Get live prices for BIST stocks
		_, err := client.GetLivePriceForBIST(ctx, []string{"THYAO", "GARAN"})
		if err != nil {
			t.Errorf("GetLivePriceForBIST failed: %v", err)
		}

		// Get live prices for US stocks
		_, err = client.GetLivePriceForUS(ctx, []string{"AAPL", "GOOGL"})
		if err != nil {
			t.Errorf("GetLivePriceForUS failed: %v", err)
		}
	})

	// Test Brokers Client methods
	t.Run("Brokers Client", func(t *testing.T) {
		// Get all brokers
		_, err := client.GetBrokers(ctx, RegionTr, 1, 10)
		if err != nil {
			t.Errorf("GetBrokers failed: %v", err)
		}

		// Get market stocks with broker statistics
		_, err = client.GetMarketStocks(ctx, RegionTr, BrokerSortNetAmount, BrokerSortDirectionDesc, "2024-01-01", "2024-01-31", 1, 10)
		if err != nil {
			t.Errorf("GetMarketStocks failed: %v", err)
		}

		// Get brokers by stock
		_, err = client.GetBrokersByStock(ctx, "THYAO", RegionTr, BrokerSortNetAmount, BrokerSortDirectionDesc, "2024-01-01", "2024-01-31", 1, 10)
		if err != nil {
			t.Errorf("GetBrokersByStock failed: %v", err)
		}
	})

	// Test Search Client methods
	t.Run("Search Client", func(t *testing.T) {
		// Search across stocks, collections, sectors, and industries
		_, err := client.Search(ctx, "technology", []SearchType{SearchTypeStock, SearchTypeCollection}, RegionUs, LocaleEn, 1, PageSize20)
		if err != nil {
			t.Errorf("Search failed: %v", err)
		}
	})

	// Test WebSocket Client methods
	t.Run("WebSocket Client", func(t *testing.T) {
		// Get WebSocket URL for real-time data
		_, err := client.GetWebSocketUrl(ctx, "user-id", []FeedType{FeedTypeLivePriceTR}, RegionTr)
		if err != nil {
			t.Errorf("GetWebSocketUrl failed: %v", err)
		}

		// Update user details
		err = client.UpdateUserDetails(ctx, UpdateUserDetailsParams{
			ExternalUserID: "user-id",
			FirstName:      "John",
			LastName:       "Doe",
			Address:        "123 Main St",
			City:           "New York",
			CountryCode:    "US",
			AccessorType:   AccessorTypeUser,
			Active:         true,
		})
		if err != nil {
			t.Errorf("UpdateUserDetails failed: %v", err)
		}
	})

	// Test Capital Increase Client methods
	t.Run("Capital Increase Client", func(t *testing.T) {
		// Get all capital increases
		_, err := client.GetAllCapitalIncreases(ctx, 1, 10, RegionTr)
		if err != nil {
			t.Errorf("GetAllCapitalIncreases failed: %v", err)
		}

		// Get capital increases for a specific instrument
		_, err = client.GetCapitalIncreasesForInstrument(ctx, "THYAO", 1, 10, RegionTr)
		if err != nil {
			t.Errorf("GetCapitalIncreasesForInstrument failed: %v", err)
		}

		// Get active rights for an instrument
		_, err = client.GetActiveRightsForInstrument(ctx, "THYAO", "2024-01-15", RegionTr)
		if err != nil {
			t.Errorf("GetActiveRightsForInstrument failed: %v", err)
		}
	})

	// Test Custom Themes Client methods
	t.Run("Custom Themes Client", func(t *testing.T) {
		// Get all custom themes
		_, err := client.GetAllCustomThemes(ctx, RegionTr, LocaleEn)
		if err != nil {
			t.Errorf("GetAllCustomThemes failed: %v", err)
		}

		// Get custom theme detail
		_, err = client.GetCustomThemeDetail(ctx, "620f455a0187ade00bb0d55f", LocaleEn, SortByPriceChange)
		if err != nil {
			t.Errorf("GetCustomThemeDetail failed: %v", err)
		}

		// Create a custom theme
		stockID, _ := primitive.ObjectIDFromHex("648ab66e38daf3102a5a7401")
		_, err = client.CreateCustomTheme(ctx, CreateCustomThemeParams{
			Title:       LocaleString{"en": "My Tech Portfolio"},
			Description: LocaleString{"en": "Technology stocks portfolio"},
			Region:      []Region{RegionUs},
			Stocks:      []primitive.ObjectID{stockID},
			Status:      CollectionStatusActive,
		})
		if err != nil {
			t.Errorf("CreateCustomTheme failed: %v", err)
		}

		// Update a custom theme (using a valid ID)
		themeID, _ := primitive.ObjectIDFromHex("620f455a0187ade00bb0d55f")
		err = client.UpdateCustomTheme(ctx, themeID, UpdateCustomThemeParams{
			Title:  LocaleString{"en": "Updated Tech Portfolio"},
			Stocks: []primitive.ObjectID{stockID},
		})
		if err != nil {
			t.Errorf("UpdateCustomTheme failed: %v", err)
		}

		// Delete a custom theme
		err = client.DeleteCustomTheme(ctx, themeID)
		if err != nil {
			t.Errorf("DeleteCustomTheme failed: %v", err)
		}
	})

	// Test Key Insights Client methods
	t.Run("Key Insights Client", func(t *testing.T) {
		// Get key insights for a stock
		_, err := client.GetKeyInsights(ctx, "AAPL", RegionUs)
		if err != nil {
			t.Errorf("GetKeyInsights failed: %v", err)
		}
	})
}
