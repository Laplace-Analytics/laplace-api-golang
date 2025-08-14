package laplace

import (
	"context"
	"fmt"
	"net/http"
)

// GetAllThemes retrieves all investment themes available for the specified region and locale.
func (c *Client) GetAllThemes(ctx context.Context, region Region, locale Locale) ([]Collection, error) {
	endpoint := fmt.Sprintf("%s/api/v1/theme", c.baseUrl)
	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Add("region", string(region))
	q.Add("locale", string(locale))
	req.URL.RawQuery = q.Encode()

	res, err := sendRequest[[]Collection](ctx, c, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// GetThemeDetail fetches detailed information about a specific investment theme including its constituent stocks.
func (c *Client) GetThemeDetail(ctx context.Context, id string, region Region, locale Locale) (CollectionDetail, error) {
	endpoint := fmt.Sprintf("%s/api/v1/theme/%s", c.baseUrl, id)
	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return CollectionDetail{}, err
	}

	q := req.URL.Query()
	q.Add("region", string(region))
	q.Add("locale", string(locale))
	req.URL.RawQuery = q.Encode()

	res, err := sendRequest[CollectionDetail](ctx, c, req)
	if err != nil {
		return CollectionDetail{}, err
	}

	return res, nil
}
