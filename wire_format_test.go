package laplace

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// The tests in this file replay the exact payloads the Laplace gateway emits
// against an httptest server, so they need neither an API key nor network
// access. They guard the JSON field names of the streaming and REST models.

func newMockClient(t *testing.T, handler http.HandlerFunc) *Client {
	t.Helper()
	srv := httptest.NewServer(handler)
	t.Cleanup(srv.Close)

	c, err := NewClient(LaplaceConfiguration{APIKey: "test-key", BaseURL: srv.URL})
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	return c
}

// sseHandler serves events on path as an SSE stream and then ends the response.
// wantQuery lists query parameters that must be present with the given value.
func sseHandler(t *testing.T, path string, wantQuery map[string]string, events ...string) http.HandlerFunc {
	t.Helper()
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != path {
			t.Errorf("request path = %s, want %s", r.URL.Path, path)
			http.NotFound(w, r)
			return
		}
		if got := r.Header.Get("Authorization"); got != "Bearer test-key" {
			t.Errorf("Authorization header = %q", got)
		}
		for k, v := range wantQuery {
			if got := r.URL.Query().Get(k); got != v {
				t.Errorf("query %s = %q, want %q", k, got, v)
			}
		}
		w.Header().Set("Content-Type", "text/event-stream")
		w.WriteHeader(http.StatusOK)
		for _, ev := range events {
			_, _ = w.Write([]byte("data: " + ev + "\n\n"))
		}
		if f, ok := w.(http.Flusher); ok {
			f.Flush()
		}
	}
}

// receive reads n payloads from ch or fails the test.
func receive[T any](t *testing.T, ch <-chan LivePriceResult[T], n int) []T {
	t.Helper()
	timeout := time.After(5 * time.Second)
	out := make([]T, 0, n)
	for len(out) < n {
		select {
		case r, ok := <-ch:
			if !ok {
				t.Fatalf("stream closed after %d of %d messages", len(out), n)
			}
			if r.Error != nil {
				t.Fatalf("stream error: %v", r.Error)
			}
			out = append(out, r.Data)
		case <-timeout:
			t.Fatalf("timed out after %d of %d messages", len(out), n)
		}
	}
	return out
}

func TestBISTPriceStreamDecodesV2Envelope(t *testing.T) {
	c := newMockClient(t, sseHandler(t, "/api/v2/stock/price/live",
		map[string]string{"filter": "THYAO", "region": "tr"},
		`{"t":"heartbeat"}`,
		`{"t":"pr","d":{"s":"THYAO","ch":-1.17,"p":296.75,"d":1789394100282}}`,
	))

	stream, err := c.CreateLivePriceStreamForBIST(context.Background(), []string{"THYAO"})
	if err != nil {
		t.Fatalf("CreateLivePriceStreamForBIST: %v", err)
	}
	defer stream.Close()

	msgs := receive(t, stream.Receive(), 2)

	if msgs[0].Type != MessageTypeHeartbeat {
		t.Errorf("first message type = %q, want heartbeat", msgs[0].Type)
	}
	if msgs[0].Data.Symbol != "" {
		t.Errorf("heartbeat carried data: %+v", msgs[0].Data)
	}

	tick := msgs[1]
	if tick.Type != MessageTypePrice {
		t.Errorf("second message type = %q, want pr", tick.Type)
	}
	want := BISTStockLiveData{Symbol: "THYAO", DailyPercentChange: -1.17, ClosePrice: 296.75, Date: 1789394100282}
	if tick.Data != want {
		t.Errorf("decoded tick = %+v, want %+v", tick.Data, want)
	}
}

func TestDelayedPriceStreamDecodesV2Envelope(t *testing.T) {
	c := newMockClient(t, sseHandler(t, "/api/v1/stock/price/delayed",
		map[string]string{"filter": "GARAN", "region": "tr"},
		`{"t":"pr","d":{"s":"GARAN","ch":0.42,"p":118.5,"d":1789394100282}}`,
	))

	stream, err := c.CreateDelayedPriceStreamForBIST(context.Background(), []string{"GARAN"})
	if err != nil {
		t.Fatalf("CreateDelayedPriceStreamForBIST: %v", err)
	}
	defer stream.Close()

	tick := receive(t, stream.Receive(), 1)[0]
	want := BISTStockLiveData{Symbol: "GARAN", DailyPercentChange: 0.42, ClosePrice: 118.5, Date: 1789394100282}
	if tick.Type != MessageTypePrice || tick.Data != want {
		t.Errorf("decoded message = %+v, want type pr with %+v", tick, want)
	}
}

