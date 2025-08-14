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
	stream := client.GetLivePriceStreamForBIST([]string{})
	err = stream.Subscribe(ctx, []string{"AKBNK"})
	if err != nil {
		t.Fatalf("Failed to subscribe to live price stream: %v", err)
	}
	defer stream.Close()

	receivedData := []string{}

	timer := time.NewTimer(5 * time.Second)
	go func() {
		<-timer.C
		stream.Subscribe(ctx, []string{"TUPRS", "ASELS"})
		receivedData = append(receivedData, "SWITCH")

		timer.Reset(5 * time.Second)
		<-timer.C
		stream.Close()
	}()

	receiveChan := stream.Receive()
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
