package laplace

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/google/uuid"
)

// LivePriceType represents the type of live price stream
type LivePriceType string

const (
	LivePriceTypePrice        LivePriceType = "price"
	LivePriceTypeDelayedPrice LivePriceType = "delayed-price"
	LivePriceTypeOrderBook    LivePriceType = "order-book"
	LivePriceTypeBidAsk       LivePriceType = "bid-ask"
)

// MessageType represents the type of message in live data streams
type MessageType string

const (
	MessageTypePrice       MessageType = "pr"
	MessageTypeStateChange MessageType = "state_change"
	MessageTypeHeartbeat   MessageType = "heartbeat"
	MessageTypeOrderbook   MessageType = "ob"
)

// LiveMessageV2 is the envelope used by the v2 price feeds (live and delayed
// BIST prices, US prices, BIST bid/ask). On the wire every event looks like
// {"t":"<type>","d":{...}}. Check Type before using Data:
//
//   - MessageTypePrice ("pr"): Data holds a price tick.
//   - MessageTypeHeartbeat ("heartbeat"): sent every 10 seconds with no "d",
//     so Data is the zero value.
//   - MessageTypeStateChange ("state_change"): emitted on the BIST live and
//     delayed feeds when a market or stock changes state, as
//     {"t":"state_change","d":{"marketSymbol":..,"stockSymbol":..,"state":..,"time":..}}.
//     Market-level events reach every subscriber regardless of the symbol
//     filter. That payload does not match T, so Data is the zero value; the
//     state payload is not decoded yet.
type LiveMessageV2[T any] struct {
	Type MessageType `json:"t"`
	Data T           `json:"d"`
}

// LevelSide represents the side of an orderbook level
type LevelSide string

const (
	LevelSideBid LevelSide = "bid"
	LevelSideAsk LevelSide = "ask"
)

// OrderbookLevel represents a single level in the orderbook
type OrderbookLevel struct {
	ID     int       `json:"level"`
	Side   LevelSide `json:"side"`
	Volume float64   `json:"vol"`
	Orders int       `json:"orders"`
	Price  float64   `json:"p"`
}

// OrderbookDeletedLevel represents a deleted level in the orderbook
type OrderbookDeletedLevel struct {
	ID   int       `json:"level"`
	Side LevelSide `json:"side"`
}

// BISTStockOrderBookData represents BIST stock order book data. Unlike the
// price feeds the order book stream is not wrapped in an envelope and carries
// the symbol under "symbol".
type BISTStockOrderBookData struct {
	Updated []OrderbookLevel        `json:"updated"`
	Deleted []OrderbookDeletedLevel `json:"deleted"`
	Symbol  string                  `json:"symbol"`
}

// LivePriceStream handles live price streaming for a specific region and type
type LivePriceStream[T any] struct {
	mu sync.RWMutex

	// Per-subscription state, replaced on every Subscribe. cancel stops the
	// SSE request and the forwarder; done is closed by the forwarder once it
	// has exited and closed outputChan.
	cancel     context.CancelFunc
	done       chan struct{}
	outputChan chan LivePriceResult[T]

	c            *Client
	region       Region
	priceType    LivePriceType
	symbols      []string
	closed       bool
	isSubscribed bool
}

// NewLivePriceStream creates a new LivePriceStream
func NewLivePriceStream[T any](client *Client, priceType LivePriceType, region Region) *LivePriceStream[T] {
	return &LivePriceStream[T]{
		c:         client,
		priceType: priceType,
		region:    region,
		closed:    false,
	}
}

// Subscribe subscribes to live price updates for the given symbols. Calling
// Subscribe on a stream that is already subscribed switches it to the new
// symbols: the channel previously returned by Receive is closed and a fresh
// one is created, so call Receive again afterwards.
func (s *LivePriceStream[T]) Subscribe(ctx context.Context, symbols []string) error {
	if ctx == nil {
		return fmt.Errorf("context cannot be nil")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	// Cleanup existing stream
	if err := s.cleanupExistingStream(); err != nil {
		return fmt.Errorf("failed to cleanup existing stream: %w", err)
	}

	s.symbols = symbols
	s.closed = false
	s.isSubscribed = false

	// Start streaming
	if err := s.startStreaming(ctx); err != nil {
		return fmt.Errorf("failed to start streaming: %w", err)
	}

	s.isSubscribed = true
	return nil
}

// Receive returns a channel to receive live price data. The channel is closed
// when Close is called, when the context passed to Subscribe is cancelled,
// when the server ends the stream, or when Subscribe is called again.
func (s *LivePriceStream[T]) Receive() <-chan LivePriceResult[T] {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if !s.isSubscribed || s.outputChan == nil {
		// Return a closed channel if not subscribed
		ch := make(chan LivePriceResult[T])
		close(ch)
		return ch
	}

	return s.outputChan
}

// Close closes the stream and cleanup resources
func (s *LivePriceStream[T]) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.closed {
		return nil
	}

	s.closed = true
	s.isSubscribed = false
	return s.cleanupExistingStream()
}

