package laplace

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

type NewsType string

const (
	NewsTypeBriefs    NewsType = "briefs"
	NewsTypeBloomberg NewsType = "bloomberg"
	NewsTypeFDA       NewsType = "fda"
	NewsTypeReuters   NewsType = "reuters"
)

type NewsOrderBy string

const (
	NewsOrderByTimestamp NewsOrderBy = "timestamp"
)

type NewsHighlights struct {
	Consumer                []string `json:"consumer"`
	EnergyAndUtilities      []string `json:"energyAndUtilities"`
	Finance                 []string `json:"finance"`
	Healthcare              []string `json:"healthcare"`
	IndustrialsAndMaterials []string `json:"industrialsAndMaterials"`
	Tech                    []string `json:"tech"`
	Other                   []string `json:"other"`
}

type NewsPublisher struct {
	Name    string  `json:"name"`
	LogoUrl *string `json:"logoUrl"`
}

type NewsTicker struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Symbol string `json:"symbol,omitempty"`
}

type NewsCategories struct {
	Name         string  `json:"name"`
	NewsCount    int64   `json:"newsCount"`
	CategoryType *string `json:"categoryType,omitempty"`
	MeanType     *int64  `json:"meanType,omitempty"`
}

type NewsSector struct {
	Name         string  `json:"name"`
	NewsCount    int64   `json:"newsCount"`
	CategoryType *string `json:"categoryType,omitempty"`
	MeanType     *int64  `json:"meanType,omitempty"`
}

type NewsContent struct {
	Title           string   `json:"title"`
	Description     string   `json:"description"`
	Content         []string `json:"content"`
	Summary         []string `json:"summary"`
	InvestorInsight string   `json:"investorInsight"`
}

type NewsIndustry struct {
	Name     string `json:"name"`
	MeanType int64  `json:"meanType"`
}

type News struct {
	URL            string          `json:"url"`
	ImageUrl       string          `json:"imageUrl"`
	Timestamp      time.Time       `json:"timestamp"`
	PublisherUrl   string          `json:"publisherUrl"`
	Publisher      NewsPublisher   `json:"publisher"`
	RelatedTickers []NewsTicker    `json:"relatedTickers"`
	QualityScore   int64           `json:"qualityScore"`
	CreatedAt      time.Time       `json:"createdAt"`
	Tickers        []NewsTicker    `json:"tickers,omitempty"`
	Categories     *NewsCategories `json:"categories,omitempty"`
	Sectors        *NewsSector     `json:"sectors,omitempty"`
	Content        *NewsContent    `json:"content,omitempty"`
	Industries     *NewsIndustry   `json:"industries,omitempty"`
}

// GetNewsHighlights retrieves news highlights categorized by sector for the specified region and locale.
func (c *Client) GetNewsHighlights(ctx context.Context, region Region, locale Locale) (*NewsHighlights, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/api/v1/news/highlights", c.baseUrl), nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Add("region", string(region))
	q.Add("locale", string(locale))
	req.URL.RawQuery = q.Encode()

	resp, err := sendRequest[NewsHighlights](ctx, c, req)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

type GetNewsParams struct {
	Region           Region
	Locale           Locale
	NewsType         NewsType
	Page             *int
	Size             *int
	OrderBy          NewsOrderBy
	OrderByDirection SortDirection
	ExtraFilters     string
}

// GetNews retrieves a paginated list of news articles with optional filtering by type, ordering, and extra filters.
func (c *Client) GetNews(ctx context.Context, params GetNewsParams) (*PaginatedResponse[News], error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/api/v1/news", c.baseUrl), nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Add("region", string(params.Region))
	q.Add("locale", string(params.Locale))
	if params.NewsType != "" {
		q.Add("newsType", string(params.NewsType))
	}
	if params.Page != nil {
		q.Add("page", strconv.Itoa(*params.Page))
	}
	if params.Size != nil {
		q.Add("size", strconv.Itoa(*params.Size))
	}
	if params.OrderBy != "" {
		q.Add("orderBy", string(params.OrderBy))
	}
	if params.OrderByDirection != "" {
		q.Add("orderByDirection", string(params.OrderByDirection))
	}
	if params.ExtraFilters != "" {
		q.Add("extraFilters", params.ExtraFilters)
	}
	req.URL.RawQuery = q.Encode()

	resp, err := sendRequest[PaginatedResponse[News]](ctx, c, req)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}
