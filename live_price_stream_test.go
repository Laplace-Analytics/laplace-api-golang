package laplace

import (
	"context"
	"fmt"
	"net/http"
	"runtime"
	"strings"
	"testing"
	"time"
)

// The tests in this file cover the lifecycle of LivePriceStream against an
// httptest server, so they need neither an API key nor network access:
// switching symbols with a second Subscribe, the server ending the stream,
// Close, and cancellation of the caller's context.

// holdingSSEHandler serves ticks price ticks for the requested filter symbol
// and then keeps the connection open until the client goes away, like the
// real gateway does between ticks.
func holdingSSEHandler(t *testing.T, ticks int) http.HandlerFunc {
	t.Helper()
	return func(w http.ResponseWriter, r *http.Request) {
		symbol := r.URL.Query().Get("filter")
		w.Header().Set("Content-Type", "text/event-stream")
		w.WriteHeader(http.StatusOK)
		for i := 0; i < ticks; i++ {
			fmt.Fprintf(w, "data: {\"t\":\"pr\",\"d\":{\"s\":%q,\"p\":%d}}\n\n", symbol, i+1)
		}
		if f, ok := w.(http.Flusher); ok {
			f.Flush()
		}
		<-r.Context().Done()
	}
}

// assertClosed drains ch and fails the test if it is not closed within d.
func assertClosed[T any](t *testing.T, ch <-chan LivePriceResult[T], d time.Duration) {
	t.Helper()
	deadline := time.After(d)
	for {
		select {
		case _, ok := <-ch:
			if !ok {
				return
			}
		case <-deadline:
			t.Fatalf("channel not closed within %v", d)
			return
		}
	}
}

// streamGoroutines counts goroutines currently inside the live price stream
// plumbing: the SSE reader started by sendSSERequest and the forwarder.
func streamGoroutines() (int, string) {
	buf := make([]byte, 1<<20)
	n := runtime.Stack(buf, true)
	dump := string(buf[:n])
	return strings.Count(dump, "sendSSERequest[") + strings.Count(dump, "LivePriceStream[...]).forwardData"), dump
}

func TestLivePriceStreamResubscribeSwitchesSymbols(t *testing.T) {
	c := newMockClient(t, holdingSSEHandler(t, 1))

	stream, err := c.CreateLivePriceStreamForBIST(context.Background(), []string{"THYAO"})
	if err != nil {
		t.Fatalf("CreateLivePriceStreamForBIST: %v", err)
	}
	defer stream.Close()

	first := stream.Receive()
	if got := receive(t, first, 1)[0]; got.Data.Symbol != "THYAO" {
		t.Fatalf("first subscription tick = %+v, want THYAO", got)
	}

	if err := stream.Subscribe(context.Background(), []string{"GARAN"}); err != nil {
		t.Fatalf("second Subscribe: %v", err)
	}

	// A consumer ranging over the first channel must be released.
	assertClosed(t, first, 5*time.Second)

	second := stream.Receive()
	if got := receive(t, second, 1)[0]; got.Data.Symbol != "GARAN" {
		t.Fatalf("second subscription tick = %+v, want GARAN", got)
	}

	// Nothing from the torn-down first connection may leak into the new
	// channel, in particular not its "context canceled" read error.
	select {
	case msg, ok := <-second:
		if ok {
			t.Errorf("unexpected message on the new channel after switch: %+v", msg)
		} else {
			t.Error("new channel closed unexpectedly")
		}
	case <-time.After(200 * time.Millisecond):
	}

	if err := stream.Close(); err != nil {
		t.Fatalf("Close: %v", err)
	}
	assertClosed(t, second, 5*time.Second)
}

func TestLivePriceStreamEndsWhenServerDisconnects(t *testing.T) {
	c := newMockClient(t, sseHandler(t, "/api/v2/stock/price/live", nil,
		`{"t":"pr","d":{"s":"THYAO","p":296.75}}`,
	))

	stream, err := c.CreateLivePriceStreamForBIST(context.Background(), []string{"THYAO"})
	if err != nil {
		t.Fatalf("CreateLivePriceStreamForBIST: %v", err)
	}
	defer stream.Close()

	// The README pattern: range until the channel closes.
	var ticks []BISTStockLiveData
	done := make(chan struct{})
	go func() {
		defer close(done)
		for msg := range stream.Receive() {
			if msg.Error != nil {
				t.Errorf("stream error: %v", msg.Error)
				continue
			}
			ticks = append(ticks, msg.Data.Data)
		}
	}()

	select {
	case <-done:
	case <-time.After(5 * time.Second):
		t.Fatal("Receive channel was not closed after the server ended the stream")
	}
	if len(ticks) != 1 || ticks[0].Symbol != "THYAO" {
		t.Errorf("ticks = %+v, want a single THYAO tick", ticks)
	}
}

func TestLivePriceStreamContextCancelReleasesReader(t *testing.T) {
	// More events than the output buffer holds, so the reader goroutine is
	// parked on a send when the context is cancelled.
	c := newMockClient(t, holdingSSEHandler(t, 300))

	baseline, _ := streamGoroutines()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	stream, err := c.CreateLivePriceStreamForBIST(ctx, []string{"THYAO"})
	if err != nil {
		t.Fatalf("CreateLivePriceStreamForBIST: %v", err)
	}
	defer stream.Close()

	ch := stream.Receive()
	receive(t, ch, 1) // the connection is up and events are flowing

	cancel()

	assertClosed(t, ch, 5*time.Second)

	deadline := time.Now().Add(5 * time.Second)
	for {
		n, dump := streamGoroutines()
		if n <= baseline {
			return
		}
		if time.Now().After(deadline) {
			t.Fatalf("stream goroutines still running 5s after cancel (%d, baseline %d):\n%s", n, baseline, dump)
		}
		time.Sleep(10 * time.Millisecond)
	}
}
