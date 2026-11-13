package laplace

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"
)

type NewsType string

const (
	NewsTypeBriefs    NewsType = "briefs"
	NewsTypeBloomberg NewsType = "bloomberg"
	NewsTypeFDA       NewsType = "fda"
	NewsTypeReuters   NewsType = "reuters"
)

// NewsOrderBy is a field the news endpoints can sort by. Any string is accepted, so values added
// to the API later can be passed as NewsOrderBy("…").
type NewsOrderBy string

const (
	// Deprecated: timestamp is the "best known time". Use NewsOrderByPublishedAt on v2.
	NewsOrderByTimestamp NewsOrderBy = "timestamp"
	// NewsOrderByQualityScore is the older spelling of NewsOrderByQualityScoreV2. Both are accepted.
	NewsOrderByQualityScore NewsOrderBy = "quality_score"

	// The fields below are accepted by the v2 endpoint (GetNewsV2 / GetNewsV2Page).
	NewsOrderByPublishedAt    NewsOrderBy = "publishedAt"
	NewsOrderByIngestedAt     NewsOrderBy = "ingestedAt"
	NewsOrderByUpdatedAt      NewsOrderBy = "updatedAt"
	NewsOrderByQualityScoreV2 NewsOrderBy = "qualityScore"
)

// NewsOrder is one field of a multi-field v2 sort. An empty Direction lets the API apply its
// default (desc).
type NewsOrder struct {
	Field     NewsOrderBy
	Direction SortDirection
}

// NewsTickerRelevance says how central a tagged company is to an article.
type NewsTickerRelevance string

const (
	// NewsTickerRelevancePrimary: the article is about this company (at most one per article).
	NewsTickerRelevancePrimary NewsTickerRelevance = "primary"
	// NewsTickerRelevanceMentioned: the company is named with some substance.
	NewsTickerRelevanceMentioned NewsTickerRelevance = "mentioned"
	// NewsTickerRelevanceListed: the company only appears in a list, ranking or market table.
	NewsTickerRelevanceListed NewsTickerRelevance = "listed"
)

// NewsEventType is an event type slug, accepted by the EventTypes / ExcludeEventTypes filters.
// GetNewsEventTypes returns the current list with localized names; slugs not listed here can be
// passed as NewsEventType("…").
type NewsEventType string

const (
	NewsEventTypeFinancialResults     NewsEventType = "financial-results"
	NewsEventTypeGuidance             NewsEventType = "guidance"
	NewsEventTypeOperationalData      NewsEventType = "operational-data"
	NewsEventTypeOrdersInvestments    NewsEventType = "orders-investments"
	NewsEventTypeManagementChange     NewsEventType = "management-change"
	NewsEventTypeRegulation           NewsEventType = "regulation"
	NewsEventTypeLitigation           NewsEventType = "litigation"
	NewsEventTypeCreditRating         NewsEventType = "credit-rating"
	NewsEventTypeDividend             NewsEventType = "dividend"
	NewsEventTypeAnalystReport        NewsEventType = "analyst-report"
	NewsEventTypePriceActionTechnical NewsEventType = "price-action-technical"
	NewsEventTypeOther                NewsEventType = "other"
)

// NewsOriginality marks whether an article is the first report of its story or later coverage.
type NewsOriginality string

const (
	NewsOriginalityOriginal NewsOriginality = "original"
	NewsOriginalityRewrite  NewsOriginality = "rewrite"
)

// NewsLane is a curated news surface that groups upstream sources. It is accepted by the `lane`
// filter of the News, News V2, and Live News Stream endpoints. US lanes are NewsLaneGlobalMacro
// and NewsLaneFastMovers; TR lanes are NewsLaneTrEkonomi and NewsLaneBist.
type NewsLane string

const (
	NewsLaneGlobalMacro NewsLane = "global_macro"
	NewsLaneTrEkonomi   NewsLane = "tr_ekonomi"
	NewsLaneBist        NewsLane = "bist"
	NewsLaneFastMovers  NewsLane = "fast_movers"
)

