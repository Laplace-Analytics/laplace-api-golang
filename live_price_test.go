package laplace

import (
	"context"
	"slices"
	"testing"
	"time"
)

func TestGetLivePriceForBIST(t *testing.T) {
	cfg, err := loadConfig()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	client, err := NewClient(*cfg)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Use new unified streaming API
	stream, err := client.CreateLivePriceStreamForBIST(ctx, []string{"AKBNK"})
	if err != nil {
		t.Fatalf("Failed to create live price stream: %v", err)
	}
	defer stream.Close()

	receiveChan := stream.Receive()

	select {
	case data := <-receiveChan:
		if data.Error != nil {
			t.Logf("Received error: %v", data.Error)
		} else {
			t.Logf("Received data: %+v", data.Data)
		}
	case <-ctx.Done():
		t.Log("Timeout waiting for data")
	}
}

func TestGetLivePriceForUS(t *testing.T) {
	cfg, err := loadConfig()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	client, err := NewClient(*cfg)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Use new unified streaming API
	stream, err := client.CreateLivePriceStreamForUS(ctx, []string{"AAPL"})
	if err != nil {
		t.Fatalf("Failed to create live price stream: %v", err)
	}
	defer stream.Close()

	receiveChan := stream.Receive()

	select {
	case data := <-receiveChan:
		if data.Error != nil {
			t.Logf("Received error: %v", data.Error)
		} else {
			t.Logf("Received data: %+v", data.Data)
		}
	case <-ctx.Done():
		t.Log("Timeout waiting for data")
	}
}

// collectPriceSymbols drains ch for up to d and returns the symbols of the
// price ticks it saw, skipping heartbeats and state changes. It returns early
// if ch is closed and fails the test on a stream error.
func collectPriceSymbols[T any](t *testing.T, ch <-chan LivePriceResult[LiveMessageV2[T]], symbolOf func(T) string, d time.Duration) []string {
	t.Helper()

	deadline := time.After(d)
	symbols := []string{}
	for {
		select {
		case msg, ok := <-ch:
			if !ok {
				return symbols
			}
			if msg.Error != nil {
				t.Fatalf("Received error: %v", msg.Error)
			}
			if msg.Data.Type != MessageTypePrice {
				continue
			}
			symbols = append(symbols, symbolOf(msg.Data.Data))
		case <-deadline:
			return symbols
		}
	}
}

// assertSwitchedSymbols checks that ticks seen before the switch belong to the
// first subscription (AKBNK) and ticks seen after it belong to the second
// (TUPRS, ASELS). Either window may be empty outside market hours; that is
// logged rather than failed.
func assertSwitchedSymbols(t *testing.T, beforeSwitch, afterSwitch []string) {
	t.Helper()

	if len(beforeSwitch) == 0 {
		t.Log("No price ticks received before switch (market closed?)")
	} else {
		if !slices.Contains(beforeSwitch, "AKBNK") {
			t.Errorf("Did not receive AKBNK data before switch, got %v", beforeSwitch)
		}
		if slices.Contains(beforeSwitch, "TUPRS") || slices.Contains(beforeSwitch, "ASELS") {
			t.Errorf("Received data for the second subscription before switch: %v", beforeSwitch)
		}
	}

	if len(afterSwitch) == 0 {
		t.Log("No price ticks received after switch (market closed?)")
	} else {
		if !slices.Contains(afterSwitch, "TUPRS") {
			t.Errorf("Did not receive TUPRS data after switch, got %v", afterSwitch)
		}
		if !slices.Contains(afterSwitch, "ASELS") {
			t.Errorf("Did not receive ASELS data after switch, got %v", afterSwitch)
		}
		if slices.Contains(afterSwitch, "AKBNK") {
			t.Errorf("Still receiving AKBNK data after switch: %v", afterSwitch)
		}
	}
}

