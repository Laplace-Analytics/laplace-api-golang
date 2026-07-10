package laplace

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type ScreenerSortBy string

// Sort keys. Any range-filter field is sortable, plus the smrRating/adRating
// letter grades, epsAcceleration and symbol. Results are ordered by the chosen
// column (NULLS LAST) then by symbol ASC.
const (
	ScreenerSortBySymbol            ScreenerSortBy = "symbol"
	ScreenerSortByPrice             ScreenerSortBy = "price"
	ScreenerSortByDailyChange       ScreenerSortBy = "dailyChange"
	ScreenerSortByMarketCap         ScreenerSortBy = "marketCap"
	ScreenerSortByPERatio           ScreenerSortBy = "peRatio"
	ScreenerSortByPBRatio           ScreenerSortBy = "pbRatio"
	ScreenerSortByWeeklyReturn      ScreenerSortBy = "weeklyReturn"
	ScreenerSortByMonthlyReturn     ScreenerSortBy = "monthlyReturn"
	ScreenerSortByThreeMonthReturn  ScreenerSortBy = "threeMonthReturn"
	ScreenerSortByYearlyReturn      ScreenerSortBy = "yearlyReturn"
	ScreenerSortByThreeYearReturn   ScreenerSortBy = "threeYearReturn"
	ScreenerSortByFiveYearReturn    ScreenerSortBy = "fiveYearReturn"
	ScreenerSortByYTDReturn         ScreenerSortBy = "ytdReturn"
	ScreenerSortByCompositeRating   ScreenerSortBy = "compositeRating"
	ScreenerSortByCompositeScore    ScreenerSortBy = "compositeScore"
	ScreenerSortByRSRating          ScreenerSortBy = "rsRating"
	ScreenerSortByRSScore           ScreenerSortBy = "rsScore"
	ScreenerSortByPerfQ1            ScreenerSortBy = "perfQ1"
	ScreenerSortByPerfQ2            ScreenerSortBy = "perfQ2"
	ScreenerSortByPerfQ3            ScreenerSortBy = "perfQ3"
	ScreenerSortByPerfQ4            ScreenerSortBy = "perfQ4"
	ScreenerSortByEPSRating         ScreenerSortBy = "epsRating"
	ScreenerSortByEPSScore          ScreenerSortBy = "epsScore"
	ScreenerSortByEPSGrowthYoY      ScreenerSortBy = "epsGrowthYoy"
	ScreenerSortByEPSGrowthQoQ      ScreenerSortBy = "epsGrowthQoq"
	ScreenerSortByEPSTrailing4Q     ScreenerSortBy = "epsTrailing4q"
	ScreenerSortByEPSAcceleration   ScreenerSortBy = "epsAcceleration"
	ScreenerSortByADRating          ScreenerSortBy = "adRating"
	ScreenerSortByADScore           ScreenerSortBy = "adScore"
	ScreenerSortByUpVolumeRatio     ScreenerSortBy = "upVolumeRatio"
	ScreenerSortByVolumeTrend       ScreenerSortBy = "volumeTrend"
	ScreenerSortBySMRRating         ScreenerSortBy = "smrRating"
	ScreenerSortBySMRScore          ScreenerSortBy = "smrScore"
	ScreenerSortBySalesGrowth2Q     ScreenerSortBy = "salesGrowth2q"
	ScreenerSortByGrossMargin       ScreenerSortBy = "grossMargin"
	ScreenerSortByNetMargin         ScreenerSortBy = "netMargin"
	ScreenerSortByROE               ScreenerSortBy = "roe"
	ScreenerSortBySMA20             ScreenerSortBy = "sma20"
	ScreenerSortBySMA50             ScreenerSortBy = "sma50"
	ScreenerSortBySMA150            ScreenerSortBy = "sma150"
	ScreenerSortBySMA200            ScreenerSortBy = "sma200"
	ScreenerSortByVolumeSMA50       ScreenerSortBy = "volumeSma50"
	ScreenerSortByPriceVsSMA20      ScreenerSortBy = "priceVsSma20"
	ScreenerSortByPriceVsSMA50      ScreenerSortBy = "priceVsSma50"
	ScreenerSortByPriceVsSMA150     ScreenerSortBy = "priceVsSma150"
	ScreenerSortByPriceVsSMA200     ScreenerSortBy = "priceVsSma200"
	ScreenerSortByHigh52W           ScreenerSortBy = "high52w"
	ScreenerSortByLow52W            ScreenerSortBy = "low52w"
	ScreenerSortByOffHighPct        ScreenerSortBy = "offHighPct"
	ScreenerSortByVolumeVsAvg50     ScreenerSortBy = "volumeVsAvg50"
	ScreenerSortByPriceChangePct    ScreenerSortBy = "priceChangePct"
	ScreenerSortByPriceChangeAmount ScreenerSortBy = "priceChangeAmount"
	ScreenerSortByYTDChangePct      ScreenerSortBy = "ytdChangePct"
)

// ScreenerRating is an IBD-style letter grade (A is best, E is worst) used by
// the SMR and Accumulation/Distribution ratings.
type ScreenerRating string

