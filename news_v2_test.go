package laplace

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"
)

// Unit tests for the v2 news additions. They run against httptest servers and need no API key.

func ptr[T any](v T) *T { return &v }

func TestBuildNewsV2QuerySendsNewParams(t *testing.T) {
	q := buildNewsV2Query(GetNewsV2Params{
		GetNewsParams: GetNewsParams{
			Region:  RegionTr,
			Locale:  LocaleTr,
			Symbols: "THYAO,GARAN",
			Size:    ptr(73),
		},
		TickerRelevance:     []NewsTickerRelevance{NewsTickerRelevancePrimary, NewsTickerRelevanceMentioned},
		ExcludeListArticles: ptr(false),
		EventTypes:          []NewsEventType{NewsEventTypeDividend, "some-new-slug"},
		ExcludeEventTypes:   []NewsEventType{NewsEventTypeOther},
		CollapseClusters:    ptr(true),
		MaxPublisherShare:   ptr(25),
		PublishedFrom:       "2026-11-01",
		PublishedTo:         "2026-11-12T09:00:00Z",
		UpdatedSince:        "2026-11-10T00:00:00Z",
		AsOf:                "2026-11-11T00:00:00Z",
		Cursor:              "abc",
		IncludeTotal:        ptr(false),
	})

	want := map[string]string{
		"region":              "tr",
		"locale":              "tr",
		"symbols":             "THYAO,GARAN",
		"size":                "73",
		"tickerRelevance":     "primary,mentioned",
		"excludeListArticles": "false",
		"eventTypes":          "dividend,some-new-slug",
		"excludeEventTypes":   "other",
		"collapseClusters":    "true",
		"maxPublisherShare":   "25",
		"publishedFrom":       "2026-11-01",
		"publishedTo":         "2026-11-12T09:00:00Z",
		"updatedSince":        "2026-11-10T00:00:00Z",
		"asOf":                "2026-11-11T00:00:00Z",
		"cursor":              "abc",
		"includeTotal":        "false",
	}
	for k, v := range want {
		if got := q.Get(k); got != v {
			t.Errorf("%s = %q, want %q", k, got, v)
		}
	}
	if len(q) != len(want) {
		t.Errorf("unexpected extra params: %v", q)
	}
}

func TestBuildNewsV2QueryOmitsUnsetParams(t *testing.T) {
	q := buildNewsV2Query(GetNewsV2Params{GetNewsParams: GetNewsParams{Region: RegionUs, Locale: LocaleEn}})
	if len(q) != 2 {
		t.Fatalf("want only region and locale, got %v", q)
	}
}

func TestNewsV1QueryUnchanged(t *testing.T) {
	q := buildNewsQuery(GetNewsParams{Region: RegionUs, Locale: LocaleEn, OrderBy: NewsOrderByQualityScore})
	if got := q.Encode(); got != "locale=en&orderBy=quality_score&region=us" {
		t.Errorf("v1 query = %s", got)
	}
}

func TestNewsOrderByRendering(t *testing.T) {
	base := GetNewsParams{Region: RegionUs, Locale: LocaleEn}
	cases := []struct {
		name          string
		params        GetNewsV2Params
		wantOrderBy   string
		wantDirection string
	}{
		{"single field", GetNewsV2Params{GetNewsParams: withOrder(base, NewsOrderByPublishedAt, "")}, "publishedAt", ""},
		{"single field with direction", GetNewsV2Params{GetNewsParams: withOrder(base, NewsOrderByUpdatedAt, SortDirectionAsc)}, "updatedAt", "asc"},
		{"legacy quality_score", GetNewsV2Params{GetNewsParams: withOrder(base, NewsOrderByQualityScore, "")}, "quality_score", ""},
		{"multi field", GetNewsV2Params{
			GetNewsParams: withOrder(base, NewsOrderByTimestamp, SortDirectionAsc),
			OrderByFields: []NewsOrder{
				{NewsOrderByQualityScoreV2, SortDirectionDesc},
				{NewsOrderByPublishedAt, SortDirectionDesc},
				{Field: NewsOrderByIngestedAt},
			},
		}, "qualityScore desc,publishedAt desc,ingestedAt", ""},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			q := buildNewsV2Query(tc.params)
			if got := q.Get("orderBy"); got != tc.wantOrderBy {
				t.Errorf("orderBy = %q, want %q", got, tc.wantOrderBy)
			}
			if got := q.Get("orderByDirection"); got != tc.wantDirection {
				t.Errorf("orderByDirection = %q, want %q", got, tc.wantDirection)
			}
		})
	}
}

