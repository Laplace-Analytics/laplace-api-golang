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
)

// MessageType represents the type of message in live data streams
type MessageType string

const (
	MessageTypePrice       MessageType = "pr"
	MessageTypeStateChange MessageType = "state_change"
	MessageTypeHeartbeat   MessageType = "heartbeat"
	MessageTypeOrderbook   MessageType = "ob"
)

// LiveMessageV2 is a generic wrapper for live price messages
type LiveMessageV2[T any] struct {
	Data   T           `json:"data"`
	Symbol string      `json:"symbol"`
	Type   MessageType `json:"type"`
}

// LevelSide represents the side of an orderbook level
type LevelSide string

const (
	LevelSideBid LevelSide = "bid"
	LevelSideAsk LevelSide = "ask"
)

// OrderbookLevel represents a single level in the orderbook
type OrderbookLevel struct {
	ID    int       `json:"level"`
	Side  LevelSide `json:"side"`
	Price float64   `json:"price"`
	Size  float64   `json:"size"`
}

// OrderbookDeletedLevel represents a deleted level in the orderbook
type OrderbookDeletedLevel struct {
	ID   int       `json:"level"`
	Side LevelSide `json:"side"`
}

// BISTStockOrderBookData represents BIST stock order book data
type BISTStockOrderBookData struct {
	Updated []OrderbookLevel        `json:"updated"`
	Deleted []OrderbookDeletedLevel `json:"deleted"`
	Symbol  string                  `json:"s"`
}

type LivePriceClient[T any] interface {
	Close() error
	Receive() <-chan LivePriceResult[T]
	Subscribe(ctx context.Context, symbols []string) error
}

// LivePriceStream handles live price streaming for a specific region and type
type LivePriceStream[T any] struct {
	mu           sync.RWMutex
	ctx          context.Context
	cancel       context.CancelFunc
	sseChan      <-chan LivePriceResult[T]
	outputChan   chan LivePriceResult[T]
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

// Subscribe subscribes to live price updates for given symbols
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
	s.outputChan = make(chan LivePriceResult[T], 100) // Buffered channel
	s.closed = false
	s.ctx = ctx

	// Start streaming
	if err := s.startStreaming(); err != nil {
		return fmt.Errorf("failed to start streaming: %w", err)
	}

	s.isSubscribed = true
	return nil
}