func TestUSPriceStreamDecodesV2Envelope(t *testing.T) {
	c := newMockClient(t, sseHandler(t, "/api/v2/stock/price/live",
		map[string]string{"filter": "AAPL", "region": "us"},
		`{"t":"pr","d":{"s":"AAPL","p":230.5,"pc":1.2,"ac":2.7,"d":1789394100282}}`,
	))

	stream, err := c.CreateLivePriceStreamForUS(context.Background(), []string{"AAPL"})
	if err != nil {
		t.Fatalf("CreateLivePriceStreamForUS: %v", err)
	}
	defer stream.Close()

	tick := receive(t, stream.Receive(), 1)[0]
	want := USStockLiveData{Symbol: "AAPL", Price: 230.5, PercentChange: 1.2, AmountChange: 2.7, Date: 1789394100282}
	if tick.Type != MessageTypePrice || tick.Data != want {
		t.Errorf("decoded message = %+v, want type pr with %+v", tick, want)
	}
}

func TestOrderBookStreamDecodesSymbol(t *testing.T) {
	c := newMockClient(t, sseHandler(t, "/api/v1/stock/orderbook/live",
		map[string]string{"filter": "THYAO", "region": "tr"},
		`{"updated":[{"level":1,"side":"bid","vol":1500,"orders":3,"p":296.5}],"deleted":[{"level":10,"side":"ask"}],"symbol":"THYAO"}`,
	))

	stream, err := c.CreateLiveOrderBookStreamForBIST(context.Background(), []string{"THYAO"})
	if err != nil {
		t.Fatalf("CreateLiveOrderBookStreamForBIST: %v", err)
	}
	defer stream.Close()

	book := receive(t, stream.Receive(), 1)[0]
	if book.Symbol != "THYAO" {
		t.Errorf("Symbol = %q, want THYAO", book.Symbol)
	}
	if len(book.Updated) != 1 || book.Updated[0] != (OrderbookLevel{ID: 1, Side: LevelSideBid, Volume: 1500, Orders: 3, Price: 296.5}) {
		t.Errorf("Updated = %+v", book.Updated)
	}
	if len(book.Deleted) != 1 || book.Deleted[0] != (OrderbookDeletedLevel{ID: 10, Side: LevelSideAsk}) {
		t.Errorf("Deleted = %+v", book.Deleted)
	}
}

func TestBidAskStreamDecodesV2Envelope(t *testing.T) {
	c := newMockClient(t, sseHandler(t, "/api/v1/stock/price/bids",
		map[string]string{"filter": "THYAO", "region": "tr"},
		`{"t":"pr","d":{"s":"THYAO","bid":296.5,"ask":296.75,"d":1789394100282}}`,
	))

	stream, err := c.CreateLiveBidAskStreamForBIST(context.Background(), []string{"THYAO"})
	if err != nil {
		t.Fatalf("CreateLiveBidAskStreamForBIST: %v", err)
	}
	defer stream.Close()

	msg := receive(t, stream.Receive(), 1)[0]
	want := BISTBidAskLiveData{Symbol: "THYAO", Bid: 296.5, Ask: 296.75, Date: 1789394100282}
	if msg.Type != MessageTypePrice || msg.Data != want {
		t.Errorf("decoded message = %+v, want type pr with %+v", msg, want)
	}
}

func TestGetSectorDetailReturnsUpstreamError(t *testing.T) {
	c := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"message":"sector not found","error_code":"not_found"}`))
	})

	_, err := c.GetSectorDetail(context.Background(), "000000000000000000000000", RegionTr, LocaleTr)
	if err == nil {
		t.Fatal("GetSectorDetail returned nil error for a 404 response")
	}
	var httpErr *LaplaceHTTPError
	if !errors.As(err, &httpErr) {
		t.Fatalf("error type = %T, want *LaplaceHTTPError", err)
	}
	if httpErr.HTTPStatus != http.StatusNotFound || httpErr.Message.Message != "sector not found" {
		t.Errorf("error = %+v", httpErr)
	}
}