func withOrder(p GetNewsParams, field NewsOrderBy, dir SortDirection) GetNewsParams {
	p.OrderBy = field
	p.OrderByDirection = dir
	return p
}

const newsItemAllFields = `{
	"id": "3f6c1e2a",
	"url": "https://www.ekonomim.com/a",
	"imageUrl": "https://x/image.jpg",
	"publisherUrl": "https://www.ekonomim.com",
	"publisher": {"name": "Ekonomim", "logoUrl": null},
	"publishedAt": null,
	"ingestedAt": "2026-11-12T10:24:13.106706Z",
	"contentUpdatedAt": null,
	"updatedAt": "2026-11-12T10:31:02.55192Z",
	"createdAt": "2026-11-12T10:24:13.106706Z",
	"timestamp": "2026-11-12T10:24:13.106706Z",
	"content": {"title": "THY ekim ayında 8,1 milyon yolcu taşıdı", "description": "d", "content": ["p"], "summary": ["s"], "investorInsight": "i"},
	"primaryTicker": "THYAO",
	"tickers": [{"id": "61dd0d6", "symbol": "THYAO", "name": "Türk Hava Yolları", "region": "tr", "relevance": "primary"}],
	"isListArticle": false,
	"eventType": {"slug": "operational-data", "name": "Operasyonel Veriler", "language": "tr"},
	"qualityScore": 7,
	"clusterId": "9b1d",
	"isClusterLead": true,
	"clusterSize": 4,
	"originality": "original",
	"someFieldAddedLater": {"x": 1}
}`

const newsItemOldShape = `{
	"id": "old",
	"url": "https://x/a",
	"imageUrl": "",
	"publisherUrl": "https://x",
	"publisher": {"name": "X"},
	"timestamp": "2025-01-02T03:04:05Z",
	"createdAt": "2025-01-02T03:04:05Z",
	"qualityScore": 5,
	"tickers": [{"id": "1", "name": "Apple", "symbol": "AAPL"}]
}`

func TestNewsV2ItemParsesAllNewFields(t *testing.T) {
	var n NewsV2
	if err := json.Unmarshal([]byte(newsItemAllFields), &n); err != nil {
		t.Fatal(err)
	}
	if n.PublishedAt != nil || n.ContentUpdatedAt != nil {
		t.Errorf("null times should stay nil: %v %v", n.PublishedAt, n.ContentUpdatedAt)
	}
	if n.IngestedAt.IsZero() || n.UpdatedAt.IsZero() {
		t.Error("ingestedAt / updatedAt not parsed")
	}
	if n.PrimaryTicker != "THYAO" || n.IsListArticle || n.ClusterID != "9b1d" || !n.IsClusterLead ||
		n.ClusterSize != 4 || n.Originality != NewsOriginalityOriginal {
		t.Errorf("unexpected item fields: %+v", n)
	}
	if n.EventType == nil || n.EventType.Slug != NewsEventTypeOperationalData || n.EventType.Name != "Operasyonel Veriler" || n.EventType.Language != LocaleTr {
		t.Errorf("eventType = %+v", n.EventType)
	}
	tk := n.Tickers[0]
	if tk.Region != RegionTr || tk.Relevance != NewsTickerRelevancePrimary || tk.Name != "Türk Hava Yolları" {
		t.Errorf("ticker = %+v", tk)
	}
}

