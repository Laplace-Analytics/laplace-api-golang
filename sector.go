package laplace

import (
	"context"
	"fmt"
	"net/http"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Sector struct {
	ID        primitive.ObjectID `json:"id"`
	Title     string             `json:"title"`
	ImageUrl  string             `json:"imageUrl"`
	AvatarUrl string             `json:"avatarUrl"`
	NumStocks int                `json:"numStocks"`
}

// GetAllSectors retrieves all sectors available for the specified region and locale.
func (c *Client) GetAllSectors(ctx context.Context, region Region, locale Locale) ([]Sector, error) {
	endpoint := fmt.Sprintf("%s/api/v1/sector", c.baseUrl)
	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Add("locale", string(locale))
	q.Add("region", string(region))
	req.URL.RawQuery = q.Encode()

	res, err := sendRequest[[]Sector](ctx, c, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}
