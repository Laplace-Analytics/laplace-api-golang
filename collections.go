package laplace

import (
	"context"
	"fmt"
	"net/http"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CollectionType string

const (
	CollectionTypeSector      CollectionType = "sector"
	CollectionTypeIndustry    CollectionType = "industry"
	CollectionTypeTheme       CollectionType = "theme"
	CollectionTypeCustomTheme CollectionType = "custom-theme"
	CollectionTypeCollection  CollectionType = "collection"
)

type Region string

const (
	RegionTr Region = "tr"
	RegionUs Region = "us"
)

type Locale string

const (
	LocaleTr Locale = "tr"
	LocaleEn Locale = "en"
)

type Collection struct {
	ID          primitive.ObjectID `json:"id"`
	Title       string             `json:"title"`
	Description string             `json:"description"`
	Region      []Region           `json:"region"`
	AssetClass  string             `json:"assetClass"`
	ImageUrl    string             `json:"imageUrl"`
	AvatarUrl   string             `json:"avatarUrl"`
	NumStocks   int                `json:"numStocks"`

	Image    string           `json:"image"`
	Order    int              `json:"order"`
	Status   CollectionStatus `json:"status,omitempty"`
	MetaData map[string]any   `json:"metaData,omitempty"`
}

type CollectionDetail struct {
	*Collection `json:",inline"`
	Stocks      []Stock `json:"stocks"`
}

func (c *Client) getAllCollections(ctx context.Context, collectionType CollectionType, region Region, locale Locale) ([]Collection, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/api/v1/%s", c.baseUrl, collectionType), nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	if region != "" {
		q.Add("region", string(region))
	}
	q.Add("locale", string(locale))
	req.URL.RawQuery = q.Encode()

	resp, err := sendRequest[[]Collection](ctx, c, req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

type SortBy string

const (
	SortByPriceChange SortBy = "price_change"
)

func (c *Client) getCollectionDetail(ctx context.Context, id string, collectionType CollectionType, region Region, locale Locale, sortBy SortBy) (CollectionDetail, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/api/v1/%s/%s", c.baseUrl, collectionType, id), nil)
	if err != nil {
		return CollectionDetail{}, err
	}

	q := req.URL.Query()
	if region != "" {
		q.Add("region", string(region))
	}
	q.Add("locale", string(locale))
	if sortBy != "" {
		q.Add("sortBy", string(sortBy))
	}
	req.URL.RawQuery = q.Encode()

	resp, err := sendRequest[CollectionDetail](ctx, c, req)
	if err != nil {
		return CollectionDetail{}, err
	}

	return resp, nil
}

func (c *Client) GetAllSectors(ctx context.Context, region Region, locale Locale) ([]Collection, error) {
	return c.getAllCollections(ctx, CollectionTypeSector, region, locale)
}

func (c *Client) GetAllIndustries(ctx context.Context, region Region, locale Locale) ([]Collection, error) {
	return c.getAllCollections(ctx, CollectionTypeIndustry, region, locale)
}

func (c *Client) GetAllThemes(ctx context.Context, region Region, locale Locale) ([]Collection, error) {
	return c.getAllCollections(ctx, CollectionTypeTheme, region, locale)
}

func (c *Client) GetAllCollections(ctx context.Context, region Region, locale Locale) ([]Collection, error) {
	return c.getAllCollections(ctx, CollectionTypeCollection, region, locale)
}

func (c *Client) GetSectorDetail(ctx context.Context, id string, region Region, locale Locale) (CollectionDetail, error) {
	return c.getCollectionDetail(ctx, id, CollectionTypeSector, region, locale, "")
}

func (c *Client) GetIndustryDetail(ctx context.Context, id string, region Region, locale Locale) (CollectionDetail, error) {
	return c.getCollectionDetail(ctx, id, CollectionTypeIndustry, region, locale, "")
}

func (c *Client) GetThemeDetail(ctx context.Context, id string, region Region, locale Locale) (CollectionDetail, error) {
	return c.getCollectionDetail(ctx, id, CollectionTypeTheme, region, locale, "")
}

func (c *Client) GetCollectionDetail(ctx context.Context, id string, region Region, locale Locale) (CollectionDetail, error) {
	return c.getCollectionDetail(ctx, id, CollectionTypeCollection, region, locale, "")
}