func TestNewsV2ItemParsesOldShape(t *testing.T) {
	var n NewsV2
	if err := json.Unmarshal([]byte(newsItemOldShape), &n); err != nil {
		t.Fatal(err)
	}
	if n.ID != "old" || n.Timestamp.IsZero() || n.PublishedAt != nil || n.EventType != nil || n.ClusterSize != 0 {
		t.Errorf("unexpected item: %+v", n)
	}
	if n.Tickers[0].Region != "" || n.Tickers[0].Relevance != "" {
		t.Errorf("ticker = %+v", n.Tickers[0])
	}
}

func TestNewsV2PageParsing(t *testing.T) {
	t.Run("null recordCount with cursor and publisherCap", func(t *testing.T) {
		var p NewsV2Page
		body := `{"recordCount": null, "items": [` + newsItemOldShape + `], "nextCursor": "c1",
			"publisherCap": {"share": 25, "perPublisher": 42, "achievedShare": 24.9}}`
		if err := json.Unmarshal([]byte(body), &p); err != nil {
			t.Fatal(err)
		}
		if p.HasRecordCount || p.RecordCount != 0 {
			t.Errorf("recordCount: has=%v value=%d", p.HasRecordCount, p.RecordCount)
		}
		if len(p.Items) != 1 || p.NextCursor != "c1" {
			t.Errorf("page = %+v", p)
		}
		if p.PublisherCap == nil || *p.PublisherCap != (NewsPublisherCap{Share: 25, PerPublisher: 42, AchievedShare: 24.9}) {
			t.Errorf("publisherCap = %+v", p.PublisherCap)
		}
	})

	t.Run("old page shape", func(t *testing.T) {
		var p NewsV2Page
		if err := json.Unmarshal([]byte(`{"recordCount": 549, "items": []}`), &p); err != nil {
			t.Fatal(err)
		}
		if !p.HasRecordCount || p.RecordCount != 549 || p.NextCursor != "" || p.PublisherCap != nil {
			t.Errorf("page = %+v", p)
		}
		// The embedded PaginatedResponse keeps the old type reachable.
		var _ PaginatedResponse[NewsV2] = p.PaginatedResponse
	})
}

func TestForEachNewsV2FollowsCursor(t *testing.T) {
	var queries []url.Values
	c := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v2/news" {
			t.Errorf("path = %s", r.URL.Path)
		}
		q := r.URL.Query()
		queries = append(queries, q)
		w.Header().Set("Content-Type", "application/json")
		switch q.Get("cursor") {
		case "":
			_, _ = w.Write([]byte(`{"recordCount": null, "items": [{"id": "a"}, {"id": "b"}], "nextCursor": "next-1"}`))
		case "next-1":
			_, _ = w.Write([]byte(`{"recordCount": null, "items": [{"id": "c"}]}`))
		default:
			t.Errorf("unexpected cursor %q", q.Get("cursor"))
			w.WriteHeader(http.StatusBadRequest)
		}
	})

	var ids []string
	err := c.ForEachNewsV2(context.Background(), GetNewsV2Params{
		GetNewsParams: GetNewsParams{Region: RegionTr, Locale: LocaleTr, Symbols: "THYAO", Size: ptr(2)},
		EventTypes:    []NewsEventType{NewsEventTypeDividend},
	}, func(n NewsV2) error {
		ids = append(ids, n.ID)
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if strings.Join(ids, ",") != "a,b,c" {
		t.Errorf("ids = %v", ids)
	}
	if len(queries) != 2 {
		t.Fatalf("requests = %d, want 2", len(queries))
	}
	for i, q := range queries {
		for k, v := range map[string]string{"region": "tr", "locale": "tr", "symbols": "THYAO", "size": "2", "eventTypes": "dividend", "includeTotal": "false"} {
			if q.Get(k) != v {
				t.Errorf("request %d: %s = %q, want %q", i, k, q.Get(k), v)
			}
		}
		if q.Has("page") || q.Has("skip") {
			t.Errorf("request %d sent page/skip: %v", i, q)
		}
	}
	if queries[1].Get("cursor") != "next-1" {
		t.Errorf("second request cursor = %q", queries[1].Get("cursor"))
	}
}

func TestForEachNewsV2RejectsPage(t *testing.T) {
	c := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		t.Error("no request expected")
	})
	err := c.ForEachNewsV2(context.Background(), GetNewsV2Params{GetNewsParams: GetNewsParams{Page: ptr(1)}}, func(NewsV2) error { return nil })
	if err == nil {
		t.Fatal("want error when Page is set")
	}
}