// NewsHighlight is a single dated highlights record, with the top news items bucketed by sector
// group. It is one item of the paginated News Highlights response.
type NewsHighlight struct {
	ID                      string    `json:"id"`
	CreatedAt               time.Time `json:"createdAt"`
	Consumer                []string  `json:"consumer"`
	EnergyAndUtilities      []string  `json:"energyAndUtilities"`
	Finance                 []string  `json:"finance"`
	Healthcare              []string  `json:"healthcare"`
	IndustrialsAndMaterials []string  `json:"industrialsAndMaterials"`
	Tech                    []string  `json:"tech"`
	Other                   []string  `json:"other"`
}

type NewsPublisher struct {
	Name    string  `json:"name"`
	LogoUrl *string `json:"logoUrl"`
}

type NewsTicker struct {
	// ID is the Laplace stock ID. Empty if the symbol has no listing.
	ID     string `json:"id"`
	Name   string `json:"name"`
	Symbol string `json:"symbol,omitempty"`
	// Region is the listing market. Set on v2 and stream items.
	Region Region `json:"region,omitempty"`
	// Relevance is empty when the article hasn't been classified. Set on v2 and stream items.
	Relevance NewsTickerRelevance `json:"relevance,omitempty"`
}

// NewsEventTypeInfo is an event type with its localized name, as attached to a news item and as
// returned by GetNewsEventTypes.
type NewsEventTypeInfo struct {
	Slug     NewsEventType `json:"slug"`
	Name     string        `json:"name"`
	Language Locale        `json:"language"`
}

type NewsCategories struct {
	ID           string  `json:"id"`
	Name         string  `json:"name"`
	CategoryType *string `json:"categoryType,omitempty"`
}

type NewsSector struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type NewsContent struct {
	Title           string   `json:"title"`
	Description     string   `json:"description"`
	Content         []string `json:"content"`
	Summary         []string `json:"summary"`
	InvestorInsight string   `json:"investorInsight"`
}

