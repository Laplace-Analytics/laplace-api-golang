package laplace

import (
	"context"
	"slices"
	"testing"
	"time"
)

func TestGetDelayedPriceForBIST(t *testing.T) {
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

	dc, err := client.GetDelayedPriceForBIST(ctx, []string{"AKBNK"})
	if err != nil {
		t.Fatalf("Failed to create delayed price client: %v", err)
	}
	defer dc.Close()

	receiveChan := dc.Receive()

	select {
	case data := <-receiveChan:
		if data.Error != nil {
			t.Logf("Received error: %v", data.Error)
		} else {
			t.Logf("Received delayed data: %+v", data.Data)
			// Verify it's the expected symbol
			if data.Data.Symbol != "AKBNK" {
				t.Errorf("Expected symbol AKBNK, got %s", data.Data.Symbol)
			}
			// Verify basic data structure
			if data.Data.ClosePrice <= 0 {
				t.Errorf("Expected positive close price, got %f", data.Data.ClosePrice)
			}
			if data.Data.Date <= 0 {
				t.Errorf("Expected positive date, got %d", data.Data.Date)
			}
		}
	case <-ctx.Done():
		t.Log("Timeout waiting for delayed data")
	}
}

func TestGetDelayedPriceForUS(t *testing.T) {
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

	dc, err := client.GetDelayedPriceForUS(ctx, []string{"AAPL"})
	if err != nil {
		t.Fatalf("Failed to create delayed price client: %v", err)
	}
	defer dc.Close()

	receiveChan := dc.Receive()

	select {
	case data := <-receiveChan:
		if data.Error != nil {
			t.Logf("Received error: %v", data.Error)
		} else {
			t.Logf("Received delayed data: %+v", data.Data)
			// Verify it's the expected symbol
			if data.Data.Symbol != "AAPL" {
				t.Errorf("Expected symbol AAPL, got %s", data.Data.Symbol)
			}
			// Verify basic data structure
			if data.Data.Price <= 0 {
				t.Errorf("Expected positive price, got %f", data.Data.Price)
			}
			if data.Data.Date <= 0 {
				t.Errorf("Expected positive date, got %d", data.Data.Date)
			}
		}
	case <-ctx.Done():
		t.Log("Timeout waiting for delayed data")
	}
}

func TestDelayedPriceSubscribe(t *testing.T) {
	cfg, err := loadConfig()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	client, err := NewClient(*cfg)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	dc, err := client.GetDelayedPriceForBIST(ctx, []string{"AKBNK"})
	if err != nil {
		t.Fatalf("Failed to create delayed price client: %v", err)
	}
	defer dc.Close()

	receivedData := []string{}

	timer := time.NewTimer(5 * time.Second)
	go func() {
		<-timer.C
		// Subscribe to different symbols
		err := dc.Subscribe(ctx, []string{"TUPRS", "ASELS"})
		if err != nil {
			t.Errorf("Subscribe failed: %v", err)
		}
		receivedData = append(receivedData, "SWITCH")

		timer.Reset(5 * time.Second)
		<-timer.C
		dc.Close()
	}()

	receiveChan := dc.Receive()
	for data := range receiveChan {
		if data.Error != nil {
			t.Fatalf("Received error: %v", data.Error)
		}

		receivedData = append(receivedData, data.Data.Symbol)
	}

	idxOfSwitch := slices.Index(receivedData, "SWITCH")

	// Check data received before switch
	if idxOfSwitch > 0 {
		beforeSwitch := receivedData[:idxOfSwitch]
		if !slices.Contains(beforeSwitch, "AKBNK") {
			t.Error("Did not receive AKBNK delayed data before switch")
		}
	}

	// Check data received after switch
	if idxOfSwitch >= 0 && idxOfSwitch < len(receivedData)-1 {
		afterSwitch := receivedData[idxOfSwitch+1:]
		if !slices.Contains(afterSwitch, "TUPRS") {
			t.Error("Did not receive TUPRS delayed data after switch")
		}
		if !slices.Contains(afterSwitch, "ASELS") {
			t.Error("Did not receive ASELS delayed data after switch")
		}
	}
}

func TestDelayedPriceClose(t *testing.T) {
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

	dc, err := client.GetDelayedPriceForBIST(ctx, []string{"AKBNK"})
	if err != nil {
		t.Fatalf("Failed to create delayed price client: %v", err)
	}

	err = dc.Close()
	if err != nil {
		t.Fatalf("Close failed: %v", err)
	}

	// Test double close - should not error
	err = dc.Close()
	if err != nil {
		t.Fatalf("Double close failed: %v", err)
	}
}

func TestDelayedPriceMultipleSymbols(t *testing.T) {
	cfg, err := loadConfig()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	client, err := NewClient(*cfg)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// Test with multiple symbols
	symbols := []string{"AKBNK", "TUPRS", "ASELS"}
	dc, err := client.GetDelayedPriceForBIST(ctx, symbols)
	if err != nil {
		t.Fatalf("Failed to create delayed price client: %v", err)
	}
	defer dc.Close()

	receiveChan := dc.Receive()
	receivedSymbols := make(map[string]bool)

	timeout := time.After(10 * time.Second)

	for len(receivedSymbols) < len(symbols) {
		select {
		case data := <-receiveChan:
			if data.Error != nil {
				t.Logf("Received error: %v", data.Error)
				continue
			}
			receivedSymbols[data.Data.Symbol] = true
			t.Logf("Received delayed data for %s: price=%f, change=%f%%",
				data.Data.Symbol, data.Data.ClosePrice, data.Data.DailyPercentChange)
		case <-timeout:
			t.Log("Timeout waiting for all symbols")
			break
		}
	}

	// Verify we received data for all requested symbols
	for _, symbol := range symbols {
		if !receivedSymbols[symbol] {
			t.Errorf("Did not receive delayed data for symbol: %s", symbol)
		}
	}
}

func TestDelayedPriceErrorHandling(t *testing.T) {
	cfg, err := loadConfig()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	client, err := NewClient(*cfg)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	// Test with nil context
	_, err = client.GetDelayedPriceForBIST(nil, []string{"AKBNK"})
	if err == nil {
		t.Error("Expected error with nil context")
	}

	// Test with invalid symbol
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	dc, err := client.GetDelayedPriceForBIST(ctx, []string{"INVALID_SYMBOL"})
	if err != nil {
		t.Logf("Expected error with invalid symbol: %v", err)
		return
	}
	defer dc.Close()

	// If no error during creation, check if we get error in stream
	receiveChan := dc.Receive()
	select {
	case data := <-receiveChan:
		if data.Error != nil {
			t.Logf("Received expected error for invalid symbol: %v", data.Error)
		} else {
			t.Logf("Received data for supposedly invalid symbol: %+v", data.Data)
		}
	case <-ctx.Done():
		t.Log("Timeout - no error or data received for invalid symbol")
	}
}