// cleanupExistingStream cancels the current subscription, if any, and waits
// for its forwarding goroutine to exit. The forwarder owns outputChan and
// closes it on the way out, so once this returns any consumer ranging over
// the old channel has been released and nothing can write to it anymore.
func (s *LivePriceStream[T]) cleanupExistingStream() error {
	if s.cancel != nil {
		s.cancel()
		s.cancel = nil
	}

	if s.done != nil {
		<-s.done
		s.done = nil
	}

	s.outputChan = nil
	return nil
}

// buildStreamURL builds the streaming URL for the given symbols and region
func (s *LivePriceStream[T]) buildStreamURL() string {
	streamID := uuid.New().String()
	symbolsParam := strings.Join(s.symbols, ",")

	baseURL := s.c.baseUrl
	var endpoint string

	switch {
	case s.priceType == LivePriceTypePrice && s.region == RegionTr:
		endpoint = "/api/v2/stock/price/live"
	case s.priceType == LivePriceTypeDelayedPrice:
		endpoint = "/api/v1/stock/price/delayed"
	case s.priceType == LivePriceTypeOrderBook:
		endpoint = "/api/v1/stock/orderbook/live"
	case s.priceType == LivePriceTypeBidAsk:
		endpoint = "/api/v1/stock/price/bids"
	default:
		endpoint = "/api/v2/stock/price/live"
	}

	return fmt.Sprintf("%s%s?filter=%s&region=%s&stream=%s",
		baseURL, endpoint, symbolsParam, string(s.region), streamID)
}

// startStreaming starts the SSE streaming connection
func (s *LivePriceStream[T]) startStreaming(ctx context.Context) error {
	url := s.buildStreamURL()

	ctxWithCancel, cancel := context.WithCancel(ctx)

	sseChan, _, err := sendSSERequest[T](ctxWithCancel, s.c, url)
	if err != nil {
		cancel()
		return fmt.Errorf("failed to establish SSE connection: %w", err)
	}

	outputChan := make(chan LivePriceResult[T], 100) // Buffered channel
	done := make(chan struct{})

	s.cancel = cancel
	s.outputChan = outputChan
	s.done = done

	go s.forwardData(ctxWithCancel, sseChan, outputChan, done)

	return nil
}

// forwardData forwards data from the SSE channel to outputChan until the
// subscription context is cancelled or the SSE reader closes its channel
// (the server ended the stream). It takes the channels as arguments rather
// than reading them off the struct so a later Subscribe cannot swap them
// underneath it. On exit it closes outputChan, releasing consumers ranging
// over Receive(), and then signals done.
func (s *LivePriceStream[T]) forwardData(ctx context.Context, sseChan <-chan LivePriceResult[T], outputChan chan<- LivePriceResult[T], done chan<- struct{}) {
	defer close(done)
	defer close(outputChan)
	defer func() {
		if r := recover(); r != nil {
			s.c.logger.Error("panic in forwardData", r)
		}
	}()

	for {
		select {
		case data, ok := <-sseChan:
			if !ok {
				return
			}

			select {
			case outputChan <- data:
			case <-ctx.Done():
				return
			}
		case <-ctx.Done():
			return
		}
	}
}

type BISTStockLiveData struct {
	Symbol             string  `json:"s"`
	DailyPercentChange float64 `json:"ch"`
	ClosePrice         float64 `json:"p"`
	Date               int64   `json:"d"`
}

type USStockLiveData struct {
	Symbol        string  `json:"s"`
	Price         float64 `json:"p"`
	Date          int64   `json:"d"`
	PercentChange float64 `json:"pc"`
	AmountChange  float64 `json:"ac"`
}

// ===== NEW UNIFIED STREAMING API =====

// GetLivePriceStreamForBIST creates a new live price stream for BIST stocks.
// Call Subscribe(ctx, symbols) on the returned stream to start receiving data.
func (c *Client) GetLivePriceStreamForBIST() *LivePriceStream[LiveMessageV2[BISTStockLiveData]] {
	stream := NewLivePriceStream[LiveMessageV2[BISTStockLiveData]](c, LivePriceTypePrice, RegionTr)
	return stream
}

// GetLivePriceStreamForUS creates a new live price stream for US stocks.
// Call Subscribe(ctx, symbols) on the returned stream to start receiving data.
func (c *Client) GetLivePriceStreamForUS() *LivePriceStream[LiveMessageV2[USStockLiveData]] {
	stream := NewLivePriceStream[LiveMessageV2[USStockLiveData]](c, LivePriceTypePrice, RegionUs)
	return stream
}