type NewsIndustry struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type News struct {
	ID             string          `json:"id"`
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

// NewsV2 is a news item from the v2 endpoint and the live news stream.
//
// Fields added after the first v2 release are left at their zero value when a response doesn't
// carry them.
type NewsV2 struct {
	// ID is stable: it stays the same when the article is updated or replaced by a higher-trust copy.
	ID       string `json:"id"`
	URL      string `json:"url"`
	ImageUrl string `json:"imageUrl"`
	// Deprecated: Timestamp is PublishedAt if known, else IngestedAt. Use PublishedAt.
	Timestamp    time.Time     `json:"timestamp"`
	PublisherUrl string        `json:"publisherUrl"`
	Publisher    NewsPublisher `json:"publisher"`
	// QualityScore is 0–10; 0 means not scored yet.
	QualityScore int64 `json:"qualityScore"`
	// CreatedAt is the same as IngestedAt.
	CreatedAt  time.Time       `json:"createdAt"`
	Tickers    []NewsTicker    `json:"tickers,omitempty"`
	Categories *NewsCategories `json:"categories,omitempty"`
	Sectors    *NewsSector     `json:"sectors,omitempty"`
	Content    *NewsContent    `json:"content,omitempty"`
	Industries *NewsIndustry   `json:"industries,omitempty"`

	// PublishedAt is when the source published the article; nil when the source gave no usable date.
	PublishedAt *time.Time `json:"publishedAt,omitempty"`
	// IngestedAt is when Laplace first stored the article.
	IngestedAt time.Time `json:"ingestedAt"`
	// ContentUpdatedAt is the last time the title or text changed; nil if it never changed.
	ContentUpdatedAt *time.Time `json:"contentUpdatedAt,omitempty"`
	// UpdatedAt is the last change of any kind (content, summary, translation, tags, score or
	// cluster). Use it with GetNewsV2Params.UpdatedSince.
	UpdatedAt time.Time `json:"updatedAt"`
	// PrimaryTicker is the one company the article is about; empty if none.
	PrimaryTicker string `json:"primaryTicker,omitempty"`
	// IsListArticle is true for market wraps, rankings and "most traded" lists.
	IsListArticle bool `json:"isListArticle"`
	// EventType is nil until the article is classified.
	EventType *NewsEventTypeInfo `json:"eventType,omitempty"`
	// ClusterID identifies the story. An article with no related coverage is its own cluster.
	ClusterID     string          `json:"clusterId,omitempty"`
	IsClusterLead bool            `json:"isClusterLead"`
	ClusterSize   int             `json:"clusterSize"`
	Originality   NewsOriginality `json:"originality,omitempty"`
}

// NewsCategory is a canonical news category with a localized name, as returned by the
// News Categories endpoint. The ID value is what the `categoryIds` filter of the News, News V2
// and Live News Stream endpoints accepts.
type NewsCategory struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// GetNewsCategories returns the full canonical news category list with localized names for the
// given locale. All categories are always returned regardless of whether they currently have
// tagged news, making the result suitable for populating dropdowns and filter UIs.
func (c *Client) GetNewsCategories(ctx context.Context, locale Locale) ([]NewsCategory, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/api/v1/news/categories", c.baseUrl), nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Add("locale", string(locale))
	req.URL.RawQuery = q.Encode()

	resp, err := sendRequest[[]NewsCategory](ctx, c, req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// GetNewsEventTypes returns the event type slugs accepted by the EventTypes and ExcludeEventTypes
// filters, with names in the given locale. An empty locale uses the API default (en).
func (c *Client) GetNewsEventTypes(ctx context.Context, locale Locale) ([]NewsEventTypeInfo, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/api/v1/news/event-types", c.baseUrl), nil)
	if err != nil {
		return nil, err
	}

	if locale != "" {
		q := req.URL.Query()
		q.Add("locale", string(locale))
		req.URL.RawQuery = q.Encode()
	}

	return sendRequest[[]NewsEventTypeInfo](ctx, c, req)
}

// NewsLaneInfo is a selectable news lane as returned by the News Lanes endpoint. The ID value is
// the exact value accepted by the `lane` filter of the News, News V2, and Live News Stream endpoints.
type NewsLaneInfo struct {
	ID    NewsLane `json:"id"`
	Label string   `json:"label"`
}

// GetNewsLanesParams holds the optional parameters for the news lanes endpoint.
type GetNewsLanesParams struct {
	// Region, when set (RegionTr or RegionUs), restricts the result to lanes valid for that region.
	Region Region
}

// GetNewsLanes returns the fixed list of news lanes (id + label) suitable for building a lane
// filter UI. The returned IDs feed the `lane` parameter of the News and Live News Stream endpoints.
func (c *Client) GetNewsLanes(ctx context.Context, params GetNewsLanesParams) ([]NewsLaneInfo, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/api/v1/news/lanes", c.baseUrl), nil)
	if err != nil {
		return nil, err
	}

	if params.Region != "" {
		q := req.URL.Query()
		q.Add("region", string(params.Region))
		req.URL.RawQuery = q.Encode()
	}

	resp, err := sendRequest[[]NewsLaneInfo](ctx, c, req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// NewsApiSource is a configured upstream news source as returned by the News API Source Names
// endpoint. The ID value is accepted by the `apiSource` filter of the News and Live News Stream
// endpoints (comma-separated for multiple); Name is the human-readable display name
// (e.g. "BBC Business", "Gazete Oksijen").
type NewsApiSource struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// GetNewsApiSourceNamesParams holds the optional parameters for the news api-source-names endpoint.
type GetNewsApiSourceNamesParams struct {
	Region   Region
	Language Locale
}

// GetNewsApiSourceNames returns the configured news sources (id + name), suitable for populating a
// source filter. The returned IDs feed the `apiSource` parameter of the News and Live News Stream
// endpoints.
func (c *Client) GetNewsApiSourceNames(ctx context.Context, params GetNewsApiSourceNamesParams) ([]NewsApiSource, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/api/v1/news/api-source-names", c.baseUrl), nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	if params.Region != "" {
		q.Add("region", string(params.Region))
	}
	if params.Language != "" {
		q.Add("language", string(params.Language))
	}
	req.URL.RawQuery = q.Encode()

	resp, err := sendRequest[[]NewsApiSource](ctx, c, req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// GetNewsHighlightsParams holds the parameters for the news highlights endpoint. Region (which must
// be RegionUs) and Locale are required. From/To (YYYY-MM-DD) narrow the result to highlights created
// in that window; Skip/Top page the result, where Top is the page size (1-20).
type GetNewsHighlightsParams struct {
	Region Region
	Locale Locale
	From   string
	To     string
	Skip   *int
	Top    *int
}

// GetNewsHighlights retrieves paginated, sector-categorized US news highlights, newest first.
func (c *Client) GetNewsHighlights(ctx context.Context, params GetNewsHighlightsParams) (*PaginatedResponse[NewsHighlight], error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/api/v1/news/highlights", c.baseUrl), nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Add("region", string(params.Region))
	q.Add("locale", string(params.Locale))
	if params.From != "" {
		q.Add("from", params.From)
	}
	if params.To != "" {
		q.Add("to", params.To)
	}
	if params.Skip != nil {
		q.Add("skip", strconv.Itoa(*params.Skip))
	}
	if params.Top != nil {
		q.Add("top", strconv.Itoa(*params.Top))
	}
	req.URL.RawQuery = q.Encode()

	resp, err := sendRequest[PaginatedResponse[NewsHighlight]](ctx, c, req)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

type GetNewsParams struct {
	Region           Region
	Locale           Locale
	NewsType         NewsType
	Lane             NewsLane
	Page             *int
	Size             *int
	OrderBy          NewsOrderBy
	OrderByDirection SortDirection
	Symbols          string
	CategoryIds      string
	SectorIds        string
	IndustryIds      string
	ApiSource        string
	QualityScoreMin  *int
	QualityScoreMax  *int
	// Deprecated: on v2 use GetNewsV2Params.PublishedFrom.
	TimestampFrom string
	// Deprecated: on v2 use GetNewsV2Params.PublishedTo.
	TimestampTo string
}

// buildNewsQuery builds the shared query parameters for the news endpoints.
func buildNewsQuery(params GetNewsParams) url.Values {
	q := url.Values{}
	q.Add("region", string(params.Region))
	q.Add("locale", string(params.Locale))
	if params.NewsType != "" {
		q.Add("newsType", string(params.NewsType))
	}
	if params.Lane != "" {
		q.Add("lane", string(params.Lane))
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
	if params.Symbols != "" {
		q.Add("symbols", params.Symbols)
	}
	if params.CategoryIds != "" {
		q.Add("categoryIds", params.CategoryIds)
	}
	if params.SectorIds != "" {
		q.Add("sectorIds", params.SectorIds)
	}
	if params.IndustryIds != "" {
		q.Add("industryIds", params.IndustryIds)
	}
	if params.ApiSource != "" {
		q.Add("apiSource", params.ApiSource)
	}
	if params.QualityScoreMin != nil {
		q.Add("qualityScoreMin", strconv.Itoa(*params.QualityScoreMin))
	}
	if params.QualityScoreMax != nil {
		q.Add("qualityScoreMax", strconv.Itoa(*params.QualityScoreMax))
	}
	if params.TimestampFrom != "" {
		q.Add("timestampFrom", params.TimestampFrom)
	}
	if params.TimestampTo != "" {
		q.Add("timestampTo", params.TimestampTo)
	}
	return q
}

// GetNews retrieves a paginated list of news articles with optional filtering, ordering, and pagination.
func (c *Client) GetNews(ctx context.Context, params GetNewsParams) (*PaginatedResponse[News], error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/api/v1/news", c.baseUrl), nil)
	if err != nil {
		return nil, err
	}

	req.URL.RawQuery = buildNewsQuery(params).Encode()

	resp, err := sendRequest[PaginatedResponse[News]](ctx, c, req)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

// GetNewsV2 retrieves a paginated list of news articles from the v2 endpoint, excluding related tickers.
//
// It only sends the filters shared with v1. GetNewsV2Page takes the v2-only filters and cursor
// paging.
func (c *Client) GetNewsV2(ctx context.Context, params GetNewsParams) (*PaginatedResponse[NewsV2], error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/api/v2/news", c.baseUrl), nil)
	if err != nil {
		return nil, err
	}

	req.URL.RawQuery = buildNewsQuery(params).Encode()

	resp, err := sendRequest[PaginatedResponse[NewsV2]](ctx, c, req)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

// GetNewsV2Params holds the parameters for GetNewsV2Page and ForEachNewsV2. It embeds the filters
// shared with v1 and adds the ones only v2 accepts. Unset fields aren't sent.
//
// Size accepts any value from 1 to 100 on v2.
type GetNewsV2Params struct {
	GetNewsParams

	// TickerRelevance keeps news where a ticker (one of Symbols, if given) has one of these
	// relevances.
	TickerRelevance []NewsTickerRelevance
	// ExcludeListArticles, when true, drops market wraps, "most traded" lists and top-gainer tables.
	ExcludeListArticles *bool
	EventTypes          []NewsEventType
	// ExcludeEventTypes drops these event types. Articles not classified yet are kept.
	ExcludeEventTypes []NewsEventType
	// CollapseClusters, when true, returns each story once instead of once per outlet.
	CollapseClusters *bool
	// MaxPublisherShare (1–100, percent) limits how much of the result one publisher can take.
	// The page's PublisherCap reports the share actually achieved.
	MaxPublisherShare *int

	// Dates take RFC3339 or YYYY-MM-DD (UTC).
	PublishedFrom string
	// PublishedTo given as a date includes that whole day.
	PublishedTo string
	// UpdatedSince filters on UpdatedAt and changes the default order to updatedAt ascending.
	UpdatedSince string
	// AsOf returns the result as it was at that time.
	AsOf string

	// OrderByFields sorts by up to 3 fields. When set, it replaces OrderBy / OrderByDirection.
	// Cursor paging only works with a single field.
	OrderByFields []NewsOrder

	// Cursor is the NextCursor of the previous page. It can't be combined with Page.
	Cursor string
	// IncludeTotal set to false skips counting all matches, which makes requests faster.
	// RecordCount is then not set; see NewsV2Page.HasRecordCount.
	IncludeTotal *bool
}

// NewsPublisherCap reports how the MaxPublisherShare filter was applied.
type NewsPublisherCap struct {
	// Share is the requested percentage.
	Share float64 `json:"share"`
	// PerPublisher is the item limit applied to each publisher.
	PerPublisher int `json:"perPublisher"`
	// AchievedShare is the largest share any publisher actually has. It can be above Share when
	// few publishers cover the query.
	AchievedShare float64 `json:"achievedShare"`
}

// NewsV2Page is a page of the v2 news endpoint. It embeds PaginatedResponse, so RecordCount and
// Items read the same as on GetNewsV2's result.
type NewsV2Page struct {
	PaginatedResponse[NewsV2]

	// HasRecordCount is false when the API returned no total (IncludeTotal=false). RecordCount is
	// then 0.
	HasRecordCount bool `json:"-"`
	// NextCursor is empty on the last page. Pass it back unchanged as GetNewsV2Params.Cursor.
	NextCursor string `json:"nextCursor,omitempty"`
	// PublisherCap is set only when MaxPublisherShare was sent.
	PublisherCap *NewsPublisherCap `json:"publisherCap,omitempty"`
}

func (p *NewsV2Page) UnmarshalJSON(data []byte) error {
	var raw struct {
		RecordCount  *int              `json:"recordCount"`
		Items        []NewsV2          `json:"items"`
		NextCursor   string            `json:"nextCursor"`
		PublisherCap *NewsPublisherCap `json:"publisherCap"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	*p = NewsV2Page{
		PaginatedResponse: PaginatedResponse[NewsV2]{Items: raw.Items},
		HasRecordCount:    raw.RecordCount != nil,
		NextCursor:        raw.NextCursor,
		PublisherCap:      raw.PublisherCap,
	}
	if raw.RecordCount != nil {
		p.RecordCount = *raw.RecordCount
	}
	return nil
}

func joinStrings[T ~string](values []T) string {
	parts := make([]string, len(values))
	for i, v := range values {
		parts[i] = string(v)
	}
	return strings.Join(parts, ",")
}

func renderNewsOrder(fields []NewsOrder) string {
	parts := make([]string, len(fields))
	for i, f := range fields {
		parts[i] = string(f.Field)
		if f.Direction != "" {
			parts[i] += " " + string(f.Direction)
		}
	}
	return strings.Join(parts, ",")
}

func buildNewsV2Query(params GetNewsV2Params) url.Values {
	q := buildNewsQuery(params.GetNewsParams)
	if len(params.TickerRelevance) > 0 {
		q.Set("tickerRelevance", joinStrings(params.TickerRelevance))
	}
	if params.ExcludeListArticles != nil {
		q.Set("excludeListArticles", strconv.FormatBool(*params.ExcludeListArticles))
	}
	if len(params.EventTypes) > 0 {
		q.Set("eventTypes", joinStrings(params.EventTypes))
	}
	if len(params.ExcludeEventTypes) > 0 {
		q.Set("excludeEventTypes", joinStrings(params.ExcludeEventTypes))
	}
	if params.CollapseClusters != nil {
		q.Set("collapseClusters", strconv.FormatBool(*params.CollapseClusters))
	}
	if params.MaxPublisherShare != nil {
		q.Set("maxPublisherShare", strconv.Itoa(*params.MaxPublisherShare))
	}
	if params.PublishedFrom != "" {
		q.Set("publishedFrom", params.PublishedFrom)
	}
	if params.PublishedTo != "" {
		q.Set("publishedTo", params.PublishedTo)
	}
	if params.UpdatedSince != "" {
		q.Set("updatedSince", params.UpdatedSince)
	}
	if params.AsOf != "" {
		q.Set("asOf", params.AsOf)
	}
	if len(params.OrderByFields) > 0 {
		q.Set("orderBy", renderNewsOrder(params.OrderByFields))
		q.Del("orderByDirection")
	}
	if params.Cursor != "" {
		q.Set("cursor", params.Cursor)
	}
	if params.IncludeTotal != nil {
		q.Set("includeTotal", strconv.FormatBool(*params.IncludeTotal))
	}
	return q
}

// GetNewsV2Page retrieves one page of news from the v2 endpoint with every v2 filter, ordering
// and cursor paging. Newest first by default.
//
// To read every page, use ForEachNewsV2, or pass the returned NextCursor back as params.Cursor
// until it is empty.
func (c *Client) GetNewsV2Page(ctx context.Context, params GetNewsV2Params) (*NewsV2Page, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/api/v2/news", c.baseUrl), nil)
	if err != nil {
		return nil, err
	}

	req.URL.RawQuery = buildNewsV2Query(params).Encode()

	resp, err := sendRequest[NewsV2Page](ctx, c, req)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

// ForEachNewsV2 follows cursors from params.Cursor (the first page when empty) until the last
// page, calling fn for each item in order. It sends includeTotal=false and keeps every other
// parameter the same between pages. params.Page must be nil, since the API rejects it together
// with a cursor.
//
// It stops at the first error from the API or from fn and returns it.
func (c *Client) ForEachNewsV2(ctx context.Context, params GetNewsV2Params, fn func(NewsV2) error) error {
	if params.Page != nil {
		return errors.New("laplace: ForEachNewsV2 uses cursor paging; leave Page unset")
	}
	includeTotal := false
	params.IncludeTotal = &includeTotal

	for {
		page, err := c.GetNewsV2Page(ctx, params)
		if err != nil {
			return err
		}
		for _, item := range page.Items {
			if err := fn(item); err != nil {
				return err
			}
		}
		if page.NextCursor == "" {
			return nil
		}
		params.Cursor = page.NextCursor
	}
}

// NewsStreamResult is the result type for news streams
type NewsStreamResult struct {
	Data  []NewsV2
	Error error
}

func sendNewsSSERequest(
	ctx context.Context,
	c *Client,
	url string,
) (<-chan NewsStreamResult, func(), error) {

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, nil, err
	}

	// Set headers
	req.Header.Set("Accept", "text/event-stream")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	// Send the request
	resp, err := c.cli.Do(req.WithContext(ctx))
	if err != nil {
		return nil, nil, err
	}

	// Check the response status
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		defer resp.Body.Close()
		return nil, nil, fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(bodyBytes))
	}

	// Create a single channel for results
	results := make(chan NewsStreamResult)

	ctxWithCancel, cancel := context.WithCancel(ctx)

	// Start a goroutine to read the SSE stream
	go func() {
		defer resp.Body.Close()
		defer close(results)
		defer cancel()

		scanner := bufio.NewScanner(resp.Body)
		for scanner.Scan() {
			select {
			case <-ctxWithCancel.Done():
				return
			default:
				line := scanner.Text()
				if strings.HasPrefix(line, "data:") {
					data := strings.TrimPrefix(line, "data:")
					var event []NewsV2
					if err := json.Unmarshal([]byte(data), &event); err != nil {
						results <- NewsStreamResult{Error: fmt.Errorf("error unmarshalling event: %w", err)}
						continue
					}
					results <- NewsStreamResult{Data: event}
				}
			}
		}

		if err := scanner.Err(); err != nil {
			results <- NewsStreamResult{Error: fmt.Errorf("error reading SSE stream: %w", err)}
		}
	}()

	return results, cancel, nil
}

// StreamNewsParams holds the parameters for the news stream endpoint. The filters mean the same
// as on GetNewsV2Params; unset ones aren't sent.
type StreamNewsParams struct {
	Region      Region
	Locale      Locale
	Lane        NewsLane
	Symbols     []string
	CategoryIds []string
	SectorIds   []string
	IndustryIds []string
	ApiSource   []string

	TickerRelevance     []NewsTickerRelevance
	ExcludeListArticles *bool
	EventTypes          []NewsEventType
	ExcludeEventTypes   []NewsEventType
	CollapseClusters    *bool
}

// NewsStream handles live news streaming for a specific locale and filters
type NewsStream struct {
	mu           sync.RWMutex
	ctx          context.Context
	cancel       context.CancelFunc
	sseChan      <-chan NewsStreamResult
	outputChan   chan NewsStreamResult
	c            *Client
	params       StreamNewsParams
	closed       bool
	isSubscribed bool
}

// Subscribe starts receiving news from the stream
func (s *NewsStream) Subscribe(ctx context.Context) error {
	if ctx == nil {
		return fmt.Errorf("context cannot be nil")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	// Cleanup existing stream
	if err := s.cleanupExistingStream(); err != nil {
		return fmt.Errorf("failed to cleanup existing stream: %w", err)
	}

	s.outputChan = make(chan NewsStreamResult, 100) // Buffered channel
	s.closed = false
	s.ctx = ctx

	// Start streaming
	if err := s.startStreaming(); err != nil {
		return fmt.Errorf("failed to start streaming: %w", err)
	}

	s.isSubscribed = true
	return nil
}

// Receive returns a channel to receive news data
func (s *NewsStream) Receive() <-chan NewsStreamResult {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if !s.isSubscribed {
		// Return a closed channel if not subscribed
		ch := make(chan NewsStreamResult)
		close(ch)
		return ch
	}

	return s.outputChan
}

// Close closes the stream and cleanup resources
func (s *NewsStream) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.closed {
		return nil
	}

	s.closed = true
	s.isSubscribed = false
	return s.cleanupExistingStream()
}

// cleanupExistingStream cancels and cleans up existing streaming task
func (s *NewsStream) cleanupExistingStream() error {
	if s.cancel != nil {
		s.cancel()
		s.cancel = nil
	}

	if s.outputChan != nil {
		close(s.outputChan)
		s.outputChan = nil
	}

	return nil
}

// startStreaming starts the SSE streaming connection
func (s *NewsStream) startStreaming() error {
	reqURL, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/api/v1/news/stream", s.c.baseUrl), nil)
	if err != nil {
		return fmt.Errorf("failed to create SSE request URL: %w", err)
	}

	q := reqURL.URL.Query()
	q.Add("locale", string(s.params.Locale))
	q.Add("region", string(s.params.Region))
	if s.params.Lane != "" {
		q.Add("lane", string(s.params.Lane))
	}
	if len(s.params.Symbols) > 0 {
		q.Add("symbols", strings.Join(s.params.Symbols, ","))
	}
	if len(s.params.CategoryIds) > 0 {
		q.Add("categoryIds", strings.Join(s.params.CategoryIds, ","))
	}
	if len(s.params.SectorIds) > 0 {
		q.Add("sectorIds", strings.Join(s.params.SectorIds, ","))
	}
	if len(s.params.IndustryIds) > 0 {
		q.Add("industryIds", strings.Join(s.params.IndustryIds, ","))
	}
	if len(s.params.ApiSource) > 0 {
		q.Add("apiSource", strings.Join(s.params.ApiSource, ","))
	}
	if len(s.params.TickerRelevance) > 0 {
		q.Add("tickerRelevance", joinStrings(s.params.TickerRelevance))
	}
	if s.params.ExcludeListArticles != nil {
		q.Add("excludeListArticles", strconv.FormatBool(*s.params.ExcludeListArticles))
	}
	if len(s.params.EventTypes) > 0 {
		q.Add("eventTypes", joinStrings(s.params.EventTypes))
	}
	if len(s.params.ExcludeEventTypes) > 0 {
		q.Add("excludeEventTypes", joinStrings(s.params.ExcludeEventTypes))
	}
	if s.params.CollapseClusters != nil {
		q.Add("collapseClusters", strconv.FormatBool(*s.params.CollapseClusters))
	}
	reqURL.URL.RawQuery = q.Encode()

	ctxWithCancel, cancel := context.WithCancel(s.ctx)
	s.cancel = cancel

	channel, _, err := sendNewsSSERequest(ctxWithCancel, s.c, reqURL.URL.String())
	if err != nil {
		return fmt.Errorf("failed to establish SSE connection: %w", err)
	}

	s.sseChan = channel
	go s.forwardData()

	return nil
}

// forwardData forwards data from SSE channel to output channel
func (s *NewsStream) forwardData() {
	defer func() {
		if r := recover(); r != nil {
			s.c.logger.Error("panic in news stream forwardData", r)
		}
	}()

	for {
		select {
		case data, ok := <-s.sseChan:
			if !ok {
				return
			}

			s.mu.RLock()
			outputChan := s.outputChan
			closed := s.closed
			s.mu.RUnlock()

			if closed || outputChan == nil {
				return
			}

			select {
			case outputChan <- data:
			case <-s.ctx.Done():
				return
			}
		case <-s.ctx.Done():
			return
		}
	}
}

// GetNewsStream creates a new news stream.
// Call Subscribe(ctx) on the returned stream to start receiving data.
func (c *Client) GetNewsStream(params StreamNewsParams) *NewsStream {
	return &NewsStream{
		c:      c,
		params: params,
		closed: false,
	}
}

// CreateNewsStream creates and subscribes to a news stream.
func (c *Client) CreateNewsStream(ctx context.Context, params StreamNewsParams) (*NewsStream, error) {
	stream := c.GetNewsStream(params)
	if err := stream.Subscribe(ctx); err != nil {
		return nil, fmt.Errorf("failed to subscribe to news stream: %w", err)
	}
	return stream, nil
}
