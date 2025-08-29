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

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	lc, err := client.GetLivePriceForBIST(ctx, []string{"AKBNK"})
	if err != nil {
		t.Fatalf("Failed to create live price client: %v", err)
	}
	defer lc.Close()

	receiveChan := lc.Receive()

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

	lc, err := client.GetLivePriceForUS(ctx, []string{"AAPL"})
	if err != nil {
		t.Fatalf("Failed to create live price client: %v", err)
	}
	defer lc.Close()

	receiveChan := lc.Receive()

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

func TestLivePriceSubscribe(t *testing.T) {
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

	lc, err := client.GetLivePriceForBIST(ctx, []string{"AKBNK"})
	if err != nil {
		t.Fatalf("Failed to create live price client: %v", err)
	}
	defer lc.Close()

	receivedData := []string{}

	timer := time.NewTimer(5 * time.Second)
	go func() {
		<-timer.C
		lc.Subscribe(ctx, []string{"TUPRS", "ASELS"})
		receivedData = append(receivedData, "SWITCH")

		timer.Reset(5 * time.Second)
		<-timer.C
		lc.Close()
	}()

	receiveChan := lc.Receive()
	for data := range receiveChan {
		if data.Error != nil {
			t.Fatalf("Received error: %v", data.Error)
		}

		receivedData = append(receivedData, data.Data.Symbol)
	}

	idxOfSwitch := slices.Index(receivedData, "SWITCH")

	if idxOfSwitch > 0 {
		beforeSwitch := receivedData[:idxOfSwitch]
		if !slices.Contains(beforeSwitch, "AKBNK") {
			t.Error("Did not receive AKBNK data before switch")
		}
	}

	if idxOfSwitch >= 0 && idxOfSwitch < len(receivedData)-1 {
		afterSwitch := receivedData[idxOfSwitch+1:]
		if !slices.Contains(afterSwitch, "TUPRS") {
			t.Error("Did not receive TUPRS data after switch")
		}
		if !slices.Contains(afterSwitch, "ASELS") {
			t.Error("Did not receive ASELS data after switch")
		}
	}
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

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	lc, err := client.GetLivePriceForBIST(ctx, []string{"AKBNK"})
	if err != nil {
		t.Fatalf("Failed to create live price client: %v", err)
	}

	err = lc.Close()
	if err != nil {
		t.Fatalf("Close failed: %v", err)
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

	lc, err := client.GetLiveBidAskForBIST(ctx, []string{"AKBNK", "ISCTR"})
	if err != nil {
		t.Fatalf("Failed to create live bid/ask client: %v", err)
	}
	defer lc.Close()

	receiveChan := lc.Receive()

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
	lc, err := client.GetLiveBidAskForBIST(ctx, []string{})
	if err != nil {
		t.Fatalf("Failed to create live bid/ask client for all symbols: %v", err)
	}
	defer lc.Close()

	receiveChan := lc.Receive()

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

	lc, err := client.GetLiveBidAskForBIST(ctx, []string{"AKBNK"})
	if err != nil {
		t.Fatalf("Failed to create live bid/ask client: %v", err)
	}
	defer lc.Close()

	receivedData := []string{}

	timer := time.NewTimer(5 * time.Second)
	go func() {
		<-timer.C
		// Switch to different symbols
		lc.Subscribe(ctx, []string{"TUPRS", "ASELS"})
		receivedData = append(receivedData, "SWITCH")

		timer.Reset(5 * time.Second)
		<-timer.C
		lc.Close()
	}()

	receiveChan := lc.Receive()
	for data := range receiveChan {
		if data.Error != nil {
			t.Fatalf("Received error: %v", data.Error)
		}

		receivedData = append(receivedData, data.Data.Data.Symbol)
	}

	idxOfSwitch := slices.Index(receivedData, "SWITCH")

	if idxOfSwitch > 0 {
		beforeSwitch := receivedData[:idxOfSwitch]
		if !slices.Contains(beforeSwitch, "AKBNK") {
			t.Error("Did not receive AKBNK bid/ask data before switch")
		}
	}

	if idxOfSwitch >= 0 && idxOfSwitch < len(receivedData)-1 {
		afterSwitch := receivedData[idxOfSwitch+1:]
		if !slices.Contains(afterSwitch, "TUPRS") {
			t.Error("Did not receive TUPRS bid/ask data after switch")
		}
		if !slices.Contains(afterSwitch, "ASELS") {
			t.Error("Did not receive ASELS bid/ask data after switch")
		}
	}
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

	lc, err := client.GetLiveBidAskForBIST(ctx, []string{"AKBNK"})
	if err != nil {
		t.Fatalf("Failed to create live bid/ask client: %v", err)
	}

	err = lc.Close()
	if err != nil {
		t.Fatalf("Close failed: %v", err)
	}

	// Test multiple closes (should not error)
	err = lc.Close()
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

	_, err = client.GetLiveBidAskForBIST(nil, []string{"AKBNK"})
	if err == nil {
		t.Fatal("Expected error for nil context")
	}

	expectedError := "context cannot be nil"
	if err.Error() != expectedError {
		t.Errorf("Expected error message '%s', got '%s'", expectedError, err.Error())
	}
}

func TestLiveBidAsk_NilClient(t *testing.T) {
	ctx := context.Background()

	var client *Client
	_, err := client.GetLiveBidAskForBIST(ctx, []string{"AKBNK"})
	if err == nil {
		t.Fatal("Expected error for nil client")
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
