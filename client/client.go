package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

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
