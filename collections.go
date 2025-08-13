package laplace

import (
	"context"
	"fmt"
	"net/http"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Region string

const (
	RegionTr   Region = "tr"
	RegionUs   Region = "us"
	RegionNone Region = "none"
)

type Locale string

const (
	LocaleTr   Locale = "tr"
	LocaleEn   Locale = "en"
	LocaleNone Locale = "none"
)

type Collection struct {
	ID          primitive.ObjectID `json:"id"`
	Title       string             `json:"title"`
	Region      []Region           `json:"region"`
	ImageUrl    string             `json:"imageUrl"`
	AvatarUrl   string             `json:"avatarUrl"`
	NumStocks   int                `json:"numStocks"`
	AssetClass  AssetClass         `json:"assetClass,omitempty"`
	Description string             `json:"description,omitempty"`
	Image       string             `json:"image,omitempty"`
	Order       *int               `json:"order,omitempty"`
	Status      CollectionStatus   `json:"status,omitempty"`
	MetaData    map[string]any     `json:"metaData,omitempty"`
}

type CollectionDetail struct {
	*Collection `json:",inline"`
	Stocks      []Stock `json:"stocks"`
}

type SortBy string

const (
	SortByPriceChange SortBy = "price_change"
)

// GetAllCollections retrieves all collections available for the specified region and locale.
func (c *Client) GetAllCollections(ctx context.Context, region Region, locale Locale) ([]Collection, error) {
	endpoint := fmt.Sprintf("%s/api/v1/collection", c.baseUrl)
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

// GetCollectionDetail fetches detailed information about a specific collection including its constituent stocks.
func (c *Client) GetCollectionDetail(ctx context.Context, id string, region Region, locale Locale) (CollectionDetail, error) {
	endpoint := fmt.Sprintf("%s/api/v1/collection/%s", c.baseUrl, id)
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
