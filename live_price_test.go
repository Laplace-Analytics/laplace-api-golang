package laplace

import (
	"context"
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

	err = lc.Subscribe(ctx, []string{"TUPRS", "ASELS"})
	if err != nil {
		t.Fatalf("Subscribe failed: %v", err)
	}

	time.Sleep(2 * time.Second)
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