func TestLivePriceSubscribe(t *testing.T) {
	cfg, err := loadConfig()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	client, err := NewClient(*cfg)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Use new manual stream creation for more control
	stream := client.GetLivePriceStreamForBIST()
	if err := stream.Subscribe(ctx, []string{"AKBNK"}); err != nil {
		t.Fatalf("Failed to subscribe to live price stream: %v", err)
	}
	defer stream.Close()

	symbolOf := func(d BISTStockLiveData) string { return d.Symbol }

	// Collect from the first subscription, then switch symbols from this
	// goroutine. Subscribe closes the old channel, so Receive is called again
	// to get the new one.
	beforeSwitch := collectPriceSymbols(t, stream.Receive(), symbolOf, 5*time.Second)

	if err := stream.Subscribe(ctx, []string{"TUPRS", "ASELS"}); err != nil {
		t.Fatalf("Failed to switch symbols: %v", err)
	}
	afterSwitch := collectPriceSymbols(t, stream.Receive(), symbolOf, 5*time.Second)

	assertSwitchedSymbols(t, beforeSwitch, afterSwitch)
}

func TestLivePriceClose(t *testing.T) {
	cfg, err := loadConfig()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	client, err := NewClient(*cfg)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	// Use new unified streaming API
	stream, err := client.CreateLivePriceStreamForBIST(ctx, []string{"AKBNK"})
	if err != nil {
		t.Fatalf("Failed to create live price stream: %v", err)
	}

	err = stream.Close()
	if err != nil {
		t.Fatalf("Close failed: %v", err)
	}
}

// Test new unified streaming API for order book
func TestOrderBookStream(t *testing.T) {
	cfg, err := loadConfig()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	client, err := NewClient(*cfg)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Test order book streaming
	stream, err := client.CreateLiveOrderBookStreamForBIST(ctx, []string{"THYAO"})
	if err != nil {
		t.Fatalf("Failed to create order book stream: %v", err)
	}
	defer stream.Close()

	receiveChan := stream.Receive()

	select {
	case data := <-receiveChan:
		if data.Error != nil {
			t.Logf("Received error: %v", data.Error)
		} else {
			t.Logf("Received order book data: Symbol=%s, Updated=%d, Deleted=%d",
				data.Data.Symbol, len(data.Data.Updated), len(data.Data.Deleted))
		}
	case <-ctx.Done():
		t.Log("Timeout waiting for order book data")
	}
}

// Test new unified streaming API for delayed price
func TestDelayedPriceStream(t *testing.T) {
	cfg, err := loadConfig()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	client, err := NewClient(*cfg)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Test delayed price streaming
	stream, err := client.CreateDelayedPriceStreamForBIST(ctx, []string{"THYAO"})
	if err != nil {
		t.Fatalf("Failed to create delayed price stream: %v", err)
	}
	defer stream.Close()

	receiveChan := stream.Receive()

	select {
	case data := <-receiveChan:
		if data.Error != nil {
			t.Logf("Received error: %v", data.Error)
		} else {
			t.Logf("Received delayed price data: %+v", data.Data)
		}
	case <-ctx.Done():
		t.Log("Timeout waiting for delayed price data")
	}
}

func TestGetLiveBidAskForBIST(t *testing.T) {
	cfg, err := loadConfig()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	client, err := NewClient(*cfg)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	stream, err := client.CreateLiveBidAskStreamForBIST(ctx, []string{"AKBNK", "ISCTR"})
	if err != nil {
		t.Fatalf("Failed to create live bid/ask stream: %v", err)
	}
	defer stream.Close()

	receiveChan := stream.Receive()

	select {
	case data := <-receiveChan:
		if data.Error != nil {
			t.Logf("Received error: %v", data.Error)
		} else {
			t.Logf("Received bid/ask data: %+v", data.Data)

			// Verify the structure of the response
			if data.Data.Data.Symbol == "" {
				t.Error("Symbol should not be empty")
			}
			if data.Data.Data.Bid <= 0 {
				t.Error("Bid price should be greater than 0")
			}
			if data.Data.Data.Ask <= 0 {
				t.Error("Ask price should be greater than 0")
			}
			if data.Data.Data.Ask <= data.Data.Data.Bid {
				t.Error("Ask price should be greater than bid price")
			}
			if data.Data.Data.Date <= 0 {
				t.Error("Date should be a valid timestamp")
			}
			if data.Data.Type == "" {
				t.Error("Type should not be empty")
			}
		}
	case <-ctx.Done():
		t.Log("Timeout waiting for bid/ask data")
	}
}

