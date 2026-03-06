package laplace

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

type NewsType string

const (
	NewsTypeBriefs    NewsType = "briefs"
	NewsTypeBloomberg NewsType = "bloomberg"
	NewsTypeFDA       NewsType = "fda"
	NewsTypeReuters   NewsType = "reuters"
)

type NewsOrderBy string

const (
	NewsOrderByTimestamp NewsOrderBy = "timestamp"
)

type NewsHighlights struct {
	Consumer                []string `json:"consumer"`
	EnergyAndUtilities      []string `json:"energyAndUtilities"`
	Finance                 []string `json:"finance"`
	Healthcare              []string `json:"healthcare"`
	IndustrialsAndMaterials []string `json:"industrialsAndMaterials"`
	Tech                    []string `json:"tech"`
	Other                   []string `json:"other"`
}

type NewsPublisher struct {
	Name    string  `json:"name"`
	LogoUrl *string `json:"logoUrl"`
}

type NewsTicker struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Symbol string `json:"symbol,omitempty"`
}

type NewsCategories struct {
	Name         string  `json:"name"`
	NewsCount    int64   `json:"newsCount"`
	CategoryType *string `json:"categoryType,omitempty"`
	MeanType     *int64  `json:"meanType,omitempty"`
}

type NewsSector struct {
	Name         string  `json:"name"`
	NewsCount    int64   `json:"newsCount"`
	CategoryType *string `json:"categoryType,omitempty"`
	MeanType     *int64  `json:"meanType,omitempty"`
}

type NewsContent struct {
	Title           string   `json:"title"`
	Description     string   `json:"description"`
	Content         []string `json:"content"`
	Summary         []string `json:"summary"`
	InvestorInsight string   `json:"investorInsight"`
}

type NewsIndustry struct {
	Name     string `json:"name"`
	MeanType int64  `json:"meanType"`
}

type News struct {
	URL            string          `json:"url"`
	ImageUrl       string          `json:"imageUrl"`
	Timestamp      time.Time       `json:"timestamp"`
	PublisherUrl   string          `json:"publisherUrl"`
	Publisher      NewsPublisher   `json:"publisher"`
	RelatedTickers []NewsTicker    `json:"relatedTickers"`
	QualityScore   int64           `json:"qualityScore"`
	CreatedAt      time.Time       `json:"createdAt"`
	Tickers        []NewsTicker    `json:"tickers,omitempty"`
	Categories     *NewsCategories `json:"categories,omitempty"`
	Sectors        *NewsSector     `json:"sectors,omitempty"`
	Content        *NewsContent    `json:"content,omitempty"`
	Industries     *NewsIndustry   `json:"industries,omitempty"`
}

// GetNewsHighlights retrieves news highlights categorized by sector for the specified region and locale.
func (c *Client) GetNewsHighlights(ctx context.Context, region Region, locale Locale) (*NewsHighlights, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/api/v1/news/highlights", c.baseUrl), nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Add("region", string(region))
	q.Add("locale", string(locale))
	req.URL.RawQuery = q.Encode()

	resp, err := sendRequest[NewsHighlights](ctx, c, req)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

type GetNewsParams struct {
	Region           Region
	Locale           Locale
	NewsType         NewsType
	Page             *int
	Size             *int
	OrderBy          NewsOrderBy
	OrderByDirection SortDirection
	ExtraFilters     string
}

// GetNews retrieves a paginated list of news articles with optional filtering by type, ordering, and extra filters.
func (c *Client) GetNews(ctx context.Context, params GetNewsParams) (*PaginatedResponse[News], error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/api/v1/news", c.baseUrl), nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Add("region", string(params.Region))
	q.Add("locale", string(params.Locale))
	if params.NewsType != "" {
		q.Add("newsType", string(params.NewsType))
	}
	if params.Page != nil {
		q.Add("page", strconv.Itoa(*params.Page))
	}
	if params.Size != nil {
		q.Add("size", strconv.Itoa(*params.Size))
	}
	if params.OrderBy != "" {
		q.Add("orderBy", string(params.OrderBy))
	}
	if params.OrderByDirection != "" {
		q.Add("orderByDirection", string(params.OrderByDirection))
	}
	if params.ExtraFilters != "" {
		q.Add("extraFilters", params.ExtraFilters)
	}
	req.URL.RawQuery = q.Encode()

	resp, err := sendRequest[PaginatedResponse[News]](ctx, c, req)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

// NewsStreamResult is the result type for news streams
type NewsStreamResult struct {
	Data  []News
	Error error
}

func sendNewsSSERequest(
	ctx context.Context,
	c *Client,
	url string,
) (<-chan NewsStreamResult, func(), error) {

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, nil, err
	}

	// Set headers
	req.Header.Set("Accept", "text/event-stream")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	// Send the request
	resp, err := c.cli.Do(req.WithContext(ctx))
	if err != nil {
		return nil, nil, err
	}

	// Check the response status
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		defer resp.Body.Close()
		return nil, nil, fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(bodyBytes))
	}

	// Create a single channel for results
	results := make(chan NewsStreamResult)

	ctxWithCancel, cancel := context.WithCancel(ctx)

	// Start a goroutine to read the SSE stream
	go func() {
		defer resp.Body.Close()
		defer close(results)
		defer cancel()

		scanner := bufio.NewScanner(resp.Body)
		for scanner.Scan() {
			select {
			case <-ctxWithCancel.Done():
				return
			default:
				line := scanner.Text()
				if strings.HasPrefix(line, "data:") {
					data := strings.TrimPrefix(line, "data:")
					var event []News
					if err := json.Unmarshal([]byte(data), &event); err != nil {
						results <- NewsStreamResult{Error: fmt.Errorf("error unmarshalling event: %w", err)}
						continue
					}
					results <- NewsStreamResult{Data: event}
				}
			}
		}

		if err := scanner.Err(); err != nil {
			results <- NewsStreamResult{Error: fmt.Errorf("error reading SSE stream: %w", err)}
		}
	}()

	return results, cancel, nil
}

