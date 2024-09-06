package client

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"finfree.co/laplace/utilities"
	"github.com/sirupsen/logrus"
)

type Client struct {
	cli     *http.Client
	baseUrl string
	apiKey  string
	logger  *logrus.Logger
}

func NewClient(
	cfg utilities.LaplaceConfiguration,
	logger *logrus.Logger,
) *Client {
	return &Client{
		cli:     &http.Client{},
		baseUrl: cfg.BaseURL,
		apiKey:  cfg.APIKey,
		logger:  logger,
	}
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
		return resp, fmt.Errorf("unexpected status code: %d, body: %s", res.StatusCode, string(body))
	}

	if err := json.Unmarshal(body, &resp); err != nil {
		return resp, err
	}

	return resp, nil
}

func sendSSERequest[T any](
	ctx context.Context,
	c *Client,
	r *http.Request) (<-chan T, <-chan error, error) {
	// Set headers
	r.Header.Set("Accept", "text/event-stream")
	r.Header.Set("Cache-Control", "no-cache")
	r.Header.Set("Connection", "keep-alive")
	r.Header.Set("Authorization", "Bearer "+c.apiKey)

	// Send the request
	resp, err := c.cli.Do(r.WithContext(ctx))
	if err != nil {
		return nil, nil, err
	}

	// Check the response status
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		return nil, nil, fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(bodyBytes))
	}

	// Create a channel to send LivePriceEvents
	events := make(chan T)
	errorChan := make(chan error)

	// Start a goroutine to read the SSE stream
	go func() {
		defer resp.Body.Close()
		defer close(events)

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
	}()

	return events, errorChan, nil
}
