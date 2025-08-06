package laplace

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/google/uuid"
)

type LivePriceClient[T any] interface {
	Close() error
	Receive() <-chan LivePriceResult[T]
	Subscribe(ctx context.Context, symbols []string) error
}

type livePriceClient[T any] struct {
	mu         sync.RWMutex
	ctx        context.Context
	cancel     context.CancelFunc
	sseChan    <-chan LivePriceResult[T]
	outputChan chan LivePriceResult[T]
	c          *Client
	region     Region
	symbols    []string
	closed     bool
}

func (c *livePriceClient[T]) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.closed {
		return nil
	}

	c.closed = true
	c.cancel()
	close(c.outputChan)

	return nil
}

func (c *livePriceClient[T]) Receive() <-chan LivePriceResult[T] {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.outputChan
}

func (c *livePriceClient[T]) Subscribe(ctx context.Context, symbols []string) error {
	if ctx == nil {
		return fmt.Errorf("context cannot be nil")
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	c.cancel()

	newClient, err := GetLivePrice[T](c.c, ctx, symbols, c.region)
	if err != nil {
		return fmt.Errorf("failed to create live price client: %w", err)
	}

	c.ctx = ctx
	c.cancel = func() { newClient.Close() }
	c.sseChan = newClient.Receive()
	c.symbols = symbols
	c.closed = false
	go c.forwardData()

	return nil
}

func (c *livePriceClient[T]) forwardData() {
	defer func() {
		if r := recover(); r != nil {
			c.c.logger.Error("panic in forwardData", r)
		}
	}()

	for {
		select {
		case data, ok := <-c.sseChan:
			if !ok {
				return
			}

			select {
			case c.outputChan <- data:
			case <-c.ctx.Done():
				return
			}
		case <-c.ctx.Done():
			return
		}
	}
}

func GetLivePrice[T any](c *Client, ctx context.Context, symbols []string, region Region) (LivePriceClient[T], error) {
	if c == nil {
		return nil, fmt.Errorf("client cannot be nil")
	}
	if ctx == nil {
		return nil, fmt.Errorf("context cannot be nil")
	}

	streamID := uuid.New().String()
	url := fmt.Sprintf("%s/api/v1/stock/price/live?filter=%s&region=%s&stream=%s",
		c.baseUrl, strings.Join(symbols, ","), string(region), streamID)

	channel, cancelFunc, err := sendSSERequest[T](ctx, c, url)
	if err != nil {
		return nil, fmt.Errorf("failed to establish SSE connection: %w", err)
	}

	client := &livePriceClient[T]{
		ctx:        ctx,
		sseChan:    channel,
		c:          c,
		region:     region,
		symbols:    symbols,
		outputChan: make(chan LivePriceResult[T]),
		closed:     false,
		cancel:     cancelFunc,
	}

	go client.forwardData()

	return client, nil
}

type BISTStockLiveData struct {
	Symbol             string  `json:"s"`
	DailyPercentChange float64 `json:"ch"`
	ClosePrice         float64 `json:"p"`
	Date               int64   `json:"d"`
}

// GetLivePriceForBIST streams real-time price data for BIST (Turkish) stock symbols via Server-Sent Events.
// Sending no symbols means all BIST stocks will be streamed.
func (c *Client) GetLivePriceForBIST(ctx context.Context, symbols []string) (LivePriceClient[BISTStockLiveData], error) {
	return GetLivePrice[BISTStockLiveData](c, ctx, symbols, RegionTr)
}

type USStockLiveData struct {
	Symbol string  `json:"s"`
	Price  float64 `json:"p"`
	Date   int64   `json:"d"`
}

// GetLivePriceForUS streams real-time price data for US stock symbols via Server-Sent Events.
// Sending no symbols means all US stocks will be streamed.
func (c *Client) GetLivePriceForUS(ctx context.Context, symbols []string) (LivePriceClient[USStockLiveData], error) {
	return GetLivePrice[USStockLiveData](c, ctx, symbols, RegionUs)
}