const (
	ScreenerRatingA ScreenerRating = "A"
	ScreenerRatingB ScreenerRating = "B"
	ScreenerRatingC ScreenerRating = "C"
	ScreenerRatingD ScreenerRating = "D"
	ScreenerRatingE ScreenerRating = "E"
)

// ScreenerRange is an inclusive [Min, Max] filter; both bounds are optional. If
// both are set, Min must be <= Max. Rows whose value is NULL in a column are
// excluded by any range filter touching that column.
type ScreenerRange struct {
	Min *float64 `json:"min,omitempty"`
	Max *float64 `json:"max,omitempty"`
}

// ScreenerFilters holds the optional filter set. Every field is a range filter
// except SMRRating/ADRating (letter-grade IN-lists) and EPSAcceleration (bool).
type ScreenerFilters struct {
	Price             *ScreenerRange `json:"price,omitempty"`
	DailyChange       *ScreenerRange `json:"dailyChange,omitempty"`
	MarketCap         *ScreenerRange `json:"marketCap,omitempty"`
	PERatio           *ScreenerRange `json:"peRatio,omitempty"`
	PBRatio           *ScreenerRange `json:"pbRatio,omitempty"`
	WeeklyReturn      *ScreenerRange `json:"weeklyReturn,omitempty"`
	MonthlyReturn     *ScreenerRange `json:"monthlyReturn,omitempty"`
	ThreeMonthReturn  *ScreenerRange `json:"threeMonthReturn,omitempty"`
	YearlyReturn      *ScreenerRange `json:"yearlyReturn,omitempty"`
	ThreeYearReturn   *ScreenerRange `json:"threeYearReturn,omitempty"`
	FiveYearReturn    *ScreenerRange `json:"fiveYearReturn,omitempty"`
	YTDReturn         *ScreenerRange `json:"ytdReturn,omitempty"`
	CompositeRating   *ScreenerRange `json:"compositeRating,omitempty"`
	CompositeScore    *ScreenerRange `json:"compositeScore,omitempty"`
	RSRating          *ScreenerRange `json:"rsRating,omitempty"`
	RSScore           *ScreenerRange `json:"rsScore,omitempty"`
	PerfQ1            *ScreenerRange `json:"perfQ1,omitempty"`
	PerfQ2            *ScreenerRange `json:"perfQ2,omitempty"`
	PerfQ3            *ScreenerRange `json:"perfQ3,omitempty"`
	PerfQ4            *ScreenerRange `json:"perfQ4,omitempty"`
	EPSRating         *ScreenerRange `json:"epsRating,omitempty"`
	EPSScore          *ScreenerRange `json:"epsScore,omitempty"`
	EPSGrowthYoY      *ScreenerRange `json:"epsGrowthYoy,omitempty"`
	EPSGrowthQoQ      *ScreenerRange `json:"epsGrowthQoq,omitempty"`
	EPSTrailing4Q     *ScreenerRange `json:"epsTrailing4q,omitempty"`
	ADScore           *ScreenerRange `json:"adScore,omitempty"`
	UpVolumeRatio     *ScreenerRange `json:"upVolumeRatio,omitempty"`
	VolumeTrend       *ScreenerRange `json:"volumeTrend,omitempty"`
	SMRScore          *ScreenerRange `json:"smrScore,omitempty"`
	SalesGrowth2Q     *ScreenerRange `json:"salesGrowth2q,omitempty"`
	GrossMargin       *ScreenerRange `json:"grossMargin,omitempty"`
	NetMargin         *ScreenerRange `json:"netMargin,omitempty"`
	ROE               *ScreenerRange `json:"roe,omitempty"`
	SMA20             *ScreenerRange `json:"sma20,omitempty"`
	SMA50             *ScreenerRange `json:"sma50,omitempty"`
	SMA150            *ScreenerRange `json:"sma150,omitempty"`
	SMA200            *ScreenerRange `json:"sma200,omitempty"`
	VolumeSMA50       *ScreenerRange `json:"volumeSma50,omitempty"`
	PriceVsSMA20      *ScreenerRange `json:"priceVsSma20,omitempty"`
	PriceVsSMA50      *ScreenerRange `json:"priceVsSma50,omitempty"`
	PriceVsSMA150     *ScreenerRange `json:"priceVsSma150,omitempty"`
	PriceVsSMA200     *ScreenerRange `json:"priceVsSma200,omitempty"`
	High52W           *ScreenerRange `json:"high52w,omitempty"`
	Low52W            *ScreenerRange `json:"low52w,omitempty"`
	OffHighPct        *ScreenerRange `json:"offHighPct,omitempty"`
	VolumeVsAvg50     *ScreenerRange `json:"volumeVsAvg50,omitempty"`
	PriceChangePct    *ScreenerRange `json:"priceChangePct,omitempty"`
	PriceChangeAmount *ScreenerRange `json:"priceChangeAmount,omitempty"`
	YTDChangePct      *ScreenerRange `json:"ytdChangePct,omitempty"`

	// Letter-grade IN-lists; values must be A..E.
	SMRRating []ScreenerRating `json:"smrRating,omitempty"`
	ADRating  []ScreenerRating `json:"adRating,omitempty"`

	// EPSAcceleration filters for stocks whose earnings growth is accelerating.
	EPSAcceleration *bool `json:"epsAcceleration,omitempty"`
}