// NewsStream handles live news streaming for a specific locale
type NewsStream struct {
	mu           sync.RWMutex
	ctx          context.Context
	cancel       context.CancelFunc
	sseChan      <-chan NewsStreamResult
	outputChan   chan NewsStreamResult
	c            *Client
	locale       Locale
	closed       bool
	isSubscribed bool
}

// Subscribe starts receiving news from the stream
func (s *NewsStream) Subscribe(ctx context.Context) error {
	if ctx == nil {
		return fmt.Errorf("context cannot be nil")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	// Cleanup existing stream
	if err := s.cleanupExistingStream(); err != nil {
		return fmt.Errorf("failed to cleanup existing stream: %w", err)
	}

	s.outputChan = make(chan NewsStreamResult, 100) // Buffered channel
	s.closed = false
	s.ctx = ctx

	// Start streaming
	if err := s.startStreaming(); err != nil {
		return fmt.Errorf("failed to start streaming: %w", err)
	}

	s.isSubscribed = true
	return nil
}

// Receive returns a channel to receive news data
func (s *NewsStream) Receive() <-chan NewsStreamResult {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if !s.isSubscribed {
		// Return a closed channel if not subscribed
		ch := make(chan NewsStreamResult)
		close(ch)
		return ch
	}

	return s.outputChan
}

// Close closes the stream and cleanup resources
func (s *NewsStream) Close() error {
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
func (s *NewsStream) cleanupExistingStream() error {
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

// startStreaming starts the SSE streaming connection
func (s *NewsStream) startStreaming() error {
	url := fmt.Sprintf("%s/api/v1/news/stream?locale=%s", s.c.baseUrl, string(s.locale))

	ctxWithCancel, cancel := context.WithCancel(s.ctx)
	s.cancel = cancel

	channel, _, err := sendNewsSSERequest(ctxWithCancel, s.c, url)
	if err != nil {
		return fmt.Errorf("failed to establish SSE connection: %w", err)
	}

	s.sseChan = channel
	go s.forwardData()

	return nil
}

// forwardData forwards data from SSE channel to output channel
func (s *NewsStream) forwardData() {
	defer func() {
		if r := recover(); r != nil {
			s.c.logger.Error("panic in news stream forwardData", r)
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

// GetNewsStream creates a new news stream.
// Call Subscribe(ctx) on the returned stream to start receiving data.
func (c *Client) GetNewsStream(locale Locale) *NewsStream {
	return &NewsStream{
		c:      c,
		locale: locale,
		closed: false,
	}
}

// CreateNewsStream creates and subscribes to a news stream.
func (c *Client) CreateNewsStream(ctx context.Context, locale Locale) (*NewsStream, error) {
	stream := c.GetNewsStream(locale)
	if err := stream.Subscribe(ctx); err != nil {
		return nil, fmt.Errorf("failed to subscribe to news stream: %w", err)
	}
	return stream, nil
}
