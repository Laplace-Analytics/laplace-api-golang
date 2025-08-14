package laplace

import (
	"context"
	"fmt"
	"net/http"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Industry struct {
	ID        primitive.ObjectID `json:"id"`
	Title     string             `json:"title"`
	ImageUrl  string             `json:"imageUrl"`
	AvatarUrl string             `json:"avatarUrl"`
	NumStocks int                `json:"numStocks"`
}

// GetAllIndustries retrieves all industries available for the specified region and locale.
func (c *Client) GetAllIndustries(ctx context.Context, region Region, locale Locale) ([]Industry, error) {
	endpoint := fmt.Sprintf("%s/api/v1/industry", c.baseUrl)
	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Add("locale", string(locale))
	q.Add("region", string(region))
	req.URL.RawQuery = q.Encode()

	res, err := sendRequest[[]Industry](ctx, c, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// GetIndustryDetail fetches detailed information about a specific industry including its constituent stocks.
func (c *Client) GetIndustryDetail(ctx context.Context, id string, region Region, locale Locale) (CollectionDetail, error) {
	endpoint := fmt.Sprintf("%s/api/v1/industry/%s", c.baseUrl, id)
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