func TestGetLiveBidAskForBIST_AllSymbols(t *testing.T) {
	cfg, err := loadConfig()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	client, err := NewClient(*cfg)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Test with empty symbols array (should stream all BIST stocks)
	stream, err := client.CreateLiveBidAskStreamForBIST(ctx, []string{})
	if err != nil {
		t.Fatalf("Failed to create live bid/ask stream for all symbols: %v", err)
	}
	defer stream.Close()

	receiveChan := stream.Receive()

	select {
	case data := <-receiveChan:
		if data.Error != nil {
			t.Logf("Received error: %v", data.Error)
		} else {
			t.Logf("Received bid/ask data for all symbols: %+v", data.Data)

			// Verify we got some data
			if data.Data.Data.Symbol == "" {
				t.Error("Should receive data for at least one symbol")
			}
		}
	case <-ctx.Done():
		t.Log("Timeout waiting for bid/ask data (all symbols)")
	}
}

func TestLiveBidAskSubscribe(t *testing.T) {
	cfg, err := loadConfig()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	client, err := NewClient(*cfg)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	stream := client.GetLiveBidAskStreamForBIST()
	if err := stream.Subscribe(ctx, []string{"AKBNK"}); err != nil {
		t.Fatalf("Failed to subscribe to bid/ask stream: %v", err)
	}
	defer stream.Close()

	symbolOf := func(d BISTBidAskLiveData) string { return d.Symbol }

	beforeSwitch := collectPriceSymbols(t, stream.Receive(), symbolOf, 5*time.Second)

	// Switch to different symbols
	if err := stream.Subscribe(ctx, []string{"TUPRS", "ASELS"}); err != nil {
		t.Fatalf("Failed to switch symbols: %v", err)
	}
	afterSwitch := collectPriceSymbols(t, stream.Receive(), symbolOf, 5*time.Second)

	assertSwitchedSymbols(t, beforeSwitch, afterSwitch)
}

func TestLiveBidAskClose(t *testing.T) {
	cfg, err := loadConfig()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	client, err := NewClient(*cfg)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	stream, err := client.CreateLiveBidAskStreamForBIST(ctx, []string{"AKBNK"})
	if err != nil {
		t.Fatalf("Failed to create live bid/ask stream: %v", err)
	}

	err = stream.Close()
	if err != nil {
		t.Fatalf("Close failed: %v", err)
	}

	// Test multiple closes (should not error)
	err = stream.Close()
	if err != nil {
		t.Fatalf("Second close failed: %v", err)
	}
}

func TestLiveBidAsk_NilContext(t *testing.T) {
	cfg, err := loadConfig()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	client, err := NewClient(*cfg)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	stream := client.GetLiveBidAskStreamForBIST()
	err = stream.Subscribe(nil, []string{"AKBNK"}) //nolint:staticcheck // intentionally testing nil context rejection
	if err == nil {
		t.Fatal("Expected error for nil context")
	}

	expectedError := "context cannot be nil"
	if err.Error() != expectedError {
		t.Errorf("Expected error message '%s', got '%s'", expectedError, err.Error())
	}
}

func TestLiveBidAskDataStructure(t *testing.T) {
	// Test BISTBidAskResponse and BISTBidAskLiveData structures
	response := BISTBidAskResponse{
		Data: BISTBidAskLiveData{
			Symbol: "AKBNK",
			Ask:    45.6,
			Bid:    45.5,
			Date:   1740414373252,
		},
		Type: "pr",
	}

	if response.Data.Symbol != "AKBNK" {
		t.Errorf("Expected symbol AKBNK, got %s", response.Data.Symbol)
	}
	if response.Data.Ask != 45.6 {
		t.Errorf("Expected ask 45.6, got %f", response.Data.Ask)
	}
	if response.Data.Bid != 45.5 {
		t.Errorf("Expected bid 45.5, got %f", response.Data.Bid)
	}
	if response.Data.Date != 1740414373252 {
		t.Errorf("Expected date 1740414373252, got %d", response.Data.Date)
	}
	if response.Type != "pr" {
		t.Errorf("Expected type pr, got %s", response.Type)
	}
}
