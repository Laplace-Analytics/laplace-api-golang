package client

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PaginationPageSize int

const (
	PageSize10 PaginationPageSize = 10
	PageSize20 PaginationPageSize = 20
	PageSize50 PaginationPageSize = 50
)

type SearchType string

const (
	SearchTypeStock      SearchType = "stock"
	SearchTypeCollection SearchType = "collection"
	SearchTypeSector     SearchType = "sector"
	SearchTypeIndustry   SearchType = "industry"
)

type SearchResponse struct {
	Stocks      []SearchResponseStock      `json:"stocks"`
	Collections []SearchResponseCollection `json:"collections"`
	Sectors     []SearchResponseCollection `json:"sectors"`
	Industries  []SearchResponseCollection `json:"industries"`
}

type SearchResponseStock struct {
	ID         primitive.ObjectID `json:"id"`
	Name       string             `json:"name"`
	Symbol     string             `json:"title"`
	Region     string             `json:"region"`
	AssetClass string             `json:"assetType"`
	AssetType  string             `json:"type"`
}

type SearchResponseCollection struct {
	ID         primitive.ObjectID `json:"id"`
	Title      string             `json:"title"`
	Region     []string           `json:"region"`
	AssetClass string             `json:"assetClass"`
	ImageUrl   string             `json:"imageUrl"`
	AvatarUrl  string             `json:"avatarUrl"`
}

func (c *Client) Search(ctx context.Context, query string, types []SearchType, region Region, locale Locale, page int, pageSize PaginationPageSize) (*SearchResponse, error) {
	typesStr := strings.Join(lo.Map(types, func(key SearchType, _ int) string {
		return string(key)
	}), ",")

	/// build path properly with query params url encoded
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/api/v1/search", c.baseUrl), nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Add("filter", query)
	q.Add("types", typesStr)
	q.Add("region", string(region))
	q.Add("locale", string(locale))
	q.Add("page", fmt.Sprintf("%d", page))
	q.Add("size", fmt.Sprintf("%d", pageSize))
	req.URL.RawQuery = q.Encode()

	resp, err := sendRequest[SearchResponse](ctx, c, req)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}