// Receive returns a channel to receive live price data
func (s *LivePriceStream[T]) Receive() <-chan LivePriceResult[T] {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if !s.isSubscribed {
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

// cleanupExistingStream cancels and cleans up existing streaming task
func (s *LivePriceStream[T]) cleanupExistingStream() error {
	if s.cancel != nil {
		s.cancel()
		s.cancel = nil
	}

	if s.outputChan != nil {
		close(s.outputChan)
		s.outputChan = nil
	}

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
	default:
		endpoint = "/api/v2/stock/price/live"
	}

	return fmt.Sprintf("%s%s?filter=%s&region=%s&stream=%s",
		baseURL, endpoint, symbolsParam, string(s.region), streamID)
}

// startStreaming starts the SSE streaming connection
func (s *LivePriceStream[T]) startStreaming() error {
	url := s.buildStreamURL()

	ctxWithCancel, cancel := context.WithCancel(s.ctx)
	s.cancel = cancel

	channel, _, err := sendSSERequest[T](ctxWithCancel, s.c, url)
	if err != nil {
		return fmt.Errorf("failed to establish SSE connection: %w", err)
	}

	s.sseChan = channel
	go s.forwardData()

	return nil
}

// forwardData forwards data from SSE channel to output channel
func (s *LivePriceStream[T]) forwardData() {
	defer func() {
		if r := recover(); r != nil {
			s.c.logger.Error("panic in forwardData", r)
		}
	}()

	for {
		select {
		case data, ok := <-s.sseChan:
			if !ok {
				return
			}

			s.mu.RLock()
			outputChan := s.outputChan
			closed := s.closed
			s.mu.RUnlock()

			if closed || outputChan == nil {
				return
			}

			select {
			case outputChan <- data:
			case <-s.ctx.Done():
				return
			}
		case <-s.ctx.Done():
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
	Symbol string  `json:"s"`
	Price  float64 `json:"p"`
	Date   int64   `json:"d"`
}

// ===== NEW UNIFIED STREAMING API =====

// GetLivePriceStreamForBIST creates a new live price stream for BIST stocks
func (c *Client) GetLivePriceStreamForBIST(symbols []string) *LivePriceStream[LiveMessageV2[BISTStockLiveData]] {
	stream := NewLivePriceStream[LiveMessageV2[BISTStockLiveData]](c, LivePriceTypePrice, RegionTr)
	return stream
}

// GetLivePriceStreamForUS creates a new live price stream for US stocks
func (c *Client) GetLivePriceStreamForUS(symbols []string) *LivePriceStream[USStockLiveData] {
	stream := NewLivePriceStream[USStockLiveData](c, LivePriceTypePrice, RegionUs)
	return stream
}

// GetLiveOrderBookStreamForBIST creates a new order book stream for BIST stocks
func (c *Client) GetLiveOrderBookStreamForBIST(symbols []string) *LivePriceStream[BISTStockOrderBookData] {
	stream := NewLivePriceStream[BISTStockOrderBookData](c, LivePriceTypeOrderBook, RegionTr)
	return stream
}

// GetDelayedPriceStreamForBIST creates a new delayed price stream for BIST stocks
func (c *Client) GetDelayedPriceStreamForBIST(symbols []string) *LivePriceStream[LiveMessageV2[BISTStockLiveData]] {
	stream := NewLivePriceStream[LiveMessageV2[BISTStockLiveData]](c, LivePriceTypeDelayedPrice, RegionTr)
	return stream
}

// ===== CONVENIENCE METHODS (Python-style API) =====

// CreateLivePriceStreamForBIST creates and subscribes to live price stream for BIST
func (c *Client) CreateLivePriceStreamForBIST(ctx context.Context, symbols []string) (*LivePriceStream[LiveMessageV2[BISTStockLiveData]], error) {
	stream := c.GetLivePriceStreamForBIST(symbols)
	if err := stream.Subscribe(ctx, symbols); err != nil {
		return nil, fmt.Errorf("failed to subscribe to live price stream: %w", err)
	}
	return stream, nil
}

// CreateLivePriceStreamForUS creates and subscribes to live price stream for US stocks
func (c *Client) CreateLivePriceStreamForUS(ctx context.Context, symbols []string) (*LivePriceStream[USStockLiveData], error) {
	stream := c.GetLivePriceStreamForUS(symbols)
	if err := stream.Subscribe(ctx, symbols); err != nil {
		return nil, fmt.Errorf("failed to subscribe to live price stream: %w", err)
	}
	return stream, nil
}

// CreateLiveOrderBookStreamForBIST creates and subscribes to order book stream for BIST
func (c *Client) CreateLiveOrderBookStreamForBIST(ctx context.Context, symbols []string) (*LivePriceStream[BISTStockOrderBookData], error) {
	stream := c.GetLiveOrderBookStreamForBIST(symbols)
	if err := stream.Subscribe(ctx, symbols); err != nil {
		return nil, fmt.Errorf("failed to subscribe to order book stream: %w", err)
	}
	return stream, nil
}

// CreateDelayedPriceStreamForBIST creates and subscribes to delayed price stream for BIST
func (c *Client) CreateDelayedPriceStreamForBIST(ctx context.Context, symbols []string) (*LivePriceStream[LiveMessageV2[BISTStockLiveData]], error) {
	stream := c.GetDelayedPriceStreamForBIST(symbols)
	if err := stream.Subscribe(ctx, symbols); err != nil {
		return nil, fmt.Errorf("failed to subscribe to delayed price stream: %w", err)
	}
	return stream, nil
}
