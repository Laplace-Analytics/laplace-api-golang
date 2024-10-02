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

func WithLogger(logger *logrus.Logger) clientOption {
	return func(c *Client) {
		c.logger = logger
	}
}

func NewClient(
	cfg LaplaceConfiguration,
	opts ...clientOption,
) *Client {

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

	return c
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
		httpErr := &LaplaceHTTPError{
			HTTPStatus: res.StatusCode,
			Message:    string(body),
		}
		return resp, WrapError(httpErr)
	}

	if err := json.Unmarshal(body, &resp); err != nil {
		return resp, err
	}

	return resp, nil
}

func sendSSERequest[T any](
	ctx context.Context,
	c *Client,
	r *http.Request) (<-chan T, <-chan error, func(), error) {
	// Set headers
	r.Header.Set("Accept", "text/event-stream")
	r.Header.Set("Cache-Control", "no-cache")
	r.Header.Set("Connection", "keep-alive")
	r.Header.Set("Authorization", "Bearer "+c.apiKey)

	// Send the request
	resp, err := c.cli.Do(r.WithContext(ctx))
	if err != nil {
		return nil, nil, nil, err
	}

	// Check the response status
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		return nil, nil, nil, fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(bodyBytes))
	}

	// Create a channel to send LivePriceEvents
	events := make(chan T)
	errorChan := make(chan error)

	// Create a new context with cancellation
	ctxWithCancel, cancel := context.WithCancel(ctx)

	// Start a goroutine to read the SSE stream
	go func() {
		defer resp.Body.Close()
		defer close(events)
		defer close(errorChan)

		scanner := bufio.NewScanner(resp.Body)
		for scanner.Scan() {
			line := scanner.Text()
			if strings.HasPrefix(line, "data:") {
				data := strings.TrimPrefix(line, "data:")
				var event T
				if err := json.Unmarshal([]byte(data), &event); err != nil {
					errorChan <- fmt.Errorf("error unmarshalling event: %w", err)
					continue
				}
				events <- event
			}
		}

		if err := scanner.Err(); err != nil {
			errorChan <- fmt.Errorf("error reading SSE stream: %w", err)
		}

		// Add a select statement to handle cancellation
		select {
		case <-ctxWithCancel.Done():
			return
		default:
			// Continue with the existing logic
		}
	}()

	// Return the channels, a cancellation function, and error
	return events, errorChan, cancel, nil
}
