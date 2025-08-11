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