// GetLiveOrderBookStreamForBIST creates a new order book stream for BIST stocks.
// Call Subscribe(ctx, symbols) on the returned stream to start receiving data.
func (c *Client) GetLiveOrderBookStreamForBIST() *LivePriceStream[BISTStockOrderBookData] {
	stream := NewLivePriceStream[BISTStockOrderBookData](c, LivePriceTypeOrderBook, RegionTr)
	return stream
}

// GetDelayedPriceStreamForBIST creates a new delayed price stream for BIST stocks.
// Call Subscribe(ctx, symbols) on the returned stream to start receiving data.
func (c *Client) GetDelayedPriceStreamForBIST() *LivePriceStream[LiveMessageV2[BISTStockLiveData]] {
	stream := NewLivePriceStream[LiveMessageV2[BISTStockLiveData]](c, LivePriceTypeDelayedPrice, RegionTr)
	return stream
}

// ===== CONVENIENCE METHODS (Python-style API) =====

// CreateLivePriceStreamForBIST creates and subscribes to live price stream for BIST
func (c *Client) CreateLivePriceStreamForBIST(ctx context.Context, symbols []string) (*LivePriceStream[LiveMessageV2[BISTStockLiveData]], error) {
	stream := c.GetLivePriceStreamForBIST()
	if err := stream.Subscribe(ctx, symbols); err != nil {
		return nil, fmt.Errorf("failed to subscribe to live price stream: %w", err)
	}
	return stream, nil
}

// CreateLivePriceStreamForUS creates and subscribes to live price stream for US stocks
func (c *Client) CreateLivePriceStreamForUS(ctx context.Context, symbols []string) (*LivePriceStream[LiveMessageV2[USStockLiveData]], error) {
	stream := c.GetLivePriceStreamForUS()
	if err := stream.Subscribe(ctx, symbols); err != nil {
		return nil, fmt.Errorf("failed to subscribe to live price stream: %w", err)
	}
	return stream, nil
}

// CreateLiveOrderBookStreamForBIST creates and subscribes to order book stream for BIST
func (c *Client) CreateLiveOrderBookStreamForBIST(ctx context.Context, symbols []string) (*LivePriceStream[BISTStockOrderBookData], error) {
	stream := c.GetLiveOrderBookStreamForBIST()
	if err := stream.Subscribe(ctx, symbols); err != nil {
		return nil, fmt.Errorf("failed to subscribe to order book stream: %w", err)
	}
	return stream, nil
}

// CreateDelayedPriceStreamForBIST creates and subscribes to delayed price stream for BIST
func (c *Client) CreateDelayedPriceStreamForBIST(ctx context.Context, symbols []string) (*LivePriceStream[LiveMessageV2[BISTStockLiveData]], error) {
	stream := c.GetDelayedPriceStreamForBIST()
	if err := stream.Subscribe(ctx, symbols); err != nil {
		return nil, fmt.Errorf("failed to subscribe to delayed price stream: %w", err)
	}
	return stream, nil
}

// BISTBidAskResponse is the envelope of the bid/ask feed. It is the same
// LiveMessageV2 envelope as the other price feeds; the alias is kept so the
// bid/ask wire format cannot drift from the shared one again.
type BISTBidAskResponse = LiveMessageV2[BISTBidAskLiveData]

type BISTBidAskLiveData struct {
	Symbol string  `json:"s"`
	Ask    float64 `json:"ask"`
	Bid    float64 `json:"bid"`
	Date   int64   `json:"d"`
}

// GetLiveBidAskStreamForBIST creates a new bid/ask price stream for BIST stocks.
// Call Subscribe(ctx, symbols) on the returned stream to start receiving data.
// Passing no symbols to Subscribe means all BIST stocks will be streamed.
func (c *Client) GetLiveBidAskStreamForBIST() *LivePriceStream[BISTBidAskResponse] {
	stream := NewLivePriceStream[BISTBidAskResponse](c, LivePriceTypeBidAsk, RegionTr)
	return stream
}

// CreateLiveBidAskStreamForBIST creates and subscribes to bid/ask price stream for BIST.
func (c *Client) CreateLiveBidAskStreamForBIST(ctx context.Context, symbols []string) (*LivePriceStream[BISTBidAskResponse], error) {
	stream := c.GetLiveBidAskStreamForBIST()
	if err := stream.Subscribe(ctx, symbols); err != nil {
		return nil, fmt.Errorf("failed to subscribe to bid/ask stream: %w", err)
	}
	return stream, nil
}