func TestForEachNewsV2StopsOnError(t *testing.T) {
	calls := 0
	c := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		calls++
		if r.URL.Query().Get("cursor") == "" {
			_, _ = w.Write([]byte(`{"items": [{"id": "a"}], "nextCursor": "n"}`))
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"message": "invalid cursor", "error_code": "validation"}`))
	})

	err := c.ForEachNewsV2(context.Background(), GetNewsV2Params{}, func(NewsV2) error { return nil })
	var httpErr *LaplaceHTTPError
	if !errors.As(err, &httpErr) || httpErr.HTTPStatus != 400 || httpErr.Message.Message != "invalid cursor" {
		t.Fatalf("err = %v", err)
	}

	stop := errors.New("stop")
	calls = 0
	err = c.ForEachNewsV2(context.Background(), GetNewsV2Params{}, func(NewsV2) error { return stop })
	if !errors.Is(err, stop) || calls != 1 {
		t.Fatalf("err = %v, calls = %d", err, calls)
	}
}

func TestGetNewsV2PageSurfacesValidationMessage(t *testing.T) {
	const msg = `invalid eventTypes value "foo": must be one of financial-results, guidance`
	c := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"message": msg, "error_code": "validation"})
	})
	_, err := c.GetNewsV2Page(context.Background(), GetNewsV2Params{EventTypes: []NewsEventType{"foo"}})
	var httpErr *LaplaceHTTPError
	if !errors.As(err, &httpErr) || httpErr.Message.Message != msg || httpErr.Message.ErrorCode != "validation" {
		t.Fatalf("err = %v", err)
	}
}

func TestGetNewsEventTypes(t *testing.T) {
	c := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/news/event-types" || r.URL.Query().Get("locale") != "tr" {
			t.Errorf("request = %s", r.URL)
		}
		_, _ = w.Write([]byte(`[{"slug": "financial-results", "name": "Finansal Sonuçlar", "language": "tr"},
			{"slug": "guidance", "name": "Beklenti ve Hedefler", "language": "tr"}]`))
	})
	types, err := c.GetNewsEventTypes(context.Background(), LocaleTr)
	if err != nil {
		t.Fatal(err)
	}
	if len(types) != 2 || types[0] != (NewsEventTypeInfo{Slug: NewsEventTypeFinancialResults, Name: "Finansal Sonuçlar", Language: LocaleTr}) {
		t.Errorf("types = %+v", types)
	}
}

func TestNewsStreamSendsNewFilters(t *testing.T) {
	c := newMockClient(t, sseHandler(t, "/api/v1/news/stream", map[string]string{
		"locale":              "tr",
		"region":              "tr",
		"symbols":             "THYAO,GARAN",
		"tickerRelevance":     "primary",
		"excludeListArticles": "true",
		"eventTypes":          "dividend,guidance",
		"excludeEventTypes":   "other",
		"collapseClusters":    "false",
	}, "["+compactJSON(t, newsItemAllFields)+"]"))

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	stream, err := c.CreateNewsStream(ctx, StreamNewsParams{
		Region:              RegionTr,
		Locale:              LocaleTr,
		Symbols:             []string{"THYAO", "GARAN"},
		TickerRelevance:     []NewsTickerRelevance{NewsTickerRelevancePrimary},
		ExcludeListArticles: ptr(true),
		EventTypes:          []NewsEventType{NewsEventTypeDividend, NewsEventTypeGuidance},
		ExcludeEventTypes:   []NewsEventType{NewsEventTypeOther},
		CollapseClusters:    ptr(false),
	})
	if err != nil {
		t.Fatal(err)
	}
	defer stream.Close()

	select {
	case r := <-stream.Receive():
		if r.Error != nil {
			t.Fatal(r.Error)
		}
		if len(r.Data) != 1 || r.Data[0].PrimaryTicker != "THYAO" || r.Data[0].EventType == nil {
			t.Errorf("data = %+v", r.Data)
		}
	case <-ctx.Done():
		t.Fatal("timed out")
	}
}

func compactJSON(t *testing.T, s string) string {
	t.Helper()
	var b bytes.Buffer
	if err := json.Compact(&b, []byte(s)); err != nil {
		t.Fatal(err)
	}
	return b.String()
}

func TestNewsStreamOmitsUnsetFilters(t *testing.T) {
	c := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		for _, k := range []string{"tickerRelevance", "excludeListArticles", "eventTypes", "excludeEventTypes", "collapseClusters"} {
			if r.URL.Query().Has(k) {
				t.Errorf("unexpected %s", k)
			}
		}
		w.WriteHeader(http.StatusOK)
	})
	stream, err := c.CreateNewsStream(context.Background(), StreamNewsParams{Region: RegionUs, Locale: LocaleEn})
	if err != nil {
		t.Fatal(err)
	}
	stream.Close()
}

// Integration tests against the live API. They are skipped without an API key.

func newsIntegrationClient(t *testing.T) *Client {
	t.Helper()
	cfg, err := loadConfig()
	if err != nil || cfg.APIKey == "" {
		t.Skip("API key is not set; skipping live API test")
	}
	return newTestClient(*cfg)
}

func TestNewsV2IntegrationTickerNames(t *testing.T) {
	c := newsIntegrationClient(t)
	page, err := c.GetNewsV2Page(context.Background(), GetNewsV2Params{
		GetNewsParams: GetNewsParams{Region: RegionTr, Locale: LocaleTr, Symbols: "THYAO", Size: ptr(10)},
	})
	if err != nil {
		t.Fatal(err)
	}
	for _, item := range page.Items {
		for _, tk := range item.Tickers {
			if tk.Symbol == "THYAO" {
				if tk.Name == "" {
					t.Errorf("THYAO ticker has no name on %s", item.ID)
				}
				return
			}
		}
	}
	t.Fatal("no THYAO ticker in result")
}

func TestNewsV2IntegrationCursorRunEnds(t *testing.T) {
	c := newsIntegrationClient(t)
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	count := 0
	err := c.ForEachNewsV2(ctx, GetNewsV2Params{
		GetNewsParams: GetNewsParams{Region: RegionTr, Locale: LocaleTr, Symbols: "THYAO", Size: ptr(100)},
		PublishedFrom: time.Now().UTC().AddDate(0, 0, -7).Format("2006-01-02"),
	}, func(NewsV2) error {
		count++
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("cursor run returned %d items", count)
}

func TestNewsV2IntegrationInvalidEventType(t *testing.T) {
	c := newsIntegrationClient(t)
	_, err := c.GetNewsV2Page(context.Background(), GetNewsV2Params{
		GetNewsParams: GetNewsParams{Region: RegionTr, Locale: LocaleTr},
		EventTypes:    []NewsEventType{"foo"},
	})
	var httpErr *LaplaceHTTPError
	if !errors.As(err, &httpErr) || httpErr.HTTPStatus != http.StatusBadRequest || !strings.Contains(httpErr.Message.Message, "eventTypes") {
		t.Fatalf("err = %v", err)
	}
}

func TestNewsIntegrationEventTypes(t *testing.T) {
	c := newsIntegrationClient(t)
	types, err := c.GetNewsEventTypes(context.Background(), LocaleTr)
	if err != nil {
		t.Fatal(err)
	}
	if len(types) == 0 || types[0].Slug == "" || types[0].Name == "" {
		t.Errorf("types = %+v", types)
	}
}