type ScreenerRequest struct {
	Filters   *ScreenerFilters `json:"filters,omitempty"`
	SortBy    ScreenerSortBy   `json:"sortBy,omitempty"`
	SortOrder SortDirection    `json:"sortOrder,omitempty"`
	Page      int              `json:"page,omitempty"`
	PageSize  int              `json:"pageSize,omitempty"`
}

// ScreenerItem is one row of screener results. Decimal fields return 0 when a
// value is absent; the integer ratings, letter grades and EPSAcceleration are
// pointers because they are null when absent.
type ScreenerItem struct {
	Symbol            string          `json:"symbol"`
	Price             *float64        `json:"price"`
	DailyChange       *float64        `json:"dailyChange"`
	MarketCap         *float64        `json:"marketCap"`
	PERatio           *float64        `json:"peRatio"`
	PBRatio           *float64        `json:"pbRatio"`
	WeeklyReturn      *float64        `json:"weeklyReturn"`
	MonthlyReturn     *float64        `json:"monthlyReturn"`
	ThreeMonthReturn  *float64        `json:"threeMonthReturn"`
	YearlyReturn      *float64        `json:"yearlyReturn"`
	ThreeYearReturn   *float64        `json:"threeYearReturn"`
	FiveYearReturn    *float64        `json:"fiveYearReturn"`
	YTDReturn         *float64        `json:"ytdReturn"`
	CompositeRating   *int            `json:"compositeRating"`
	CompositeScore    *float64        `json:"compositeScore"`
	RSRating          *int            `json:"rsRating"`
	RSScore           *float64        `json:"rsScore"`
	PerfQ1            *float64        `json:"perfQ1"`
	PerfQ2            *float64        `json:"perfQ2"`
	PerfQ3            *float64        `json:"perfQ3"`
	PerfQ4            *float64        `json:"perfQ4"`
	EPSRating         *int            `json:"epsRating"`
	EPSScore          *float64        `json:"epsScore"`
	EPSGrowthYoY      *float64        `json:"epsGrowthYoy"`
	EPSGrowthQoQ      *float64        `json:"epsGrowthQoq"`
	EPSTrailing4Q     *float64        `json:"epsTrailing4q"`
	EPSAcceleration   *bool           `json:"epsAcceleration"`
	ADRating          *ScreenerRating `json:"adRating"`
	ADScore           *float64        `json:"adScore"`
	UpVolumeRatio     *float64        `json:"upVolumeRatio"`
	VolumeTrend       *float64        `json:"volumeTrend"`
	SMRRating         *ScreenerRating `json:"smrRating"`
	SMRScore          *float64        `json:"smrScore"`
	SalesGrowth2Q     *float64        `json:"salesGrowth2q"`
	GrossMargin       *float64        `json:"grossMargin"`
	NetMargin         *float64        `json:"netMargin"`
	ROE               *float64        `json:"roe"`
	SMA20             *float64        `json:"sma20"`
	SMA50             *float64        `json:"sma50"`
	SMA150            *float64        `json:"sma150"`
	SMA200            *float64        `json:"sma200"`
	VolumeSMA50       *float64        `json:"volumeSma50"`
	PriceVsSMA20      *float64        `json:"priceVsSma20"`
	PriceVsSMA50      *float64        `json:"priceVsSma50"`
	PriceVsSMA150     *float64        `json:"priceVsSma150"`
	PriceVsSMA200     *float64        `json:"priceVsSma200"`
	High52W           *float64        `json:"high52w"`
	Low52W            *float64        `json:"low52w"`
	OffHighPct        *float64        `json:"offHighPct"`
	VolumeVsAvg50     *float64        `json:"volumeVsAvg50"`
	PriceChangePct    *float64        `json:"priceChangePct"`
	PriceChangeAmount *float64        `json:"priceChangeAmount"`
	YTDChangePct      *float64        `json:"ytdChangePct"`
}

type ScreenerResponse struct {
	Items       []ScreenerItem `json:"items"`
	RecordCount int            `json:"recordCount"`
}

// Screener returns a filtered and sorted list of stocks for the given region.
// Region defaults to "tr" if empty. US is not currently supported.
func (c *Client) Screener(ctx context.Context, region Region, params ScreenerRequest) (ScreenerResponse, error) {
	bodyJSON, err := json.Marshal(params)
	if err != nil {
		return ScreenerResponse{}, err
	}

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/api/v1/screener", c.baseUrl), bytes.NewBuffer(bodyJSON))
	if err != nil {
		return ScreenerResponse{}, err
	}

	if region != "" {
		q := req.URL.Query()
		q.Add("region", string(region))
		req.URL.RawQuery = q.Encode()
	}

	req.Header.Set("Content-Type", "application/json")

	return sendRequest[ScreenerResponse](ctx, c, req)
}
