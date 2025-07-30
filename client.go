package laplace

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/sirupsen/logrus"
)

type Client struct {
	cli     *http.Client
	baseUrl string
	apiKey  string
	logger  *logrus.Logger
}

type clientOption func(*Client)

// WithLogger configures the client to use a custom logger instead of the default one.
func WithLogger(logger *logrus.Logger) clientOption {
	return func(c *Client) {
		c.logger = logger
	}
}

// NewClient creates a new Laplace API client with the provided configuration and optional settings.
func NewClient(
	cfg LaplaceConfiguration,
	opts ...clientOption,
) (*Client, error) {
	cfg.ApplyDefaults()
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	defaultLogger := logrus.New()
	defaultLogger.SetLevel(logrus.DebugLevel)
	defaultLogger.Out = io.Discard

	c := &Client{
		cli:     &http.Client{},
		baseUrl: cfg.BaseURL,
		apiKey:  cfg.APIKey,
		logger:  defaultLogger,
	}

	for _, opt := range opts {
		opt(c)
	}

	return c, nil
}

func sendRequest[T any](
	ctx context.Context,
	c *Client,
	r *http.Request,
) (T, error) {
	var resp T

	r.Header.Set("Authorization", "Bearer "+c.apiKey)

	res, err := c.cli.Do(r.WithContext(ctx))
	if err != nil {
		return resp, err
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return resp, err
	}

	if res.StatusCode != http.StatusOK {
		var msg LaplaceHTTPErrorMsg
		if err := json.Unmarshal(body, &msg); err != nil {
			return resp, err
		}

		return resp, getLaplaceError(&LaplaceHTTPError{
			HTTPStatus: res.StatusCode,
			Message:    msg,
		})
	}

	if err := json.Unmarshal(body, &resp); err != nil {
		return resp, err
	}

	return resp, nil
}

type LivePriceResult[T any] struct {
	Data  T
	Error error
}

func sendSSERequest[T any](
	ctx context.Context,
	c *Client,
	url string,
) (<-chan LivePriceResult[T], func(), error) {

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
	results := make(chan LivePriceResult[T])

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
					var event T
					if err := json.Unmarshal([]byte(data), &event); err != nil {
						results <- LivePriceResult[T]{Error: fmt.Errorf("error unmarshalling event: %w", err)}
						continue
					}
					results <- LivePriceResult[T]{Data: event}
				}
			}
		}

		if err := scanner.Err(); err != nil {
			results <- LivePriceResult[T]{Error: fmt.Errorf("error reading SSE stream: %w", err)}
		}
	}()

	return results, cancel, nil
}
